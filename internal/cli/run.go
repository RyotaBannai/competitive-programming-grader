package cli

import (
	"RyotaBannai/competitive-programming-grader/internal/consts"
	"RyotaBannai/competitive-programming-grader/internal/pkg/appio"
	"RyotaBannai/competitive-programming-grader/internal/pkg/lib"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/mattn/go-shellwords"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

type CommandFields struct {
	args          []string // コマンド実行時の引数
	filename      string   // 実行対象のファイル名
	command       string   // ビルド時のコマンド
	excutableName string   // ビルド後の artifact 名
}

type Command interface {
	Fields() CommandFields
	Exec() (*exec.Cmd, error)
	Build() error
}

/**
* C++
* and other compiled languages
*
* expect `./a.out` as compiled excutable filepath (or artifact)
 */
type CppCommand struct {
	CommandFields
}

func (cmd *CppCommand) Build() error {
	cmd.command = strings.Replace(conf.Compile.Command, consts.COMMAND_FILENAME_PLACEHOLDER, cmd.filename, 1)
	// 文字列をコマンド、オプション単位でスライス化
	if words, err := shellwords.Parse(cmd.command); err != nil {
		return err
	} else {
		exec.Command(words[0], words[1:]...).Output()
		return nil
	}
}

func (cmd *CppCommand) Fields() CommandFields {
	return cmd.CommandFields
}

func (cmd *CppCommand) Exec() (*exec.Cmd, error) {
	return exec.Command(appio.Join(conf.Compile.OutputDir, cmd.excutableName)), nil
}

/**
* Python
* and other script languages
 */
type PythonCommand struct {
	CommandFields
}

func (cmd *PythonCommand) Build() error {
	/** Nothing to build */
	return nil
}

func (cmd *PythonCommand) Fields() CommandFields {
	return cmd.CommandFields
}

func (cmd *PythonCommand) Exec() (*exec.Cmd, error) {
	cmd.command = strings.Replace(conf.Compile.Command, consts.COMMAND_FILENAME_PLACEHOLDER, cmd.filename, 1)
	// 文字列をコマンド、オプション単位でスライス化
	if words, err := shellwords.Parse(cmd.command); err != nil {
		return nil, err
	} else {
		return exec.Command(words[0], words[1:]...), nil
	}
}

/**
* Rust
*
* Rust doesn't allow user to set output artifact name
 */
type RustCommand struct {
	CommandFields
}

func (cmd *RustCommand) Build() error {
	// Expecting using a bin name..
	cmd.command = strings.Replace(conf.Compile.Command, consts.COMMAND_FILENAME_PLACEHOLDER, cmd.excutableName, 1)
	// 文字列をコマンド、オプション単位でスライス化
	if words, err := shellwords.Parse(cmd.command); err != nil {
		return err
	} else {
		exec.Command(words[0], words[1:]...).Output()
		return nil
	}
}

func (cmd *RustCommand) Fields() CommandFields {
	return cmd.CommandFields
}

func (cmd *RustCommand) Exec() (*exec.Cmd, error) {
	return exec.Command(strings.Join([]string{conf.Compile.OutputDir, cmd.excutableName}, "")), nil
}

var runTestCmd = &cobra.Command{
	Use:     "run",
	Short:   "Run tests for a problem X",
	Example: "  cpg run -p d.cpp",
	Run: func(cmd *cobra.Command, args []string) {
		p, err := getProb()
		if err != nil {
			fmt.Println("Please set [p] flag")
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

		dirSpec, peekResult := appio.CheckDirSpecAnnotation(handle)
		testDir := conf.Test.TestfileDir

		// use filename as test file dir i.g. a.cpp -> a
		// first, remove all dir path
		targetName := strings.Split(filepath.Base(p), ".")[0]

		var dirSpecPath string
		if peekResult {
			// found @cpg_dirspec annotation
			dirSpecPath = filepath.Join(testDir, dirSpec)
		} else {
			dirSpecPath = filepath.Join(testDir, targetName)
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
		outfiles, _ := ioutil.ReadDir(outDir)
		if len(infiles) != len(outfiles) {
			fmt.Println("mismatch number of in and out cases")
			return
		}

		var command Command
		lang := strings.ToLower(conf.Lang.Lang)
		if lib.Contains(consts.RUN_TYPES.COMPILE, lang) {
			command = &CppCommand{
				CommandFields{
					args:          args,
					filename:      p,
					excutableName: consts.EXECUTABLE_NAME,
				},
			}
		} else if lib.Contains(consts.RUN_TYPES.SCRIPT, lang) {
			command = &PythonCommand{
				CommandFields{
					args:     args,
					filename: p,
				},
			}
		} else if lib.Contains(consts.RUN_TYPES.RUST, lang) {
			command = &RustCommand{
				CommandFields{
					args:          args,
					filename:      p,
					excutableName: targetName,
				},
			}
		} else {
			fmt.Printf("language is invalid or unsupported yet for [%v]\n", conf.Lang.Lang)
			return
		}

		if err := command.Build(); err != nil {
			fmt.Println(err.Error())
			fmt.Printf("failed to compile [%v]\n", command.Fields().command)
			return
		}

		lib.SortFilebyName(infiles)
		lib.SortFilebyName(outfiles)
		testIgnoreCnt := 0
		testAllowCnt := 0
		testIgnoreMap := map[int]bool{}
		testAllowMap := map[int]bool{}
		nTestCases := len(infiles)
		for i := 0; i < nTestCases; i++ {
			iTestCasePath := filepath.Join(inDir, infiles[i].Name())
			handle, err := appio.FileOpener(iTestCasePath)
			defer appio.FileCloser(handle)
			if err != nil {
				return
			}
			testAnnotationMap := appio.CheckTestAnnotations(handle)
			testIgnoreMap[i] = testAnnotationMap["ignore"]
			testAllowMap[i] = testAnnotationMap["allow"]
			if testIgnoreMap[i] {
				testIgnoreCnt++
			}
			if testAllowMap[i] {
				testAllowCnt++
			}
		}

		// Test Message
		message := fmt.Sprintf(" Found %v", nTestCases)
		if testAllowCnt > 0 {
			message += fmt.Sprintf("/ Skip %v\n", nTestCases-testAllowCnt)
		} else if testIgnoreCnt > 0 {
			message += fmt.Sprintf("/ Skip %v\n", testIgnoreCnt)
		} else {
			message += "\n"
		}

		color.Bold.Println("\nTest Cases:")
		color.Println(message)
		for i := 0; i < nTestCases; i++ {
			iTestCasePath := filepath.Join(inDir, infiles[i].Name())
			displaySKipMessage := func(path string) {
				color.New(color.Gray, color.BgDefault, color.Bold).Print(" SKIP ")
				color.Gray.Print(" -")
				color.Gray.Printf(" %v\n\n", path)
			}

			// check Annotations. prioritize Allow over Ignore Annotion
			if testAllowCnt > 0 {
				if !testAllowMap[i] {
					displaySKipMessage(iTestCasePath)
					continue
				}
			} else {
				if testIgnoreMap[i] {
					displaySKipMessage(iTestCasePath)
					continue
				}
			}

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

			cmd, err := command.Exec()
			if err != nil {
				fmt.Println("command execution failed.")
				fmt.Println(err.Error())
				return
			}

			stdin, _ := cmd.StdinPipe()
			io.WriteString(stdin, strings.Join(ifc.Contents[:], "\n"))
			stdin.Close()
			out, err := cmd.Output()
			if err != nil {
				fmt.Println("command run failed.")
				fmt.Println(err.Error())
				return
			}

			// clean outputs up
			// e.g. from "1 4 7 10 \n2 5 8 11 \n3 6 9 12" to "1 4 7 10\n2 5 8 11\n3 6 9 12"
			var cleaned []string
			for _, l := range strings.Split(string(out), "\n") {
				cleaned = append(cleaned, strings.TrimSpace(l))
			}
			expect := strings.Join(ofc.Contents, "\n")               // ReadFileContents does Trim when reading file content
			actual := strings.TrimSpace(strings.Join(cleaned, "\n")) // trim last newline

			if expect == actual { // show success message
				color.New(color.Gray, color.BgGreen, color.Bold).Print(" PASS ")
				color.Green.Print(" ✔")
				color.Gray.Printf(" %v\n\n", iTestCasePath)
			} else { // show failed message
				color.Error.Print(" Fail ")
				color.Red.Print(" ×")
				color.HiWhite.Printf(" %v\n\n", iTestCasePath)
				// show expect and actual
				ttl := color.New(color.HiWhite, color.Bold)
				body := color.HiWhite
				ttl.Println("Expect:")
				body.Printf("%v\n\n", expect)
				ttl.Println("Actual:")
				body.Printf("%v\n\n", actual)
				// show diff
				ttl.Println("Diff:")
				dmp := diffmatchpatch.New()
				diffs := dmp.DiffMain(expect, actual, false)
				fmt.Print(dmp.DiffPrettyText(diffs) + "\n\n")
				return
			}
		}
		fmt.Println("Done.")
	},
}
