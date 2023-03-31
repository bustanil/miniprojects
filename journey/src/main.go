package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
  "time"
	"github.com/yuin/goldmark"
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
  title string
	path string
  createdDate time.Time
}

func EndsWith(s string, endsWith string) bool {
  var l = strings.LastIndex(s, endsWith) 
	return l >= 0 && l == (len(s) - len(endsWith))
}

func main() {
	args := os.Args
	input := args[1]
	output := args[2]


  blogs := make([]BlogEntry, 0)

  // walk through the input directory and find markdown files
	rootFilePath, _ := filepath.Abs(input)
  re := regexp.MustCompile(`(?P<Year>\d{4})/(?P<Month>\d{2})/(?P<Day>\d{2})/(?P<Filename>.*)\.md`)
	filepath.WalkDir(input, func(path string, d fs.DirEntry, err error) error {
		subDirPath := strings.ReplaceAll(path, rootFilePath, "")
		if EndsWith(d.Name(), ".md") {
      matches := re.FindStringSubmatch(subDirPath)
      day := matches[re.SubexpIndex("Day")]
      month := matches[re.SubexpIndex("Month")]
      year := matches[re.SubexpIndex("Year")]
      filename := matches[re.SubexpIndex("Filename")]
      title := strings.ReplaceAll(filename, "-", "")
      createdDate, err := time.Parse("1-2-2006", fmt.Sprintf("%s-%s-%s", day, month, year))
      if err != nil {
        fmt.Println("error", err)
        return err
      }

      blog := BlogEntry{
        title: title,
        path: subDirPath,
        createdDate: createdDate,
      }

      var bs, err = ioutil.ReadFile(path)
      if err != nil {
        fmt.Println("error", err)
        return err
      }

      blogs = append(blogs, blog)

      // for each markdown file, generate the HTML page
      fmt.Printf(`input: %s, output: %s\n`, input, output)
      var buf bytes.Buffer
      if err := goldmark.Convert(bs, &buf); err != nil {
        panic(err)
      }
      out := fmt.Sprintf(template, string(buf.Bytes()))
      outputDir := fmt.Sprintf("/%s/%s/%s/%s", output, year, month, day)
      outputPath := fmt.Sprintf("/%s/%s.html", outputDir, filename)
      os.MkdirAll(outputDir, 0777)
      fmt.Println(outputPath)
      err = ioutil.WriteFile(outputPath, []byte(out), 0766)
      if err != nil {
        fmt.Println("error", err)
        return err
      }
      fmt.Println(blogs)
		}


		return nil
	})

}
