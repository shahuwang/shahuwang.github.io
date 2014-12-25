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

var ignore []string = []string{".git", "css", "js", "images", "exec", "index.html"}

func main() {
	var filename, dirname string
	var latex bool
	flag.StringVar(&filename, "file", "", "The path of the file you want to generate")
	flag.StringVar(&dirname, "dir", "", "Where the generated HTML page you would like to put on")
	flag.BoolVar(&latex, "latex", false, "whether to include mathjax to render LaTex")
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
			Render(string(output), title, fp, latex)
		} else {
			fmt.Println(" stop ")
			return
		}
	} else {
		Render(string(output), title, fp, latex)
	}
	t := template.New("indextpl.html")
	tp, err := t.ParseFiles("indextpl.html")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	GenIndex("../", tp)
}

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

func Render(content, title, fp string, latex bool) {
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
		Latex   template.HTML
	}
	tc.Content = template.HTML(content)
	tc.Title = template.HTML(title)
	if latex {
		tc.Latex = ` 
		<script type="text/javascript"
  			src="http://cdn.mathjax.org/mathjax/latest/MathJax.js?config=TeX-AMS-MML_HTMLorMML">
		</script>
		`
	} else {
		tc.Latex = ""
	}
	err = tp.Execute(out, tc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error  execute template rendering: %v", "template.html", err)
		os.Exit(-1)
	}
}
