package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
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
		if !path.IsAbs(filename) {
			filename, _ = filepath.Abs(filename)
		}
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
	base := path.Base(filename)
	title := base[:len(base)-3]
	render := blackfriday.HtmlRenderer(htmlFlags, title, "")
	output := blackfriday.Markdown(input, render, extensions)
	fp := dirname + "/" + title + ".html"
	if !path.IsAbs(fp) {
		fp, _ = filepath.Abs(fp)
	}
	if _, err := os.Stat(fp); err == nil {
		//fmt.Fprintf(os.Stderr, "%s exist, do you want to continue?", fp, err)
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("%s exist, do you want to continue which will overwrite it? (y/n)", fp)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		text = strings.ToLower(text)
		if text == "y" {
			Render(string(output), title, fp)
		} else {
			fmt.Println(" stop ")
		}
	}
}

func Render(content, title, fp string) {
	t := template.New("template.html")
	tp, err := t.ParseFiles("template.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing template.html ", err)
		os.Exit(-1)
	}
	var out *os.File
	if out, err = os.Create(fp); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating %s: %v", fp, err)
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

func GenIndex(root string) {
	fs, err := ioutil.ReadDir(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, " error scanning the root dir %v", err)
		os.Exit(-1)
	}
	for f := range fs {
		fmt.Println(f.Name())
	}
}
