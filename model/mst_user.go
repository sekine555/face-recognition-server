package model

import (
	"time"
)

// memo：LastLoginAtはポインタをつけないと、テーブルの列値がnullの時、omitemptyを指定しても
// 最初の日付（0000年...）が返却されてしまう
type MstUser struct {
	Id          float64  `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"-"`
	Photo       string `json:"photo"`
	S3Key       string `json:"s3Key"`
	IsAdmin     bool   `json:"isAdmin"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"-"`
}

// GORMではテーブル名が複数形になってしまうため、実テーブル名を明示する
func (MstUser) TableName() string {
	return "mst_user"
}

// ログイン認証APIのRequestBody
type LoginParams struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// ユーザ登録APIのRequestBody
type UserParams struct {
	Email string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Photo    string `json:"photo" validate:"required,base64"`
}
