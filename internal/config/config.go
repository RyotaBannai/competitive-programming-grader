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
	Compile        bool
	Command        string
	OutputDir      string `mapstructure:"output_dir"`
	OutputFilename string `mapstructure:"output_filename"`
}

type TestConfig struct {
	TestfileDir   string `mapstructure:"testfile_dir"`
	ParseComments bool   `mapstructure:"parse_comments"`
}

type Config struct {
	Compile CompileConfig
	Lang    LangConfig
	Test    TestConfig
}

// accept format:
// $HOME/cpg_conf
// $HOME/cpg_conf.toml
// If there is no `CPG_CONF_PATH` env varaible, then looking for `cpg_conf.toml` in current directry
func LoadConf() (c Config) {
	var (
		file string
		dir  string
	)
	envfilepath := os.Getenv(consts.APP_CONF_EV)
	if result, _ := appio.Exists(envfilepath); !result { // 環境変数が指定されていない時も空で見つからない
		file = consts.DEFAULT_CONF_FILENAME
		dir = consts.DEFAULT_CONF_DIR
	} else {
		dir, file = filepath.Split(envfilepath)
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
