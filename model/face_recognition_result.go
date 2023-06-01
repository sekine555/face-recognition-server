package model

import "time"

type FaceRecognitionResult struct {
	Id               float64 `json:"id"`
	MstUser          MstUser `gorm:"foreignkey:MstUserId" json:"mstUser"`
	MstUserId        float64 `json:"mstUserId"`
	SourceImage      string  `json:"sourceImage"`
	SourceImageS3Key string  `json:"sourceImageS3Key"`
	TargetImage      string  `json:"targetImage"`
	TargetImageS3Key string  `json:"targetImageS3Key"`
	Result           float64 `json:"result"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"-"`
}

func (FaceRecognitionResult) TableName() string {
	return "face_recognition_result"
}