package main

import "github.com/spf13/viper"

func readConfig(filename string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, val := range defaults {
		v.SetDefault(key, val)
	}

	v.SetConfigName("") //this will default to config.env
	v.AddConfigPath(".")
	v.SetConfigType("env")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}
