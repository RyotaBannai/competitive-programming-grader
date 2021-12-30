package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"RyotaBannai/competitive-programming-grader/internal/consts"
	"RyotaBannai/competitive-programming-grader/internal/pkg/appio"

	"github.com/spf13/viper"
)

type LangConfig struct {
	Lang string `mapstructure:"language"`
}

type CompileConfig struct {
	Compile   bool
	Command   string
	OutputDir string `mapstructure:"output_dir"`
}

type TestConfig struct {
	TestfileDir string `mapstructure:"testfile_dir"`
}

type Config struct {
	Compile CompileConfig
	Lang    LangConfig
	Test    TestConfig
}

// accept format:
// $HOME/cpg_conf
// $HOME/cpg_conf.toml
// ./cpg_conf.toml
// If there is no `CPG_CONF_PATH` env varaible, then
// looking for $HOME/.config/cpg/cpg_conf.toml, and then
// `cpg_conf.toml` from current directry
func LoadConf() (c Config) {
	var (
		file string
		dir  string
	)
	envFilepath := os.Getenv(consts.APP_CONF_EV)
	if result, _ := appio.Exists(envFilepath); result {
		dir, file = filepath.Split(envFilepath)
	} else { // 環境変数が指定されていない or 指定されているが存在しない
		homeDir, err := os.UserHomeDir()
		confFilepath := filepath.Join(homeDir, ".config", consts.APP_NAME, consts.APP_CONF_FILENAME)
		if envFilepath != "" {
			fmt.Printf("env variable [%v] was found, but [%v] doesn't exist\nlooking to [%v]\n", consts.APP_CONF_EV, envFilepath, confFilepath)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if result, _ := appio.Exists(confFilepath); result { // $HOME/.config/cpg/cpg_conf.toml, for both darwin and linux
			dir, file = filepath.Split(confFilepath)
		} else { // デフォルトの config フォルダ配下にもない
			file = consts.APP_CONF_FILENAME
			dir = consts.DEFAULT_CONF_DIR
		}
	}

	name := strings.TrimSuffix(file, filepath.Ext(file))
	v := viper.New()
	v.SetConfigName(name)
	v.AddConfigPath(dir)

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	if err := v.Unmarshal(&c); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	return
}
