package link

import (
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strings"
)

type Link struct {
	Url  string
	Text string
}

type Parser struct {
	Page string
}

func NewParser(pagePath string) (*Parser, error) {
	res, err := http.Get(pagePath)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	res.Body.Close()

	return &Parser{Page: string(body)}, nil
}

func (p *Parser) parse() (*html.Node, error) {
	doc, err := html.Parse(strings.NewReader(p.Page))
	if err != nil {
		log.Println("error parsing page in Parse()")
		return nil, err
	}
	return doc, nil
}

func (p *Parser) ExtractLinks() ([]Link, error) {
	doc, err := p.parse()
	if err != nil {
		return nil, err
	}
	// 1. Find <a> nodes in doc
	// 2. for each link node...
	//	2a. build a link
	// 3. return the links
	nodes := linkNodes(doc)
	var links []Link
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}
	return links, nil
}

func buildLink(n *html.Node) Link {
	var ret Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Url = attr.Val
			break
		}
	}
	ret.Text = text(n)
	return ret
}

func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}

	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += text(c) + " "
	}
	return strings.Join(strings.Fields(ret), " ")
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}

	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}
