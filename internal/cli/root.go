package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"RyotaBannai/competitive-programming-grader/internal/config"
)

var rootCmd = &cobra.Command{
	Use:   "cpg",
	Short: "Competitive Programming Grader",
	Long:  "Competitive Programming Grader for automating coding-build-testing loop",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

var (
	conf = config.LoadConf()
)

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(runTestCmd)
	rootCmd.AddCommand(fetchCmd)
	rootCmd.PersistentFlags().StringP("prob", "p", "", "Set problem")
	rootCmd.PersistentFlags().StringP("contest", "c", "", "Set contenst number")
	viper.BindPFlags(rootCmd.PersistentFlags())
	// or bind to private variable
	// var Source string
	// rootCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
}

func getProb() (string, error) {
	if viper.IsSet("prob") {
		return viper.GetString("prob"), nil
	} else { // finish.
		return "", errors.New("flag p is required")
	}
}

func getContest() (string, error) {
	if viper.IsSet("contest") {
		return viper.GetString("contest"), nil
	} else { // finish.
		return "", errors.New("flag c is required")
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
