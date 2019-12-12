package cmd

import (
	"github.com/pojntfx/godibs/pkg/utils"
	"github.com/spf13/viper"
)

var ConfigContent utils.Dibs

func ReadConfig(path, file string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName(file)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&ConfigContent); err != nil {
		return err
	}

	return nil
}
