package main

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gabriel-vasile/mimetype"
)

var (
	//go:embed tpl.html
	tpl   string
	wd, _ = os.Getwd()
)

func main() {
	wd = filepath.Base(wd)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(tpl))
	doc.Find("title").SetText(wd)
	doc.Find("p#title").SetText(wd)
	for _, arg := range os.Args[1:] {
		mime, err := mimetype.DetectFile(arg)
		if err != nil {
			log.Fatalf("check %s error: %s", arg, err)
			return
		}
		switch {
		case strings.HasPrefix(mime.String(), "image"):
			tag := fmt.Sprintf(`<img src="%s"></img>`, arg)
			doc.Find("div#content").AppendHtml(tag)
		case strings.HasPrefix(mime.String(), "video"):
			tag := fmt.Sprintf(`<video src="%s" controls></video>`, arg)
			doc.Find("div#content").AppendHtml(tag)
		default:
			log.Fatalf("unsupported mime")
			return
		}
	}

	target := wd + ".html"
	html, err := doc.Html()
	if err != nil {
		log.Fatalf("build %s error: %s", target, err)
		return
	}

	err = ioutil.WriteFile(target, []byte(html), 0666)
	if err != nil {
		log.Fatalf("write %s error: %s", target, err)
		return
	}
}
