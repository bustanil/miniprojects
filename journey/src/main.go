package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
  "os"
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

func main() {
	var bs, err = ioutil.ReadFile("test.md")
	if err != nil {
		fmt.Println("error")
		return
	}

  var buf bytes.Buffer
  if err := goldmark.Convert(bs, &buf); err != nil {
    panic(err)
  }

  out := fmt.Sprintf(template, string(buf.Bytes()))
  ioutil.WriteFile("out.html", []byte(out), 0666)

  args := os.Args
  input := args[1]
  output := args[2]

  fmt.Printf(`input: %s, output: %s`, input, output)
}
