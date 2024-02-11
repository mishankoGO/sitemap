package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/mishankoGO/sitemap/internal/traverse"
	"log"
	"os"
)

func main() {

	l := flag.String("l", "", "root link")
	//depth := flag.Int("depth", -1, "depth of search")
	flag.Parse()

	allLinks := make(map[string]string)

	err := traverse.Traverse(l, *l, allLinks)
	if err != nil {
		log.Fatal(err)
	}

	type URL struct {
		XMLName xml.Name `xml:"loc"`
		url     string   `xml:"url"`
	}

	type SiteMap struct {
		XMLName xml.Name `xml:"urlset"`
		root    string   `xml:"xmlns,attr"`
		URLs    []string `xml:"url>loc"`
	}

	var urls []string
	for k := range allLinks {
		urls = append(urls, k)
	}

	v := &SiteMap{root: *l, URLs: urls}

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("  ", "    ")
	if err := enc.Encode(v); err != nil {
		fmt.Printf("error: %v\n", err)
	}

}
