package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var ignore []string = []string{".git", "css", "js", "images", "exec", "index.html"}

func GenIndex(root string, tpl *template.Template) {
	fs, err := ioutil.ReadDir(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, " error scanning the root dir %v", err)
		os.Exit(-1)
	}
	var link []string
	for i := range fs {
		f := fs[i]
		br := false
		for _, ig := range ignore {
			if ig == f.Name() {
				br = true
			}
		}
		if br {
			continue
		}
		nest := filepath.Join(root, f.Name())
		if f.IsDir() {
			GenIndex(nest, tpl)
			l := fmt.Sprintf("<a href='%s'><li><h2>%s</h2></li></a>\n", f.Name()+"/index.html", f.Name())
			link = append(link, l)
		} else {
			if ext := filepath.Ext(nest); ext == ".html" {
				noext := strings.TrimSuffix(f.Name(), ext)
				l := fmt.Sprintf("<a href='%s'><li><h2>%s</h2></li></a>\n", f.Name(), noext)
				link = append(link, l)
			}
		}
	}
	var tc struct {
		Content template.HTML
	}
	tc.Content = template.HTML(strings.Join(link, "<p/>"))
	out, err := os.Create(filepath.Join(root, "index.html"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating index.html %s, %v", filepath.Join(root, "index.html"), err)
		os.Exit(-1)
	}
	err = tpl.Execute(out, tc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %v", err)
		os.Exit(-1)
	}

}
func main() {
	t := template.New("indextpl.html")
	tp, err := t.ParseFiles("indextpl.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	GenIndex("../", tp)
}
