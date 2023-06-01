package api

import (
	"face-recognition/aws"
	"face-recognition/config"
	"face-recognition/db"
	"face-recognition/logger"
	"face-recognition/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/rs/xid"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
	"time"
)

// QRトークン取得
func GetQrToken() echo.HandlerFunc {
	return func(context echo.Context) error {
		logger.Log.Info("QRトークン取得API開始")
		// DB接続
		db, err := db.SqlConnect()
		if err != nil {
			return context.String(http.StatusBadGateway, err.Error())
		} else {
			// ログ出力有効
			db.LogMode(true)
			fmt.Println("Successfully connect database..")
		}
		// DBクローズ（遅延）
		defer db.Close()

		// トークンからユーザ特定
		user := context.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		userId := claims["userId"].(float64)
		logger.Log.Info("QRトークン取得API", zap.String("User", strconv.FormatFloat(userId,'f', -1, 64)))
		// qr_tokenテーブルにレコードが存在すれば返却する
		qrToken := model.QrToken{}
		// プリロードを利用すれば、1センテンスで複数のテーブルからデータを取得
		// モデル間の関係を持っていることが前提
		db.Where("mst_user_id = ?", userId).Preload("MstUser").Find(&qrToken)
		// QRトークンが存在しなければ、データを作成して返却する
		if qrToken.Id == 0 {
			logger.Log.Info("qr_tokenテーブルにレコードが存在しないため、登録データを作成")
			// トークン生成
			token := jwt.New(jwt.SigningMethodHS256)
			claims := token.Claims.(jwt.MapClaims)
			claims["iat"] = time.Now()
			claims["userId"] = userId
			// 署名
			t, err := token.SignedString([]byte(config.Config.Secret))
			if err != nil {
				return err
			}
			// qr_tokenテーブルへの投入データ作成
			qrToken.QrToken = t
			qrToken.MstUserId = userId
			// トランザクション開始
			tx := db.Begin()
			defer tx.Close()
			if err := tx.Create(&qrToken).Error; err != nil {
				logger.Log.Info("QRトークンテーブル登録失敗")
				logger.Log.Info("QRトークン取得API終了", zap.String("User", strconv.FormatFloat(userId,'f', -1, 64)))
				return context.JSON(http.StatusInternalServerError, map[string]interface{}{
					"message": "QRトークンテーブルへ登録できませんでした",
				})
			}
			// コミット
			tx.Commit()
			db.Preload("MstUser").Find(&qrToken)
			logger.Log.Info("QRトークン取得API終了", zap.String("User", strconv.FormatFloat(userId,'f', -1, 64)))
			return context.JSON(http.StatusOK, qrToken)
		}
		// 既にQRトークンテーブルにレコードが存在する場合にはそのトークンを返却する
		logger.Log.Info("QRトークン取得API終了", zap.String("User", strconv.FormatFloat(userId,'f', -1, 64)))
		return context.JSON(http.StatusOK, qrToken)
	}
}

