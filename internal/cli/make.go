package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"RyotaBannai/competitive-programming-grader/internal/pkg/appio"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:     "make",
	Short:   "Make a test case for a problem X",
	Example: "  cpg make -p d.cpp",
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
			// first, remove all dir path
			splited := strings.Split(filepath.Base(p), ".")
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
		color.Notice.Println("[input]")
		in, err := appio.ReadLines()
		if err != nil {
			color.HiWhite.Println("goodbye.")
			return
		}
		color.Notice.Println("[output]")
		out, err := appio.ReadLines()
		if err != nil {
			color.HiWhite.Println("goodbye.")
			return
		}

		location := time.FixedZone("Asia/Tokyo", 9*60*60) // will be fixed from env variable
		now := func() time.Time { return time.Now().In(location) }
		suffix := now().Format("2006-01-02_15:04:05") // or time.RFC3339

		// ready to write
		inf := filepath.Join(inDir, "sample_in_"+fmt.Sprintf("%v", insamples)+"_"+suffix+".txt ")
		outf := filepath.Join(outDir, "sample_out_"+fmt.Sprintf("%v", outsamples)+"_"+suffix+".txt ")

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

		arrow := color.Bold.Sprint("→")
		color.New(color.HiWhite, color.Bold).Print("\nCreated a test case in the following paths: \n")
		color.HiWhite.Printf("[input]  "+arrow+" %v\n[outout] "+arrow+" %v\n", strings.TrimSpace(inf), strings.TrimSpace(outf))
	},
}
