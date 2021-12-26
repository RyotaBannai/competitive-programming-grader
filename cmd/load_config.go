package cmd

import (
	"fmt"
	"os"

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

func Debug(c Config) {
	pp.Println(c)
}

func LoadConf() (c Config) {
	v := viper.New()
	v.SetConfigName("cpg_conf") // TODO: CPG_CONF_PATH をセットする（なければデフォルトでカレント）
	v.AddConfigPath(".")        // TODO: CPG_CONF_PATH をパースして filename と dir に分けてセットする

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("couldn't load config: %s", err)
		os.Exit(1)
	}

	if err := v.Unmarshal(&c); err != nil {
		fmt.Printf("couldn't read config: %s", err)
	}
	// Debug(c)

	return
}
