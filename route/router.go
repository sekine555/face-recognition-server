package route

import (
	"face-recognition/api"
	"face-recognition/config"
	"face-recognition/logger"
	"fmt"
	"github.com/labstack/echo"
	echoMw "github.com/labstack/echo/middleware"
	"go.uber.org/zap"
)

func bodyDumpHandler(c echo.Context, reqBody, resBody []byte) {
	logger.Log.Info("Request Body", zap.String("パラメータ", string(reqBody)))
	fmt.Printf("Request Body: %v\n", string(reqBody))
}

func Init() *echo.Echo {
	// インスタンス生成
	e := echo.New()
	// アプリケーションのどこかで予期せずにpanicを起こしてしまっても、サーバは落とさずにエラーレスポンスを返せるようにリカバリーする
	e.Use(echoMw.Recover())
	// アクセスログ出力
	e.Use(echoMw.Logger())
	// リクエストボディの値をログ出力
	e.Use(echoMw.BodyDump(bodyDumpHandler))
	// ルーティング
	// バージョン管理用にパスを束ねる
	v1 := e.Group("/api/v1")
	{
		v1.POST("/users/login", api.PostLogin())
		v1.POST("/users/register", api.PostUser())
		// 認証ミドルウェア設定
		v1.Use(echoMw.JWT([]byte(config.Config.Secret)))
		// ここより下のエンドポイントはJWT認証必須
		v1.GET("/users", api.GetUser())
		v1.GET("/qr-token", api.GetQrToken())
		v1.POST("/face-recognition", api.PostFaceRecognition())
	}
	// 生成したechoを返却
	return e
}
