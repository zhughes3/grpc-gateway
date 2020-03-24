package main

import (
	"net"
)

type serverConfig struct {
	rpcHost, rpcPort, httpHost, httpPort string
}

type dbConfig struct {
	name, host, port, user, password string
}

func main() {
	config_defaults := map[string]interface{}{
		"rpc_host": "localhost",
		"http_host": "localhost",
		"db_host": "localhost",
	}
	config, err := readConfig("config.env", config_defaults)
	if err != nil {
		panic(err)
	}
	sCfg := &serverConfig{
		rpcHost:     config.GetString("rpc_host"),
		rpcPort:  config.GetString("rpc_port"),
		httpHost: config.GetString("http_host"),
		httpPort: config.GetString("http_port"),
	}
	dbCfg := &dbConfig{
		name:     config.GetString("db_name"),
		host:     config.GetString("db_host"),
		port:     config.GetString("db_port"),
		user:     config.GetString("db_user"),
		password: config.GetString("db_password"),
	}
	conn, err := net.Listen("tcp", ":" + sCfg.httpPort)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	db := NewDB(dbCfg)
	defer db.Close()

	s := NewServer(sCfg, conn, db)

	s.Start()
}