// 顔認証
func PostFaceRecognition() echo.HandlerFunc {
	return func(context echo.Context) error {
		logger.Log.Info("顔認証API開始")
		// DB接続
		db, err := db.SqlConnect()
		if err != nil {
			return context.String(http.StatusBadGateway, err.Error())
		} else {
			// ログ出力有効
			db.LogMode(true)
			fmt.Println("Successfully connect database..")
		}
		// DBクローズ（遅延）
		defer db.Close()
		// リクエストボディーを構造体にバインド
		face := new(model.FaceRecognitionParams)
		if err = context.Bind(face); err != nil {
			logger.Log.Info("顔認証情報パラメータバインド失敗")
			logger.Log.Info("顔認証API終了")
			return context.JSON(http.StatusBadRequest, err.Error())
		}
		// バリデーション
		validate := validator.New()
		var errorMessages []string
		if err := validate.Struct(face); err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				var errMsg string
				fieldName := err.Field()
				switch fieldName {
				case "QrToken":
					errMsg = "QRトークンは必須項目です"
				case "Photo":
					var tag = err.Tag()
					fmt.Println("tag", tag)
					switch tag {
					case "required":
						errMsg = "写真は必須項目です"
					case "base64":
						errMsg = "写真のフォーマットが不正です"

					}
				}
				errorMessages = append(errorMessages, errMsg)
			}
			fmt.Println("errMsg：", errorMessages)
			logger.Log.Info("パラメータエラー", zap.Strings("エラー内容", errorMessages))
			logger.Log.Info("顔認証API終了")
			return context.JSON(http.StatusBadRequest, errorMessages)
		}

		// qrトークン（jwt）をデコード
		tokenString := face.QrToken
		claims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Config.Secret), nil
		})
		if err != nil {
			logger.Log.Info("QRトークンからデコード失敗")
			logger.Log.Info("顔認証API終了", zap.String("QRトークン", face.QrToken))
			return context.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "QRトークンのフォーマットが不正のため、特定できませんでした",
			})
		}
		if len(claims) == 0 {
			logger.Log.Info("QRトークンにクレーム情報がありません")
			logger.Log.Info("顔認証API終了", zap.String("QRトークン", face.QrToken))
			return context.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "QRトークンのフォーマットが不正のため、特定できませんでした",
			})
		}
		var userId float64 = -1
		for key, val := range claims {
			fmt.Printf("Key: %v, value: %v\n", key, val)
			if key == "userId" {
				userId = val.(float64)
			}
		}
		if userId == -1 {
			logger.Log.Info("QRトークンのクレーム情報にuserIdがありません")
			logger.Log.Info("顔認証API終了", zap.String("QRトークン", face.QrToken))
			return context.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "QRトークンからユーザーを特定できませんでした",
			})
		}
		logger.Log.Info("顔認証API", zap.String("認証User", strconv.FormatFloat(userId,'f', -1, 64)))

		// 認証対象ユーザのプロフィール画像取得
		mstUser := model.MstUser{}
		db.Where("id = ?", userId).Find(&mstUser)
		if mstUser.Id == 0 {
			logger.Log.Info("ユーザマスタに存在しません", zap.String("認証User", strconv.FormatFloat(userId,'f', -1, 64)))
			logger.Log.Info("顔認証API終了", zap.String("QRトークン", face.QrToken))
			return context.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "QRトークンからユーザーを特定できませんでした",
			})
		}
		// 比較対象画像をS3へアップロード
		// 一意なファイル名（S3のキー）生成
		fileId := xid.New()
		logger.Log.Info("ファイル名", zap.String("fileId", fileId.String()))
		s3url, err := aws.PutToS3(face.Photo, fileId.String(), "png")
		if err != nil {
			logger.Log.Info("比較先画像のS3アップロードエラー発生")
			return context.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "比較先画像をアップロードできませんでした",
			})
		}
		logger.Log.Info("生成した画像URL", zap.String("s3Url", s3url))
		// 顔認証実施
		resp, err := aws.CompareFaces(mstUser.S3Key, fileId.String())
		if err != nil {
			logger.Log.Info("顔認証失敗", zap.String("error", err.Error()))
			logger.Log.Info("顔認証API終了", zap.String("QRトークン", face.QrToken))
			return context.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "顔認証に失敗しました",
			})
		}
		// 顔認証結果テーブルへ投入
		createFaceRecognitionResult := model.FaceRecognitionResult{}
		createFaceRecognitionResult.MstUserId = mstUser.Id
		createFaceRecognitionResult.SourceImage = mstUser.Photo
		createFaceRecognitionResult.SourceImageS3Key = mstUser.S3Key
		createFaceRecognitionResult.TargetImage = s3url
		createFaceRecognitionResult.TargetImageS3Key = fileId.String()
		createFaceRecognitionResult.Result = resp
		// トランザクション開始
		tx := db.Begin()
		defer tx.Close()
		if err := tx.Create(&createFaceRecognitionResult).Error; err != nil {
			logger.Log.Info("顔認証結果登録失敗")
			logger.Log.Info("顔認証API終了")
			return context.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "顔認証結果テーブルへ登録できませんでした",
			})
		}
		// コミット
		tx.Commit()
		authResult := false
		if resp != 0 {
			authResult = true
		}
		logger.Log.Info( "顔認証API終了", zap.String("結果", strconv.FormatBool(authResult)))
		return context.JSON(http.StatusOK, map[string]interface{} {
			"authResult": authResult,
		})
	}
}