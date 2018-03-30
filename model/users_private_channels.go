package model

import (
	"github.com/traPtitech/traQ/utils/validator"
)

// UsersPrivateChannel : UsersPrivateChannelsの構造体
type UsersPrivateChannel struct {
	UserID    string `xorm:"char(36) pk" validate:"uuid,required"`
	ChannelID string `xorm:"char(36) pk" validate:"uuid,required"`
}

// TableName : テーブル名を指定するメソッド
func (upc *UsersPrivateChannel) TableName() string {
	return "users_private_channels"
}

// Validate 構造体を検証します
func (upc *UsersPrivateChannel) Validate() error {
	return validator.ValidateStruct(upc)
}

// Create : データベースへ反映
func (upc *UsersPrivateChannel) Create() (err error) {
	if err = upc.Validate(); err != nil {
		return err
	}

	_, err = db.InsertOne(upc)
	return
}

// GetPrivateChannel ある二つのユーザー間のプライベートチャンネルが存在するかを調べる
func GetPrivateChannel(userID1, userID2 string) (*UsersPrivateChannel, error) {
	upc := &UsersPrivateChannel{}
	has, err := db.Table(upc).Join("LEFT", []string{"users_private_channels", "p"}, "p.user_id IN(?, ?) AND users_private_channels.user_id IN(?, ?) AND p.channel_id = users_private_channels.channel_id", userID1, userID2, userID1, userID2).
		Get(upc)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	return upc, nil
}
