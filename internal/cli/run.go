package cli

import (
	"RyotaBannai/competitive-programming-grader/internal/consts"
	"RyotaBannai/competitive-programming-grader/internal/pkg/appio"
	"RyotaBannai/competitive-programming-grader/internal/pkg/misc"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gookit/color"
	"github.com/mattn/go-shellwords"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

// e.g. fs.FileInfo
type HasName interface {
	Name() string
}

func sortFilebyName[T HasName](s []T) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].Name() < s[j].Name()
	})
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
		outfiles, _ := ioutil.ReadDir(outDir)
		if len(infiles) != len(outfiles) {
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

		sortFilebyName(infiles)
		sortFilebyName(outfiles)
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

			var cmd *exec.Cmd
			if conf.Compile.Compile { // use executable
				cmd = exec.Command(appio.Join(conf.Compile.OutputDir, consts.EXECUTABLE_NAME))
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
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(expect, actual, false)
			misc.Debug(dmp.DiffText1(diffs))
			misc.Debug(dmp.DiffText2(diffs))

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
