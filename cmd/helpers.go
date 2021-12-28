package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/k0kubun/pp"
)

type FileCotents struct {
	Contents []string // file contents removed comments
	Comments []string // all comments in file `//` or multiline comment `/** ... */`
}

func Debug(c interface{}) {
	pp.Println(c)
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func existsOrCreateFolder(dir string) (bool, error) {
	if result, _ := exists(dir); !result {
		if err := os.Mkdir(dir, 0755); err != nil {
			// failed to create
			return false, err
		} else {
			// success to newly create
			return true, nil
		}
	} else {
		// already exists
		return true, nil
	}
}

func fileCloser(file *os.File) {
	if err := file.Close(); err != nil {
		// log.Fatal(fmt.Sprintf("error occurred while closing file [%v]", file.Name()))
		log.Fatal(err)
	} else {
		// log.Println(fmt.Sprintf("file [%v] closed", file.Name()))
	}
}

func fileOpener(filepath string) (file *os.File, err error) {
	if file, err = os.Open(filepath); err != nil { // file に assign するのは メモリ上の実態データ
		fmt.Println("file doesn't exist.")
		// log.Println(fmt.Sprintf("file [%v] doesn't exist", file.Name()))
		return nil, err
	} else {
		// log.Println(fmt.Sprintf("file [%v] opened", file.Name()))
		return
	}
}

func peekComment(file *os.File) string {
	if comments := readFileContentsByParsingComments(file, true).Comments; len(comments) > 0 {
		return comments[0]
	} else {
		return ""
	}
}

func checkDirSpecAnnotation(line string) (dirSpec string, b bool) {
	if strings.Contains(line, ANNOTATIONS.DIRSPEC) {
		splitted := strings.Split(line, ANNOTATIONS.DIRSPEC)
		token := strings.Fields(splitted[1])
		if len(token) > 0 {
			return token[0], true
		}
	}
	return "", false
}

func readFileContents(file *os.File) FileCotents {
	return readFileContentsByParsingComments(file, false)
}

// read text file and load test case.
// also parse comments inserted in codes. `//` or multiline comment `/** ... */`
// so that it can manage annotations embbed in codes
func readFileContentsByParsingComments(file *os.File, takeFirstComment bool) FileCotents {
	defer func() {
		if x := recover(); x != nil {
			fmt.Printf("%v", x)
		}
	}()
	causePanic := func(name string, no int) {
		panic(fmt.Sprintf("test file broken [%v] line number [%v]\n", name, no))
	}
	var contents, comments []string
	foundOne := false
	lno := 0
	multiLFlag := false
	tmpMultiLComment := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lno++
		line := scanner.Text()
		if strings.HasSuffix(line, "*/") {
			if multiLFlag {
				tmpMultiLComment += line
				comments = append(comments, tmpMultiLComment)
				tmpMultiLComment = ""
				multiLFlag = false
				foundOne = true
			} else {
				causePanic(file.Name(), lno)
			}
		} else if multiLFlag {
			tmpMultiLComment += line + "\n"
		} else if strings.HasPrefix(line, "/*") {
			if strings.HasSuffix(line, "*/") { // i.g. /** some comments */
				comments = append(comments, tmpMultiLComment)
				foundOne = true
			} else {
				multiLFlag = true
				tmpMultiLComment += line + "\n"
			}
		} else if strings.HasPrefix(line, "//") {
			comments = append(comments, line)
			foundOne = true
		} else {
			// only statement for contents
			contents = append(contents, line)
		}

		if takeFirstComment && foundOne { // take only first comment and return
			return FileCotents{Contents: contents, Comments: comments}
		}
	}

	if multiLFlag { // file doesn't close multiline comment
		causePanic(file.Name(), lno)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return FileCotents{Contents: contents, Comments: comments}
}
