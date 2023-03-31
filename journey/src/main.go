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

// TODO extract the postTemplate into a separate file
var postTemplate string = `
<html>
  <head>
    <link rel="stylesheet" href="https://unpkg.com/@picocss/pico@1.*/css/pico.min.css">
  </head>
  <body>
    %s
  </body>
</html>
`

var postListTemplate string = `
<html>
  <head>
    <link rel="stylesheet" href="https://unpkg.com/@picocss/pico@1.*/css/pico.min.css">
  </head>
  <body>
    %s 
  </body>
</html> `

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
      createdDate, err := time.Parse("02-01-2006", fmt.Sprintf("%s-%s-%s", day, month, year))
      if err != nil {
        fmt.Println("error", err)
        return err
      }

      bs, err := ioutil.ReadFile(path)
      if err != nil {
        fmt.Println("error", err)
        return err
      }

      // for each markdown file, generate the HTML page
      fmt.Printf(`input: %s, output: %s\n`, input, output)
      var buf bytes.Buffer
      if err := goldmark.Convert(bs, &buf); err != nil {
        panic(err)
      }
      out := fmt.Sprintf(postTemplate, string(buf.Bytes()))
      outputSubDir := fmt.Sprintf("/%s/%s/%s", year, month, day)
      outputDir := fmt.Sprintf("/%s/%s", output, outputSubDir)
      outputPath := fmt.Sprintf("/%s/%s.html", outputDir, filename)
      relativePath := fmt.Sprintf("%s/%s.html", outputSubDir, filename)

      blog := BlogEntry{
        title: title,
        createdDate: createdDate,
        path: relativePath,
      }
      blogs = append(blogs, blog)

      os.MkdirAll(outputDir, 0777)
      err = ioutil.WriteFile(outputPath, []byte(out), 0766)
      if err != nil {
        fmt.Println("error", err)
        return err
      }
      fmt.Println(blogs)
		}

    // generate the post list
    var posts string = ""
    for _, blog := range blogs {
      fmt.Println(blog.path)
      var row = fmt.Sprintf(`<div><a href="%s">%s (%s)</a></div>`, blog.path[1:len(blog.path)], blog.title, blog.createdDate)
      posts = posts + row
    }

    ioutil.WriteFile(fmt.Sprintf("%s/index.html", output), []byte(fmt.Sprintf(postListTemplate, posts)), 0766)

		return nil
	})

}
