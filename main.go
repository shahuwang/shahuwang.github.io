package main

import (
	"flag"
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"os"
)

func main() {
	var filename, dirname string
	flag.StringVar(&filename, "file", "", "The path of the file you want to generate")
	flag.StringVar(&dirname, "dir", "", "Where the generated HTML page you would like to put on")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: newblog --file=test.md --dir=blog")
		flag.PrintDefaults()
	}
	flag.Parse()
	var err error
	var input []byte
	if filename != "" && dirname != "" {
		if input, err = ioutil.ReadFile(filename); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from ", filename, ":", err)
			os.Exit(-1)
		}
	} else {
		flag.Usage()
		os.Exit(-1)
	}
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_OMIT_CONTENTS
	render := blackfriday.HtmlRenderer(htmlFlags, filename[:len(filename)-3], "")
	output := blackfriday.Markdown(input, render, extensions)
	title := filename[:len(filename)-3]
	Render(string(output), title, dirname, filename)
}

func Render(content, title, dirname, filename string) {
	t := template.New("template.html")
	tp, err := t.ParseFiles("template.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing template.html ", err)
		os.Exit(-1)
	}
	var out *os.File
	outfilename := filename[:len(filename)-3] + ".html"
	if out, err = os.Create(dirname + "/" + outfilename); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating %s: %v", dirname+outfilename, err)
		os.Exit(-1)
	}
	var tc struct {
		Title   template.HTML
		Content template.HTML
	}
	tc.Content = template.HTML(content)
	tc.Title = template.HTML(title)
	err = tp.Execute(out, tc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error  execute template rendering: %v", "template.html", err)
		os.Exit(-1)
	}
}
