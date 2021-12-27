package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "cpg",
	Short: "Competitive Programming Grader for automating coding-build-testing loop. ",
	Long: `Competitive Programming Grader for automating coding-build-testing loop. 
- Created and maintained by RyotaBannai`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Create test file for Problem X",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[input]")
		fmt.Println("[output]")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().StringP("create", "c", "", "Filename to create test files. cpg -c d")
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.BindPFlags(rootCmd.PersistentFlags())
	Debug(viper.Get("create"))
	// or bind to private variable
	// var Source string
	// rootCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
}
