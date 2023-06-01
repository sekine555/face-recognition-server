package aws

import (
	"face-recognition/config"
	"face-recognition/logger"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"go.uber.org/zap"
	"strconv"
)

// AWS Rekognition顔認証（比較）
func CompareFaces(sourceImageS3Key string, targetImageS3Key string) (similarity float64, err error) {
	// セッション作成
	cred := credentials.NewStaticCredentials(config.Config.AccessKeyId, config.Config.SecretAccessKey, "")
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: cred,
		Region: aws.String(config.Config.Region),
	}))
	// 解析オブジェクト作成
	svc := rekognition.New(sess)
	// パラメータセット
	input := &rekognition.CompareFacesInput{
		// 認識度（高いほど厳しい：0-100）
		SimilarityThreshold: aws.Float64(90.000000),
		SourceImage: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(config.Config.Bucket),
				Name: aws.String(sourceImageS3Key),
			},
		},
		TargetImage: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(config.Config.Bucket),
				Name: aws.String(targetImageS3Key),
			},
		},
	}
	// 顔比較実行
	response, err := svc.CompareFaces(input)
	if err != nil {
		logger.Log.Info("rekognitionのCompareFacesエラー")
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			// 顔が写っていないエラー
			case rekognition.ErrCodeInvalidParameterException:
				logger.Log.Info("判定画像に顔が写っていない")
				return 0, nil
			default:
				logger.Log.Info(err.Error())
				return 0, err
			}
		}
		return 0, err
	}
	fmt.Println("認証結果：", response)
	logger.Log.Info(response.String())
	// 認証結果判定
	if len(response.FaceMatches) != 0 {
		ret := *response.FaceMatches[0].Similarity
		logger.Log.Info("顔認証結果OK", zap.String("認識度", strconv.FormatFloat(ret, 'f', 2, 64)))
		return ret, nil
	} else {
		logger.Log.Info("顔認証結果NG")
		return 0, nil
	}
}