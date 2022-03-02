package cli

import (
	"fmt"
	"strings"

	"RyotaBannai/competitive-programming-grader/internal/consts"
	"RyotaBannai/competitive-programming-grader/internal/pkg/appio"

	"github.com/gocolly/colly"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:     "fetch",
	Short:   "fetch all test cases for the contest X on AtCoder",
	Example: "  cpg fetch -c ",
	Run: func(cmd *cobra.Command, args []string) {
		contestId, err := getContest()
		if err != nil {
			fmt.Println("Please set [c] flag")
			return
		}

		contestTaskPage := strings.Replace(consts.URLS.AT_CODER_TASKS, consts.AT_CODER_ID_PLACEHOLDER, contestId, 1)
		controller := colly.NewCollector()

		// Find tasks
		controller.OnHTML("tbody tr td:nth-of-type(1) a", func(e *colly.HTMLElement) {
			fmt.Println("First column of a table row:", e.Text)
			taskPage := contestTaskPage + "/" + contestId + "_" + e.Text
			e.Request.Visit(taskPage)
		})

		// Find in/out
		controller.OnHTML("div.part section", func(e *colly.HTMLElement) {
			// for _, elem := range e.DOM.Children().Nodes {
			// 	fmt.Println("pre id:", elem.text)
			// }

			e.ForEach("pre", func(i int, e *colly.HTMLElement) {
				a, _ := e.DOM.Html()
				fmt.Println(a)
			})
		})

		controller.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
		})

		controller.Visit(contestTaskPage)

		// create test file dir if it doesn't exist
		testDir := conf.Test.TestfileDir
		if result, _ := appio.ExistsOrCreateFolder(testDir); !result {
			fmt.Println(err)
			return
		}

		// dirSpecPath := filepath.Join(testDir, dirSpec)

		// inDir := filepath.Join(dirSpecPath, "in")
		// outDir := filepath.Join(dirSpecPath, "out")
		// infiles, _ := ioutil.ReadDir(inDir)
		// insamples := len(infiles)
		// outfiles, _ := ioutil.ReadDir(outDir)
		// outsamples := len(outfiles)
		// if insamples != outsamples {
		// 	fmt.Println("mismatch number of in and out cases")
		// 	return
		// }

		// if result, _ := appio.ExistsOrCreateFolder(dirSpecPath); !result {
		// 	fmt.Println(err)
		// 	return
		// }
		// if result, _ := appio.ExistsOrCreateFolder(inDir); !result {
		// 	fmt.Println(err)
		// 	return
		// }
		// if result, _ := appio.ExistsOrCreateFolder(outDir); !result {
		// 	fmt.Println(err)
		// 	return
		// }

		// ready to accept inputs
		// color.Notice.Println("[input]")
		// in, err := appio.ReadLines()
		// if err != nil {
		// 	color.HiWhite.Println("goodbye.")
		// 	return
		// }
		// color.Notice.Println("[output]")
		// out, err := appio.ReadLines()
		// if err != nil {
		// 	color.HiWhite.Println("goodbye.")
		// 	return
		// }

		// location := time.FixedZone("Asia/Tokyo", 9*60*60) // will be fixed from env variable
		// now := func() time.Time { return time.Now().In(location) }
		// suffix := now().Format("2006-01-02_15:04:05") // or time.RFC3339

		// // ready to write
		// inf := filepath.Join(inDir, "sample_in_"+fmt.Sprintf("%v", insamples)+"_"+suffix+".txt ")
		// outf := filepath.Join(outDir, "sample_out_"+fmt.Sprintf("%v", outsamples)+"_"+suffix+".txt ")

		// for _, p := range []struct {
		// 	path     string
		// 	contents []string
		// 	name     string
		// }{{inf, in, "in"}, {outf, out, "out"}} {
		// 	if result, err := appio.Exists(p.path); result || err != nil {
		// 		fmt.Printf("file [%v] already exists", p.path)
		// 		return
		// 	} else {
		// 		handle, err := os.Create(p.path)
		// 		defer appio.FileCloser(handle)
		// 		if err != nil {
		// 			fmt.Printf("failed to create [%v] test case\n", p.name)
		// 			return
		// 		}
		// 		if _, err := handle.WriteString(strings.Join(p.contents[:], "\n")); err != nil {
		// 			fmt.Printf("failed to create [%v] test case\n", p.name)
		// 			return
		// 		}
		// 	}
		// }

		// arrow := color.Bold.Sprint("→")
		// color.New(color.HiWhite, color.Bold).Print("\nCreated a test case in the following paths: \n")
		// color.HiWhite.Printf("[input]  "+arrow+" %v\n[outout] "+arrow+" %v\n", strings.TrimSpace(inf), strings.TrimSpace(outf))
	},
}
