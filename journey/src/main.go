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
    <title>%s</title>
    <link rel="stylesheet" href="https://unpkg.com/@picocss/pico@1.*/css/pico.min.css">
  </head>
  <body>
    <article>
      %s
    </article>
  </body>
</html>
`

var postListTemplate string = `
<html>
  <head>
    <title>%s</title>
    <link rel="stylesheet" href="https://unpkg.com/@picocss/pico@1.*/css/pico.min.css">
  </head>
  <body>
    <article>
      <h2>%s</h2>
      <hr />
      %s 
    </article>
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
  title := args[3]

  blogs := make([]BlogEntry, 0)

  // walk through the input directory and find markdown files
	rootFilePath, _ := filepath.Abs(input)
  re := regexp.MustCompile(`(?P<Year>\d{4})/(?P<Month>\d{2})/(?P<Day>\d{2})/(?P<Filename>.*)\.md`)
	filepath.WalkDir(input, func(path string, d fs.DirEntry, err error) error {
		subDirPath := strings.ReplaceAll(path, rootFilePath, "")
		if EndsWith(d.Name(), ".md") {

      // validate and parse the file path
      matches := re.FindStringSubmatch(subDirPath)
      day := matches[re.SubexpIndex("Day")]
      month := matches[re.SubexpIndex("Month")]
      year := matches[re.SubexpIndex("Year")]
      filename := matches[re.SubexpIndex("Filename")]
      title := strings.ReplaceAll(filename, "-", " ")
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
      fmt.Printf(`input: %s, output: %s`, input, output)
      var buf bytes.Buffer
      if err := goldmark.Convert(bs, &buf); err != nil {
        panic(err)
      }
      content := fmt.Sprintf(postTemplate, title, (buf.Bytes()))

      outputSubDir := fmt.Sprintf("%s/%s/%s", year, month, day)
      outputDir := fmt.Sprintf("%s/%s", output, outputSubDir)
      outputPath := fmt.Sprintf("%s/%s.html", outputDir, filename)
      relativePath := fmt.Sprintf("%s/%s.html", outputSubDir, filename)

      blog := BlogEntry{
        title: title,
        createdDate: createdDate,
        path: relativePath,
      }
      blogs = append(blogs, blog)

      // write the HTML page
      os.MkdirAll(outputDir, 0777)
      err = ioutil.WriteFile(outputPath, []byte(content), 0766)
      if err != nil {
        fmt.Println("error", err)
        return err
      }
		}

    // generate the post list
    var posts string = ""
    for _, blog := range blogs {
      var row = fmt.Sprintf(`<div><a href="%s">%s (%s)</a></div>`, blog.path, blog.title, blog.createdDate.Format("02-01-2006"))
      posts = posts + row
    }

    ioutil.WriteFile(fmt.Sprintf("%s/index.html", output), []byte(fmt.Sprintf(postListTemplate, title, title, posts)), 0766)

		return nil
	})

}
