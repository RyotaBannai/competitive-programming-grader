package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"RyotaBannai/competitive-programming-grader/internal/pkg/appio"
)

var createCmd = &cobra.Command{
	Use:   "m",
	Short: "Make test file for Problem X i.g. cpg m -p d",
	Run: func(cmd *cobra.Command, args []string) {
		var p interface{}
		if viper.Get("p") != nil {
			p = viper.Get("p")
		} else if viper.Get("prob") != "" {
			p = viper.Get("prob")
		} else {
			// finish.
			fmt.Println("Please set problem")
			return
		}

		pStr := fmt.Sprintf("%v", p)
		if result, _ := appio.Exists(pStr); !result {
			fmt.Printf("file path [%v] is not found.\n", pStr)
			return
		}

		handle, err := appio.FileOpener(pStr)
		defer appio.FileCloser(handle)
		if err != nil {
			return
		}

		comment := appio.PeekComment(handle)
		dirSpec, peekResult := appio.CheckDirSpecAnnotation(comment)

		// create test file dir if it doesn't exist
		testDir := conf.Test.TestfileDir
		if result, _ := appio.ExistsOrCreateFolder(testDir); !result {
			fmt.Println(err)
			return
		}

		var dirSpecPath string
		if peekResult {
			// found @cpg_dirspec annotation
			dirSpecPath = filepath.Join(testDir, dirSpec)
		} else {
			// use filename as test file dir i.g. a.cpp -> a
			splited := strings.Split(pStr, ".")
			dirSpecPath = filepath.Join(testDir, splited[0])
		}

		inDir := filepath.Join(dirSpecPath, "in")
		outDir := filepath.Join(dirSpecPath, "out")
		infiles, _ := ioutil.ReadDir(inDir)
		insamples := len(infiles)
		outfiles, _ := ioutil.ReadDir(outDir)
		outsamples := len(outfiles)
		if insamples != outsamples {
			fmt.Println("mismatch number of in and out cases")
			return
		}

		if result, _ := appio.ExistsOrCreateFolder(dirSpecPath); !result {
			fmt.Println(err)
			return
		}
		if result, _ := appio.ExistsOrCreateFolder(inDir); !result {
			fmt.Println(err)
			return
		}
		if result, _ := appio.ExistsOrCreateFolder(outDir); !result {
			fmt.Println(err)
			return
		}

		// ready to accept inputs
		fmt.Println("[input]")
		in, err := appio.ReadLines()
		if err != nil {
			fmt.Println("goodbye.")
			return
		}
		fmt.Println("[output]")
		out, err := appio.ReadLines()
		if err != nil {
			fmt.Println("goodbye.")
			return
		}

		// ready to write
		inf := filepath.Join(inDir, "sample_in_"+fmt.Sprintf("%v", insamples)+".txt ")
		outf := filepath.Join(outDir, "sample_out_"+fmt.Sprintf("%v", outsamples)+".txt ")

		for _, p := range []struct {
			path     string
			contents []string
			name     string
		}{{inf, in, "in"}, {outf, out, "out"}} {
			if result, err := appio.Exists(p.path); result || err != nil {
				fmt.Printf("file [%v] already exists", p.path)
				return
			} else {
				handle, err := os.Create(p.path)
				defer appio.FileCloser(handle)
				if err != nil {
					fmt.Printf("failed to create [%v] test case\n", p.name)
					return
				}
				if _, err := handle.WriteString(strings.Join(p.contents[:], "\n")); err != nil {
					fmt.Printf("failed to create [%v] test case\n", p.name)
					return
				}
			}
		}

		fmt.Printf("Finished to create test case. \nin[ %v]\nout[ %v]\n", inf, outf)
	},
}
