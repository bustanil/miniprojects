package main

import (
	// "bytes"
	"fmt"
	"io/fs"
	// "io/ioutil"
	"os"
	"path/filepath"
	"strings"
  "regexp"

	// "github.com/yuin/goldmark"
)

var template string = `
<html>
  <head>
  </head>
  <body>
    %s
  </body>
</html>
`

type BlogEntry struct {
	// date date
	path string
}

func EndsWith(s string, endsWith string) bool {
  var l = strings.LastIndex(s, endsWith) 
	return l >= 0 && l == (len(s) - len(endsWith))
}

func main() {
	// var bs, err = ioutil.ReadFile("test.md")
	// if err != nil {
	// 	fmt.Println("error", err)
	// 	return
	// }

	args := os.Args
	input := args[1]
	// output := args[2]

	rootFilePath, _ := filepath.Abs(input)
  re := regexp.MustCompile(`(?P<Year>\d{4})/(?P<Month>\d{2})/(?P<Day>\d{2})/(?P<Filename>.*)`)
	filepath.WalkDir(input, func(path string, d fs.DirEntry, err error) error {
		subDirPath := strings.ReplaceAll(path, rootFilePath, "")
		if EndsWith(d.Name(), ".md") {
      matches := re.FindStringSubmatch(subDirPath)
      yearIndex := re.SubexpIndex("Year")
      fmt.Println(matches[yearIndex])
      fmt.Println(matches[re.SubexpIndex("Filename")])
		}
		return nil
	})

	// fmt.Printf(`input: %s, output: %s`, input, output)
	// var buf bytes.Buffer
	// if err := goldmark.Convert(bs, &buf); err != nil {
	// 	panic(err)
	// }

	// out := fmt.Sprintf(template, string(buf.Bytes()))
	// ioutil.WriteFile("out.html", []byte(out), 0666)

}
