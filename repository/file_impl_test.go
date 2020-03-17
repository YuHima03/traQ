package repository

import (
	"bytes"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/traQ/model"
	"testing"
)

func TestRepositoryImpl_GenerateIconFile(t *testing.T) {
	t.Parallel()
	repo, assert, require := setup(t, common)

	id, err := repo.GenerateIconFile("salt")
	if assert.NoError(err) {
		meta, err := repo.GetFileMeta(id)
		require.NoError(err)
		assert.Equal(model.FileTypeIcon, meta.Type)
	}
}

func TestRepositoryImpl_DeleteFile(t *testing.T) {
	t.Parallel()
	repo, _, _ := setup(t, common)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		id, err := repo.GenerateIconFile("test")
		require.NoError(t, err)

		if assert.NoError(t, repo.DeleteFile(id)) {
			_, err := repo.GetFileMeta(id)
			assert.EqualError(t, err, ErrNotFound.Error())
		}
	})

	t.Run("nil id", func(t *testing.T) {
		t.Parallel()

		assert.EqualError(t, repo.DeleteFile(uuid.Nil), ErrNilID.Error())
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		assert.EqualError(t, repo.DeleteFile(uuid.Must(uuid.NewV4())), ErrNotFound.Error())
	})
}

func TestRepositoryImpl_GetFileMeta(t *testing.T) {
	t.Parallel()
	repo, assert, _ := setup(t, common)

	file := mustMakeFile(t, repo)
	result, err := repo.GetFileMeta(file.ID)
	if assert.NoError(err) {
		assert.Equal(file.ID, result.ID)
	}

	_, err = repo.GetFileMeta(uuid.Nil)
	assert.Error(err)
}

func TestRepositoryImpl_SaveFile(t *testing.T) {
	t.Parallel()
	repo, assert, _ := setup(t, common)

	buf := bytes.NewBufferString("test message")
	f, err := repo.SaveFile(SaveFileArgs{
		FileName: "test.txt",
		FileSize: int64(buf.Len()),
		FileType: model.FileTypeUserFile,
		Src:      buf,
	})
	if assert.NoError(err) {
		assert.Equal("text/plain; charset=utf-8", f.Mime)
	}
}

func TestRepositoryImpl_IsFileAccessible(t *testing.T) {
	t.Parallel()
	repo, _, _, user := setupWithUser(t, common)

	t.Run("invalid args", func(t *testing.T) {
		t.Parallel()

		_, err := repo.IsFileAccessible(uuid.Nil, uuid.Nil)
		assert.EqualError(t, err, ErrNilID.Error())
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		_, err := repo.IsFileAccessible(uuid.Must(uuid.NewV4()), user.GetID())
		assert.EqualError(t, err, ErrNotFound.Error())
	})

	t.Run("Allow all", func(t *testing.T) {
		t.Parallel()
		f := mustMakeFile(t, repo)

		t.Run("any user", func(t *testing.T) {
			t.Parallel()

			ok, err := repo.IsFileAccessible(f.ID, uuid.Nil)
			if assert.NoError(t, err) {
				assert.True(t, ok)
			}
		})

		t.Run("user", func(t *testing.T) {
			t.Parallel()

			ok, err := repo.IsFileAccessible(f.ID, user.GetID())
			if assert.NoError(t, err) {
				assert.True(t, ok)
			}
		})
	})

	t.Run("Allow one", func(t *testing.T) {
		t.Parallel()

		buf := bytes.NewBufferString("test message")
		args := SaveFileArgs{
			FileName: "test.txt",
			FileSize: int64(buf.Len()),
			FileType: model.FileTypeUserFile,
			Src:      buf,
			ACL:      ACL{},
		}
		args.SetCreator(user.GetID())
		f, err := repo.SaveFile(args)
		require.NoError(t, err)

		t.Run("any user", func(t *testing.T) {
			t.Parallel()

			ok, err := repo.IsFileAccessible(f.ID, uuid.Nil)
			if assert.NoError(t, err) {
				assert.False(t, ok)
			}
		})

		t.Run("allowed user", func(t *testing.T) {
			t.Parallel()

			ok, err := repo.IsFileAccessible(f.ID, user.GetID())
			if assert.NoError(t, err) {
				assert.True(t, ok)
			}
		})

		t.Run("denied user", func(t *testing.T) {
			t.Parallel()

			user := mustMakeUser(t, repo, random)
			ok, err := repo.IsFileAccessible(f.ID, user.GetID())
			if assert.NoError(t, err) {
				assert.False(t, ok)
			}
		})
	})

	t.Run("Allow two", func(t *testing.T) {
		t.Parallel()

		user2 := mustMakeUser(t, repo, random)
		buf := bytes.NewBufferString("test message")
		args := SaveFileArgs{
			FileName: "test.txt",
			FileSize: int64(buf.Len()),
			FileType: model.FileTypeUserFile,
			Src:      buf,
			ACL:      ACL{user2.GetID(): true},
		}
		args.SetCreator(user.GetID())
		f, err := repo.SaveFile(args)
		require.NoError(t, err)

		t.Run("any user", func(t *testing.T) {
			t.Parallel()

			ok, err := repo.IsFileAccessible(f.ID, uuid.Nil)
			if assert.NoError(t, err) {
				assert.False(t, ok)
			}
		})

		t.Run("allowed user", func(t *testing.T) {
			t.Parallel()

			ok, err := repo.IsFileAccessible(f.ID, user.GetID())
			if assert.NoError(t, err) {
				assert.True(t, ok)
			}
		})

		t.Run("allowed user2", func(t *testing.T) {
			t.Parallel()

			ok, err := repo.IsFileAccessible(f.ID, user2.GetID())
			if assert.NoError(t, err) {
				assert.True(t, ok)
			}
		})

		t.Run("denied user", func(t *testing.T) {
			t.Parallel()

			user := mustMakeUser(t, repo, random)
			ok, err := repo.IsFileAccessible(f.ID, user.GetID())
			if assert.NoError(t, err) {
				assert.False(t, ok)
			}
		})
	})

	t.Run("Deny rule", func(t *testing.T) {
		t.Parallel()

		deninedUser := mustMakeUser(t, repo, random)
		buf := bytes.NewBufferString("test message")
		args := SaveFileArgs{
			FileName: "test.txt",
			FileSize: int64(buf.Len()),
			FileType: model.FileTypeUserFile,
			Src:      buf,
			ACL: ACL{
				uuid.Nil:            true,
				deninedUser.GetID(): false,
			},
		}
		f, err := repo.SaveFile(args)
		require.NoError(t, err)

		t.Run("any user", func(t *testing.T) {
			t.Parallel()

			ok, err := repo.IsFileAccessible(f.ID, uuid.Nil)
			if assert.NoError(t, err) {
				assert.True(t, ok)
			}
		})

		t.Run("allowed user", func(t *testing.T) {
			t.Parallel()

			ok, err := repo.IsFileAccessible(f.ID, user.GetID())
			if assert.NoError(t, err) {
				assert.True(t, ok)
			}
		})

		t.Run("denied user", func(t *testing.T) {
			t.Parallel()

			ok, err := repo.IsFileAccessible(f.ID, deninedUser.GetID())
			if assert.NoError(t, err) {
				assert.False(t, ok)
			}
		})
	})
}
