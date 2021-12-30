package cli

import (
	"RyotaBannai/competitive-programming-grader/internal/consts"
	"RyotaBannai/competitive-programming-grader/internal/pkg/appio"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mattn/go-shellwords"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

var runTestCmd = &cobra.Command{
	Use:   "run",
	Short: "Run tests for a problem X i.g. cpg run -p d.cpp",
	Run: func(cmd *cobra.Command, args []string) {
		p, err := takeProb()
		if err != nil {
			fmt.Println("Please choose problem and set [p] flag")
			return
		}

		if result, _ := appio.Exists(p); !result {
			fmt.Printf("file path [%v] is not found.\n", p)
			return
		}

		handle, err := appio.FileOpener(p)
		defer appio.FileCloser(handle)
		if err != nil {
			return
		}

		comment := appio.PeekComment(handle)
		dirSpec, peekResult := appio.CheckDirSpecAnnotation(comment)
		testDir := conf.Test.TestfileDir

		var dirSpecPath string
		if peekResult {
			// found @cpg_dirspec annotation
			dirSpecPath = filepath.Join(testDir, dirSpec)
		} else {
			// use filename as test file dir i.g. a.cpp -> a
			// first, remove all dir path
			splited := strings.Split(filepath.Base(p), ".")
			dirSpecPath = filepath.Join(testDir, splited[0])
		}

		inDir := filepath.Join(dirSpecPath, "in")
		outDir := filepath.Join(dirSpecPath, "out")

		if result, _ := appio.Exists(inDir); !result {
			fmt.Printf("inDir [%v] is not found.\n", inDir)
			return
		}
		if result, _ := appio.Exists(outDir); !result {
			fmt.Printf("outDir [%v] is not found.\n", inDir)
			return
		}

		infiles, _ := ioutil.ReadDir(inDir)
		insamples := len(infiles)
		outfiles, _ := ioutil.ReadDir(outDir)
		outsamples := len(outfiles)
		if insamples != outsamples {
			fmt.Println("mismatch number of in and out cases")
			return
		}

		command := strings.Replace(conf.Compile.Command, consts.COMMAND_FILENAME_PLACEHOLDER, p, 1)
		c, err := shellwords.Parse(command) // 文字列をコマンド、オプション単位でスライス化
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if conf.Compile.Compile {
			// expect compiled excutable filepath is `./a.out`
			if _, err := exec.Command(c[0], c[1:]...).Output(); err != nil {
				fmt.Printf("failed to compile [%v]\n", command)
				return
			}
		}

		byName := func(s []fs.FileInfo) func(int, int) bool {
			return func(i, j int) bool {
				return s[i].Name() < s[j].Name()
			}
		}
		sort.Slice(infiles, byName(infiles))
		sort.Slice(outfiles, byName(outfiles))

		for i := 0; i < len(infiles); i++ {
			iTestCasePath := filepath.Join(inDir, infiles[i].Name())
			handle1, err := appio.FileOpener(iTestCasePath)
			defer appio.FileCloser(handle1)
			if err != nil {
				return
			}
			oTestCasePath := filepath.Join(outDir, outfiles[i].Name())
			handle2, err := appio.FileOpener(oTestCasePath)
			defer appio.FileCloser(handle2)
			if err != nil {
				return
			}
			ifc := appio.ReadFileContents(handle1)
			ofc := appio.ReadFileContents(handle2)

			var cmd *exec.Cmd
			if conf.Compile.Compile { // use executable
				cmd = exec.Command("./a.out")
			} else { // use given command for script files
				cmd = exec.Command(c[0], c[1:]...)
			}

			stdin, _ := cmd.StdinPipe()
			io.WriteString(stdin, strings.Join(ifc.Contents[:], "\n"))
			stdin.Close()
			out, err := cmd.Output()
			if err != nil {
				fmt.Println("command execution failed.")
				fmt.Println(err.Error())
				return
			}

			expect := strings.TrimSpace(strings.Join(ofc.Contents[:], "\n"))
			actual := strings.TrimSpace(string(out))

			if expect == actual {
				fmt.Printf("Ok [%s]\n", iTestCasePath)
			} else {
				fmt.Printf("Failed [%v]\n", iTestCasePath)
				dmp := diffmatchpatch.New()
				diffs := dmp.DiffMain(expect, actual, false)
				fmt.Println(dmp.DiffPrettyText(diffs))

				break
			}
		}
		fmt.Println("Done.")
	},
}
