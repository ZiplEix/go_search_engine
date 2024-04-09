package search

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type CrawlData struct {
	Url          string
	Success      bool
	ResponseCode int
	CrawlData    ParsedBody
}

type ParsedBody struct {
	CrawlTime       time.Duration
	PageTitle       string
	PageDescription string
	Heading         string
	Links           Links
}

type Links struct {
	Internal []string
	External []string
}

func runCrawl(inputUrl string) CrawlData {
	resp, err := http.Get(inputUrl)
	baseUrl, _ := url.Parse(inputUrl)
	defer resp.Body.Close()
	if err != nil || resp == nil {
		fmt.Println("something went wrong fetching the body")
		return CrawlData{
			Url:          inputUrl,
			Success:      false,
			ResponseCode: 0,
			CrawlData:    ParsedBody{},
		}
	}

	if resp.StatusCode != 200 {
		fmt.Println("non 200 response code")
		return CrawlData{
			Url:          inputUrl,
			Success:      false,
			ResponseCode: resp.StatusCode,
			CrawlData:    ParsedBody{},
		}
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		fmt.Println("non html content type")
		return CrawlData{
			Url:          inputUrl,
			Success:      false,
			ResponseCode: resp.StatusCode,
			CrawlData:    ParsedBody{},
		}
	}

	// parse the body
	data, err := parseBody(resp.Body, baseUrl)
	if err != nil {
		fmt.Println("unable to parse body")
		return CrawlData{
			Url:          inputUrl,
			Success:      false,
			ResponseCode: resp.StatusCode,
			CrawlData:    ParsedBody{},
		}
	}

	return CrawlData{
		Url:          inputUrl,
		Success:      true,
		ResponseCode: resp.StatusCode,
		CrawlData:    data,
	}
}

func parseBody(body io.Reader, baseUrl *url.URL) (ParsedBody, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return ParsedBody{}, fmt.Errorf("unable to parse body: %v", err)
	}

	start := time.Now()

	// Get Links
	links := getLinks(doc, baseUrl)

	// Get page title and description
	title, desc := getPageData(doc)

	// Get H1 tags
	heading := getPageHeadings(doc)

	// return the time & data
	end := time.Now()
	return ParsedBody{
		CrawlTime:       end.Sub(start),
		PageTitle:       title,
		PageDescription: desc,
		Heading:         heading,
		Links:           links,
	}, nil
}

// Depth first search (dfs) (https://en.wikipedia.org/wiki/Depth-first_search)
//
// Recursive function for scanning the html tree
func getLinks(node *html.Node, baseUrl *url.URL) Links {
	links := Links{}
	if node == nil {
		return links
	}

	var findLinks func(*html.Node)

	findLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					url, err := url.Parse(attr.Val)
					if err != nil || strings.HasPrefix(url.String(), "#") || strings.HasPrefix(url.String(), "mail") || strings.HasPrefix(url.String(), "tel") || strings.HasPrefix(url.String(), "javascript") || strings.HasPrefix(url.String(), "pdf") || strings.HasPrefix(url.String(), ".md") {
						continue
					}

					if url.IsAbs() {
						if isSameHost(url.String(), baseUrl.String()) {
							links.Internal = append(links.Internal, url.String())
						} else {
							links.External = append(links.External, url.String())
						}
					} else {
						rel := baseUrl.ResolveReference(url)
						links.Internal = append(links.Internal, rel.String())
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findLinks(c)
		}
	}

	findLinks(node)

	return links
}

func isSameHost(absoluteUrl string, baseUrl string) bool {
	absUrl, err := url.Parse(absoluteUrl)
	if err != nil {
		return false
	}

	base, err := url.Parse(baseUrl)
	if err != nil {
		return false
	}

	return absUrl.Host == base.Host
}

func getPageData(node *html.Node) (string, string) {
	if node == nil {
		return "", ""
	}

	var title, desc string
	var findMetaAndTitle func(*html.Node)

	findMetaAndTitle = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				title = n.FirstChild.Data
			}
		} else if n.Type == html.ElementNode && n.Data == "meta" {
			var name, content string

			for _, attr := range node.Attr {
				if attr.Key == "name" {
					name = attr.Val
				} else if attr.Key == "content" {
					content = attr.Val
				}
			}

			if name == "description" {
				desc = content
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		findMetaAndTitle(child)
	}

	findMetaAndTitle(node)

	return title, desc
}

func getPageHeadings(n *html.Node) string {
	if n == nil {
		return ""
	}

	var heading strings.Builder
	var findH1 func(*html.Node)

	findH1 = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h1" {
			if n.FirstChild != nil {
				heading.WriteString(n.FirstChild.Data)
				heading.WriteString(", ")
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findH1(c)
		}
	}

	findH1(n)

	return strings.TrimSuffix(heading.String(), ", ")
}
