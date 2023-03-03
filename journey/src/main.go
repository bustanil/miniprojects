package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
  // date date
  path string
}

func EndsWith(s string, endsWith string) bool {
  return strings.LastIndex(s, endsWith) == (len(s) - len(endsWith)) 
}

func main() {
	var bs, err = ioutil.ReadFile("test.md")
	if err != nil {
		fmt.Println("error")
		return
	}

  args := os.Args
  input := args[1]
  output := args[2]

  rootFilePath, err := filepath.Abs(input)
  filepath.WalkDir(input, func(path string, d fs.DirEntry, err error) error {
    subDirPath := strings.ReplaceAll(path, rootFilePath, "")
    if EndsWith(d.Name(), ".md") {
      fmt.Println(subDirPath)
    }
    return nil
  })
   

  fmt.Printf(`input: %s, output: %s`, input, output)
  var buf bytes.Buffer
  if err := goldmark.Convert(bs, &buf); err != nil {
    panic(err)
  }

  out := fmt.Sprintf(template, string(buf.Bytes()))
  ioutil.WriteFile("out.html", []byte(out), 0666)

}
