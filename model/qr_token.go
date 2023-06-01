package model

import (
	"time"
)

type QrToken struct {
	Id          float64   `json:"id"`
	MstUser     MstUser   `gorm:"foreignkey:MstUserId" json:"mstUser"`
	MstUserId   float64   `json:"mstUserId"`
	QrToken     string    `json:"qrToken"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"-"`
}

func (QrToken) TableName() string {
	return "qr_token"
}

// 顔認証APIのRequestBody
type FaceRecognitionParams struct {
	QrToken  string `json:"qrToken" validate:"required"`
	Photo    string `json:"photo" validate:"required,base64"`
}