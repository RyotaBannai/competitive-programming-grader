package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"RyotaBannai/competitive-programming-grader/internal/config"
)

var rootCmd = &cobra.Command{
	Use:   "cpg",
	Short: "Competitive Programming Grader for automating coding-build-testing loop. ",
	Long: `Competitive Programming Grader for automating coding-build-testing loop. 
- Created and maintained by RyotaBannai`,
	Run: func(cmd *cobra.Command, args []string) {
		// pass
	},
}

var (
	conf = config.LoadConf()
)

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(runTestCmd)
	rootCmd.PersistentFlags().StringP("prob", "p", "", "Set problem")
	viper.BindPFlags(rootCmd.PersistentFlags())
	// or bind to private variable
	// var Source string
	// rootCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
