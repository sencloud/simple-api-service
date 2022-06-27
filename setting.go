package main

import (
	"fmt"
	"github.com/spf13/viper"
)

//定义结构体
type MysqlServer struct {
	HostIp   string
	Port     string
	User     string
	Password string
	Db       string
}

//定义全局配置文件
var cfg = viper.New()
var MysqlServerSetting = &MysqlServer{}

func readYaml() {
	cfg.AddConfigPath("./")
	cfg.SetConfigName("config")
	cfg.SetConfigType("yaml")
	var err error
	//加载配置文件
	err = cfg.ReadInConfig()

	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		fmt.Println("找不到配置文件..")
	} else {
		fmt.Println("Loaded...")
	}

	MysqlServerSetting.HostIp = cfg.GetString("mysql.HostIp")
	MysqlServerSetting.User = cfg.GetString("mysql.User")
	MysqlServerSetting.Password = cfg.GetString("mysql.Password")
	MysqlServerSetting.Db = cfg.GetString("mysql.Db")
	MysqlServerSetting.Port = cfg.GetString("mysql.Port")
}