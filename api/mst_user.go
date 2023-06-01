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
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"time"
)

// ユーザ情報取得
func GetUser() echo.HandlerFunc {
	return func(context echo.Context) error {
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
		// ユーザマスタからレコード取得
		// 結果を受け取るMstUser型の空のスライスを用意しておき、db.Findの引数でそのアドレスを渡す
		var users []model.MstUser
		result := db.Find(&users)
		return context.JSON(http.StatusOK, result)
	}
}

// ログイン認証
func PostLogin() echo.HandlerFunc {
	return func(context echo.Context) error {
		logger.Log.Info("ログイン認証API開始")
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
		u := new(model.LoginParams)
		if err = context.Bind(u); err != nil {
			logger.Log.Info("ログイン情報パラメータバインド失敗")
			logger.Log.Info("ログイン認証API終了")
			return context.JSON(http.StatusBadRequest, err.Error())
		}
		// バリデーション
		validate := validator.New()
		if err := validate.Struct(u); err != nil {
			// バリデーションエラーメッセージの加工
			var errorMessages []string
			for _, err := range err.(validator.ValidationErrors) {
				var errMsg string
				fieldName := err.Field()
				switch fieldName {
				case "Email":
					var tag = err.Tag()
					switch tag {
					case "required":
						errMsg = "メールアドレスは必須項目です"
					case "email":
						errMsg = "メールアドレスのフォーマットが不正です"
					}
				case "Password":
					errMsg = "パスワードは必須項目です"
				}
				errorMessages = append(errorMessages, errMsg)
			}
			fmt.Println("errMsg：", errorMessages)
			logger.Log.Info("パラメータエラー", zap.Strings("エラー内容", errorMessages))
			logger.Log.Info("ログイン認証API終了")
			return context.JSON(http.StatusBadRequest, errorMessages)
		}
		// jwt認証
		var user []model.MstUser
		db.Where("email = ?", u.Email).Find(&user)
		if len(user) > 0 && u.Email == user[0].Email {
			// パスワード検証
			if !compareHashedPassword(user[0].Password, u.Password) {
				// ログイン認証エラー（パスワード誤り）
				logger.Log.Info("パスワードが違います")
				logger.Log.Info("ログイン認証API終了")
				return echo.ErrUnauthorized
			}
			// トークン生成
			// ヘッダのセット
			token := jwt.New(jwt.SigningMethodHS256)
			// クレームのセット
			claims := token.Claims.(jwt.MapClaims)
			// トークンの発行日時
			claims["iat"] = time.Now()
			// トークンの有効期限（3日）
			claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
			claims["userId"] = user[0].Id
			// 署名
			t, err := token.SignedString([]byte(config.Config.Secret))
			if err != nil {
				return err
			}
			logger.Log.Info("ログイン認証API終了")
			return context.JSON(http.StatusOK, map[string]interface{}{
				"token": t,
				"admin": user[0].IsAdmin,
			})
		} else {
			// ログイン認証エラー（ユーザ情報なし）
			logger.Log.Info("メールアドレスかパスワードが違います")
			logger.Log.Info("ログイン認証API終了")
			return echo.ErrUnauthorized
		}
	}
}

// ユーザ登録
func PostUser() echo.HandlerFunc {
	return func(context echo.Context) error {
		logger.Log.Info("ユーザ登録API開始")
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
		u := new(model.UserParams)
		if err = context.Bind(u); err != nil {
			logger.Log.Info("ユーザ登録パラメータバインド失敗")
			logger.Log.Info("ユーザ登録API終了")
			return context.JSON(http.StatusBadRequest, err.Error())
		}
		// バリデーション
		validate := validator.New()
		// バリデーションエラーメッセージの加工
		var errorMessages []string
		if err := validate.Struct(u); err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				var errMsg string
				fieldName := err.Field()
				switch fieldName {
				case "Email":
					var tag = err.Tag()
					switch tag {
					case "required":
						errMsg = "メールアドレスは必須項目です"
					case "email":
						errMsg = "メールアドレスのフォーマットが不正です"
					}
				case "Username":
					errMsg = "ユーザー名は必須項目です"
				case "Password":
					errMsg = "パスワードは必須項目です"
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
			logger.Log.Info("ユーザ登録API終了")
			return context.JSON(http.StatusBadRequest, errorMessages)
		}
		// メールアドレス重複チェック
		chkUser := model.MstUser{}
		var count int = 0
		db.Model(&chkUser).Where("email = ?", u.Email).Count(&count)
		if count > 0 {
			errorMessages = append(errorMessages, "既に存在するメールアドレスです")
			logger.Log.Info("パラメータエラー", zap.Strings("エラー内容", errorMessages))
			logger.Log.Info("ユーザ登録API終了")
			return context.JSON(http.StatusBadRequest, errorMessages)
		}
		// 画像アップロード
		// 一意なファイル名生成
		fileId := xid.New()
		logger.Log.Info("ファイル名", zap.String("fileId", fileId.String()))
		s3url, err := aws.PutToS3(u.Photo, fileId.String(), "png")
		if err != nil {
			logger.Log.Info("アップロードエラー発生")
			return context.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "画像を登録できませんでした",
			})
		}
		logger.Log.Info("生成した画像URL", zap.String("s3Url", s3url))
		// ユーザ登録
		createUser := model.MstUser{}
		createUser.Email = u.Email
		createUser.Username = u.Username
		createUser.Password = toHashPassword(u.Password)
		createUser.Photo = s3url
		createUser.S3Key = fileId.String()
		// トランザクション開始
		tx := db.Begin()
		defer tx.Close()
		if err := tx.Create(&createUser).Error; err != nil {
			logger.Log.Info("ユーザ登録失敗")
			logger.Log.Info("ユーザ登録API終了")
			return context.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "ユーザテーブルへ登録できませんでした",
			})
		}
		// コミット
		tx.Commit()
		logger.Log.Info("ユーザ登録API終了")
		return context.String(http.StatusOK, "")
	}
}

// パスワードハッシュ化
func toHashPassword(pass string) string {
	converted, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(converted)
}

// ハッシュ化されたパスワードとの比較
func compareHashedPassword(hash string, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err == nil {
		return true
	}
	return false
}
