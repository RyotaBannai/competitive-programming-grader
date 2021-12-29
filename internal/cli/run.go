package cli

import "github.com/spf13/cobra"

var runTestCmd = &cobra.Command{
	Use:   "run",
	Short: "Run test for problem X i.g. cpg run -f d.cpp",
	Run: func(cmd *cobra.Command, args []string) {
		//
	},
}
