package parsers

import (
	"fmt"
	"github.com/osamikoyo/geass-v2/internal/models"
	"log"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type Parser struct {
	
}

// Глобальные переменные
var (
	maxDepth    = 3
	visitedURLs = make(map[string]bool)
	mu          sync.Mutex
	wg          sync.WaitGroup
)

func extractLinks(url string) ([]models.Link, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var links []models.Link
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					text := strings.TrimSpace(getText(n))
					links = append(links, models.Link{Text: text, Href: attr.Val})
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return links, nil
}

func getText(n *html.Node) string {
	var text string
	if n.Type == html.TextNode {
		text = n.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += getText(c)
	}
	return text
}

func extractContent(url string) (models.PageInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return models.PageInfo{}, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return models.PageInfo{}, err
	}

	var pageInfo models.PageInfo
	pageInfo.Url = url
	pageInfo.Technical.Code = uint32(resp.StatusCode)
	pageInfo.Technical.ContentType = resp.Header.Get("Content-Type")

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				pageInfo.Title = getText(n)
			case "meta":
				var name, content string
				for _, attr := range n.Attr {
					if attr.Key == "name" {
						name = attr.Val
					}
					if attr.Key == "content" {
						content = attr.Val
					}
				}
				if name == "description" {
					pageInfo.MetadataDescription = content
				}
				if name == "robots" {
					pageInfo.Metadata.Robots = content
				}
			case "img":
				var src, alt string
				for _, attr := range n.Attr {
					if attr.Key == "src" {
						src = attr.Val
					}
					if attr.Key == "alt" {
						alt = attr.Val
					}
				}
				pageInfo.Content.Images = append(pageInfo.Content.Images, models.Image{Src: src, Alt: alt})
			}
		}
		if n.Type == html.TextNode {
			pageInfo.Content.FullText += n.Data + " "
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return pageInfo, nil
}
func ParsePage(url string, depth int) {
	defer wg.Done()

	mu.Lock()
	if visitedURLs[url] || depth > maxDepth {
		mu.Unlock()
		return
	}
	visitedURLs[url] = true
	mu.Unlock()

	pageInfo, err := extractContent(url)
	if err != nil {
		log.Printf("Error extracting content from %s: %v\n", url, err)
		return
	}

	fmt.Printf("URL: %s\nTitle: %s\nDescription: %s\nContent: %s\n\n",
		pageInfo.Url, pageInfo.Title, pageInfo.MetadataDescription, pageInfo.Content.FullText)

	links, err := extractLinks(url)
	if err != nil {
		log.Printf("Error extracting links from %s: %v\n", url, err)
		return
	}

	for _, link := range links {
		wg.Add(1)
		go ParsePage(link.Href, depth+1)
	}
}
