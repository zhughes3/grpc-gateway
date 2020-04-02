package main

import (
	"net"

	"github.com/spf13/viper"
)

type serverConfig struct {
	rpcHost, rpcPort, httpHost, httpPort string
}

type dbConfig struct {
	name, host, port, user, password string
}

func main() {
	config_defaults := map[string]interface{}{
		"rpc_host":  "localhost",
		"http_host": "localhost",
		"db_host":   "localhost",
	}
	config, err := readConfig("config.env", config_defaults)
	if err != nil {
		panic(err)
	}
	sCfg := getServerConfig(config)
	dbCfg := getDBConfig(config)
	conn, err := net.Listen("tcp", ":"+sCfg.httpPort)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	db := NewDB(dbCfg)
	defer db.Close()

	s := NewServer(sCfg, conn, db)

	s.Start()
}

func getServerConfig(cfg *viper.Viper) *serverConfig {
	return &serverConfig{
		rpcHost:  cfg.GetString("rpc_host"),
		rpcPort:  cfg.GetString("rpc_port"),
		httpHost: cfg.GetString("http_host"),
		httpPort: cfg.GetString("http_port"),
	}
}

func getDBConfig(cfg *viper.Viper) *dbConfig {
	return &dbConfig{
		name:     cfg.GetString("db_name"),
		host:     cfg.GetString("db_host"),
		port:     cfg.GetString("db_port"),
		user:     cfg.GetString("db_user"),
		password: cfg.GetString("db_password"),
	}
}
