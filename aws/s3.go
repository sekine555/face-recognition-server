package aws

import (
	"bytes"
	"encoding/base64"
	"errors"
	"face-recognition/config"
	"face-recognition/logger"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/zap"
)

type ImageUploadToS3 struct {
	Region          string
	BucketName      string
	AccessKeyId     string
	SecretAccessKey string
}

// AWS S3へアップロード
func PutToS3(imageBase64 string, fileName string, extension string) (url string, err error) {
	// アップロードファイル名、写真（base64形式）チェック
	if fileName == "" {
		return "", errors.New("ファイルは必須項目です")
	} else if imageBase64 == "" {
		return "", errors.New("写真は必須項目です")
	}
	// Content-Typeの設定
	var contentType string
	switch extension {
	case "jpg":
		contentType = "image/jpeg"
	case "jpeg":
		contentType = "image/jpeg"
	case "png":
		contentType = "image/png"
	default:
		return "", errors.New("拡張子が無効です")
	}

	// セッション作成
	cred := credentials.NewStaticCredentials(config.Config.AccessKeyId, config.Config.SecretAccessKey, "")
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: cred,
		Region: aws.String(config.Config.Region),
	}))

	// base64から画像生成（バッファに書き込む）
	data, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		logger.Log.Info("base64から画像生成エラー", zap.String("写真", imageBase64))
		return "", err
	}
	wb := new(bytes.Buffer)
	wb.Write(data)

	// Uploaderを作成
	uploader := s3manager.NewUploader(sess)
	resp, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.Config.Bucket),
		Key:    aws.String(fileName),
		ContentType: aws.String(contentType),
		Body:   wb,
	})
	if err != nil {
		logger.Log.Info("S3画像アップロードエラー", zap.String("ファイル", fileName))
		return "", err
	}
	fmt.Println("結果：", resp.Location)
	return resp.Location, nil
}