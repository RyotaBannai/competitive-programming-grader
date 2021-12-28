package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	conf = LoadConf()
)

var versionCmd = &cobra.Command{
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
		if result, _ := exists(pStr); !result {
			fmt.Printf("file path [%v] is not found.\n", pStr)
			return
		}

		handle, err := fileOpener(pStr)
		defer fileCloser(handle)
		if err != nil {
			return
		}

		comment := peekComment(handle)
		dirSpec, peekResult := checkDirSpecAnnotation(comment)

		// create test file dir if it doesn't exist
		testDir := conf.Test.TestfileDir
		if result, _ := existsOrCreateFolder(testDir); !result {
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

		if result, _ := existsOrCreateFolder(dirSpecPath); !result {
			fmt.Println(err)
			return
		}
		if result, _ := existsOrCreateFolder(inDir); !result {
			fmt.Println(err)
			return
		}
		if result, _ := existsOrCreateFolder(outDir); !result {
			fmt.Println(err)
			return
		}

		// ready to accept inputs
		fmt.Println("[input]")
		in, err := readLines()
		if err != nil {
			fmt.Println("goodbye.")
			return
		}
		fmt.Println("[output]")
		out, err := readLines()
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
			if result, err := exists(p.path); result || err != nil {
				fmt.Printf("file [%v] already exists", p.path)
				return
			} else {
				handle, err := os.Create(p.path)
				defer fileCloser(handle)
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

func readLines() ([]string, error) {
	var txt []string
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := s.Text()
		if line != "" {
			if fmt.Sprintf("%v", strings.TrimSpace(line)) == "q!" {
				// force quit input test case
				return txt, errors.New("force quit")
			} else if fmt.Sprintf("%v", s.Bytes()[0]) == "27" {
				// finish with esc key.
				// convert byte to ascii code
				break
			}
		}
		txt = append(txt, line)
		// log.Print(strconv.Quote(s.Text()))
	}

	if s.Err() != nil {
		// non-EOF error.
		log.Fatal(s.Err())
	}
	return txt, nil
}

func init() {
	rootCmd.AddCommand(versionCmd)
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
