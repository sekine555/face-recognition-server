package main

import (
	"face-recognition/route"
)

func main() {
	// 初期設定（echoインスタンス生成などはrouteの役割）
	router := route.Init()
	// サーバ起動
	router.Logger.Fatal(router.Start(":1323"))
}