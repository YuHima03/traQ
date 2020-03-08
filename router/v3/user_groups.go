package v3

import (
	vd "github.com/go-ozzo/ozzo-validation"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traQ/rbac/permission"
	"github.com/traPtitech/traQ/repository"
	"github.com/traPtitech/traQ/router/consts"
	"github.com/traPtitech/traQ/router/extension/herror"
	"github.com/traPtitech/traQ/utils/validator"
	"gopkg.in/guregu/null.v3"
	"net/http"
)

// PostUserGroupRequest POST /groups リクエストボディ
type PostUserGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

func (r PostUserGroupRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.Name, vd.Required, vd.RuneLength(1, 30)),
		vd.Field(&r.Description, vd.RuneLength(0, 100)),
		vd.Field(&r.Type, vd.RuneLength(0, 30)),
	)
}

// PostUserGroups POST /groups
func (h *Handlers) PostUserGroups(c echo.Context) error {
	reqUserID := getRequestUserID(c)

	var req PostUserGroupRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	if req.Type == "grade" && !h.RBAC.IsGranted(getRequestUser(c).Role, permission.CreateSpecialUserGroup) {
		// 学年グループは権限が必要
		return herror.Forbidden("you are not permitted to create groups of this type")
	}

	g, err := h.Repo.CreateUserGroup(req.Name, req.Description, req.Type, reqUserID)
	if err != nil {
		switch {
		case err == repository.ErrAlreadyExists:
			return herror.Conflict("the name's group has already existed")
		case repository.IsArgError(err):
			return herror.BadRequest(err)
		default:
			return herror.InternalServerError(err)
		}
	}

	return c.JSON(http.StatusCreated, formatUserGroup(g))
}

// GetUserGroup GET /groups/:groupID
func (h *Handlers) GetUserGroup(c echo.Context) error {
	return c.JSON(http.StatusOK, formatUserGroup(getParamGroup(c)))
}

// PatchUserGroupRequest PATCH /groups/:groupID リクエストボディ
type PatchUserGroupRequest struct {
	Name        null.String `json:"name"`
	Description null.String `json:"description"`
	Type        null.String `json:"type"`
}

func (r PatchUserGroupRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.Name, vd.RuneLength(1, 30)),
		vd.Field(&r.Description, vd.RuneLength(0, 100)),
		vd.Field(&r.Type, vd.RuneLength(0, 30)),
	)
}

// EditUserGroup PATCH /groups/:groupID
func (h *Handlers) EditUserGroup(c echo.Context) error {
	g := getParamGroup(c)

	var req PatchUserGroupRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	if req.Type.ValueOrZero() == "grade" && !h.RBAC.IsGranted(getRequestUser(c).Role, permission.CreateSpecialUserGroup) {
		// 学年グループは権限が必要
		return herror.Forbidden("you are not permitted to create groups of this type")
	}

	args := repository.UpdateUserGroupNameArgs{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
	}
	if err := h.Repo.UpdateUserGroup(g.ID, args); err != nil {
		switch {
		case err == repository.ErrAlreadyExists:
			return herror.Conflict("the name's group has already existed")
		case repository.IsArgError(err):
			return herror.BadRequest(err)
		default:
			return herror.InternalServerError(err)
		}
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteUserGroup DELETE /groups/:groupID
func (h *Handlers) DeleteUserGroup(c echo.Context) error {
	g := getParamGroup(c)

	if err := h.Repo.DeleteUserGroup(g.ID); err != nil {
		return herror.InternalServerError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetUserGroupMembers GET /groups/:groupID/members
func (h *Handlers) GetUserGroupMembers(c echo.Context) error {
	return c.JSON(http.StatusOK, formatUserGroupMembers(getParamGroup(c).Members))
}

// RemoveUserGroupMember DELETE /groups/:groupID/admins/:userID
func (h *Handlers) RemoveUserGroupMember(c echo.Context) error {
	userID := getParamAsUUID(c, consts.ParamUserID)
	g := getParamGroup(c)

	if err := h.Repo.RemoveUserFromGroup(userID, g.ID); err != nil {
		return herror.InternalServerError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetUserGroupAdmins GET /groups/:groupID/admins
func (h *Handlers) GetUserGroupAdmins(c echo.Context) error {
	g := getParamGroup(c)
	result := make([]uuid.UUID, 0)
	for _, admin := range g.Admins {
		result = append(result, admin.UserID)
	}
	return c.JSON(http.StatusOK, result)
}

// PostUserGroupAdminRequest POST /groups/:groupID/admins リクエストボディ
type PostUserGroupAdminRequest struct {
	ID uuid.UUID `json:"id"`
}

func (r PostUserGroupAdminRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.ID, vd.Required, validator.NotNilUUID),
	)
}

// AddUserGroupAdmin POST /groups/:groupID/admins
func (h *Handlers) AddUserGroupAdmin(c echo.Context) error {
	g := getParamGroup(c)

	var req PostUserGroupAdminRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	// ユーザーが存在するか
	if ok, err := h.Repo.UserExists(req.ID); err != nil {
		return herror.InternalServerError(err)
	} else if !ok {
		return herror.BadRequest("this user doesn't exist")
	}

	if err := h.Repo.AddUserToGroupAdmin(req.ID, g.ID); err != nil {
		return herror.InternalServerError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// RemoveUserGroupAdmin DELETE /groups/:groupID/admins/:userID
func (h *Handlers) RemoveUserGroupAdmin(c echo.Context) error {
	userID := getParamAsUUID(c, consts.ParamUserID)
	g := getParamGroup(c)

	if err := h.Repo.RemoveUserFromGroupAdmin(userID, g.ID); err != nil {
		if err == repository.ErrForbidden {
			return herror.BadRequest()
		}
		return herror.InternalServerError(err)
	}

	return c.NoContent(http.StatusNoContent)
}
