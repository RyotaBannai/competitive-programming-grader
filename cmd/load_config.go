package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/k0kubun/pp"
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

func Debug(c interface{}) {
	pp.Println(c)
}

func LoadConf() (c Config) {
	v := viper.New()
	dir, file := filepath.Split(fmt.Sprintf("%v", os.Getenv(APP_CONF_EV)))
	if dir == "" && file == NIL { // 環境変数が設定されていない
		file = DEFAULT_CONF_FILENAME
		dir = DEFAULT_CONF_DIR
	}
	v.SetConfigName(file)
	v.AddConfigPath(dir)

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("couldn't load config: %s", err)
		os.Exit(1)
	}

	if err := v.Unmarshal(&c); err != nil {
		fmt.Printf("couldn't read config: %s", err)
	}

	return
}
