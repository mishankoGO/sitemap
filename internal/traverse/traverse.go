package traverse

import (
	"fmt"
	"github.com/mishankoGO/sitemap/internal/link"
	"strings"
)

func extractLinks(l *string, allLinks map[string]string) ([]link.Link, error) {
	p, err := link.NewParser(*l)
	if err != nil {
		return nil, err
	}

	links, err := p.ExtractLinks()
	if err != nil {
		return nil, err
	}

	var newLinks []link.Link
	for _, link_ := range links {
		if _, ok := allLinks[link_.Url]; !ok {
			newLinks = append(newLinks, link_)
		}
	}

	return newLinks, nil
}

func Traverse(l *string, root string, allLinks map[string]string) error {

	// root link traverse
	links, err := extractLinks(l, allLinks)
	if err != nil {
		return err
	}

	if len(links) == 0 {
		return nil
	}

	k := appendLinks(links, allLinks)
	if k == 0 {
		return nil
	}
	for _, link_ := range links {
		if !strings.Contains(link_.Url, root) {
			link_.Url = fmt.Sprintf("%s/%s", root, link_.Url)
		}
		if _, ok := allLinks[link_.Url]; !ok {
			err = Traverse(&link_.Url, root, allLinks)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func appendLinks(links []link.Link, allLinks map[string]string) int {
	k := 0
	for _, link_ := range links {
		if _, ok := allLinks[link_.Url]; !ok {
			allLinks[link_.Url] = link_.Text
			k++
		}
	}
	return k
}
