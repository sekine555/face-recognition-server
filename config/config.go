package config

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"gopkg.in/ini.v1"
	"os"
)

type ConfigList struct {
	DbDriverName    string
	DbName          string
	DbUserName      string
	DbUserPassword  string
	DbHost          string
	DbPort          string
	Secret          string
	LoggerFilePath  string
	LoggerLevel     string
	Region          string
	Bucket          string
	AccessKeyId     string
	SecretAccessKey string
}

var Config ConfigList

func init()  {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}
	// 実行環境取得
	goEnv := cfg.Section("env").Key("go_env").MustString("development")
	fmt.Println("実行環境：", goEnv)
	// 環境変数設定
	if goEnv == "development" {
		Config = ConfigList{
			DbDriverName:    cfg.Section("dev").Key("db_driver_name").String(),
			DbName:          cfg.Section("dev").Key("db_name").String(),
			DbUserName:      cfg.Section("dev").Key("db_user_name").String(),
			DbUserPassword:  cfg.Section("dev").Key("db_user_password").String(),
			DbHost:          cfg.Section("dev").Key("db_host").String(),
			DbPort:          cfg.Section("dev").Key("db_port").String(),
			Secret:          cfg.Section("key").Key("secret").String(),
			LoggerFilePath:  cfg.Section("log").Key("logger_file_path").String(),
			LoggerLevel:     cfg.Section("log").Key("logger_level").MustString("info"),
			Region:          cfg.Section("aws").Key("region").String(),
			Bucket:          cfg.Section("aws").Key("bucket").String(),
			AccessKeyId:     cfg.Section("aws").Key("access_key_id").String(),
			SecretAccessKey: cfg.Section("aws").Key("secret_access_key").String(),
		}
	} else {
		Config = ConfigList{
			DbDriverName:    cfg.Section("prd").Key("db_driver_name").String(),
			DbName:          cfg.Section("prd").Key("db_name").String(),
			DbUserName:      cfg.Section("prd").Key("db_user_name").String(),
			DbUserPassword:  cfg.Section("prd").Key("db_user_password").String(),
			DbHost:          cfg.Section("prd").Key("db_host").String(),
			DbPort:          cfg.Section("prd").Key("db_port").String(),
			Secret:          cfg.Section("key").Key("secret").String(),
			LoggerFilePath:  cfg.Section("log").Key("logger_file_path").String(),
			LoggerLevel:     cfg.Section("log").Key("logger_level").MustString("info"),
			Region:          cfg.Section("aws").Key("region").String(),
			Bucket:          cfg.Section("aws").Key("bucket").String(),
			AccessKeyId:     cfg.Section("aws").Key("access_key_id").String(),
			SecretAccessKey: cfg.Section("aws").Key("secret_access_key").String(),
		}
	}
}
