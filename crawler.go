package main

import (
	"fmt"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// sample url :
// with category - https://www.leboncoin.fr/<category>/offres/?q=test&&ps=0&pe=2
// no category   - https://www.leboncoin.fr/annonces/offres/?q=test&latitude=43.297&longitude=5.3875&radius=30000&ps=3&pe=7
//               - https://www.leboncoin.fr/annonces/offres/?q=test&latitude=43.297&longitude=5.3875&radius=30000
const LBCScheme = "https"
const LBCHost = "https://www.leboncoin.fr"
const LBCSearchPrefix = LBCHost + "/annonces/offres/?"
const LBCSearchCatPrefixFormat = LBCHost + "/%s/offres/?"

type Query struct {
	Search   string
	RadiusKM int // if 0, search france entière
	Lng, Lat float32

	Category string // category, as displayed in url /
	PriceMin int    // index from 0 to len(category:PriceMin) - check lbc web interface
	PriceMax int    // index from 0 to len(category:PriceMax) - check lbc web interface

	// extra args for category-specific parameters not found in Query struct otherwise
	// - check lbc/category web interface for specifics
	RawArgs map[string]string
}

type Item struct {
	URL      string
	Title    string
	Price    string
	ThumbURL string
}

// BuildURL builds query url using q information.
func (q *Query) BuildURL() string {
	var qURL string
	var bcat bool // are we searching site-wide or in category
	if q.Category != "" {
		bcat = true
		qURL = fmt.Sprintf(LBCSearchCatPrefixFormat, q.Category)
	} else {
		qURL = LBCSearchPrefix
	}

	if q.Search != "" {
		qURL += fmt.Sprintf("q=%s&", url.QueryEscape(q.Search))
	}
	if q.RadiusKM != 0 {
		if q.Lat == 0.0 && q.Lng == 0.0 {
			// try geoloc
			q.Lat, q.Lng = getLatLng()
		}

		if q.Lat == 0.0 && q.Lng == 0.0 {
			// we couldn't determine geo-loc coordinates, force site-wide search
			q.RadiusKM = 0
		} else {
			qURL += fmt.Sprintf("latitude=%f&longitude=%f&radius=%d&", q.Lat, q.Lng, getClosestRadius(q.RadiusKM))
		}
	}

	// category specific options
	if bcat {
		if q.PriceMin > 0 {
			qURL += fmt.Sprintf("ps=%d&", q.PriceMin)
		}
		if q.PriceMax > 0 {
			qURL += fmt.Sprintf("pe=%d&", q.PriceMax)
		}

		// extra-parameters
		for k, v := range q.RawArgs {
			qURL += fmt.Sprintf("%s=%s&", k, v)
		}
	}

	// remove trailing '&'
	if idx := len(qURL) - 1; qURL[idx] == '&' {
		qURL = qURL[:idx]
	}

	return qURL
}

// Run searches leboncoin using parameters from q.
// If we get a 200 http response, parse html and extracts []Item.
func (q *Query) Run() ([]Item, error) {
	resp, err := http.Get(q.BuildURL())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status: %s", resp.StatusCode)
	}

	return parseHtml(resp.Body)
}

// parseHtml does parsing magic to retreive search result items from body html tree.
func parseHtml(body io.ReadCloser) ([]Item, error) {
	defer body.Close()
	node, err := html.Parse(body)
	if err != nil {
		return nil, err
	}

	var ok bool
	node, ok = scrape.Find(node, attrMatcher("class", "tabsContent"))
	if !ok {
		return nil, fmt.Errorf(".tabsContent not found")
	}

	nodeItems := scrape.FindAll(node, elemMatcher("a"))
	items := make([]Item, len(nodeItems))
	for k, n := range nodeItems {
		item := Item{}

		// scrape title, href
		for _, attr := range n.Attr {
			if attr.Key == "title" {
				item.Title = strings.TrimSpace(attr.Val)
			}
			if attr.Key == "href" {
				item.URL = cleanURL(attr.Val)
			}
		}

		// scrape tb url
		tb, ok := scrape.Find(n, attrMatcher("class", "lazyload"))
		if ok {
			for _, attr := range tb.Attr {
				if attr.Key == "data-imgsrc" {
					item.ThumbURL = cleanURL(attr.Val)
				}
			}
		}

		// scrape price
		price, ok := scrape.Find(n, attrMatcher("class", "item_price"))
		if ok {
			for _, attr := range price.Attr {
				if attr.Key == "content" {
					item.Price = attr.Val
				}
			}
		}

		items[k] = item
	}
	return items, nil
}

// cleanURL prevents // & / uri prefixes and replaces with
// proper scheme://host as it should be more rliable.
func cleanURL(url string) string {
	if len(url) <= 2 {
		return url
	}
	if url[:2] == "//" {
		return LBCScheme + ":" + url
	} else if url[0] == '/' {
		return LBCHost + url
	}
	return url
}

// available radius values in lbc, and get closest helper
var Radius = []int{10, 20, 30, 40, 50, 100, 200}

func getClosestRadius(val int) int {
	for _, v := range Radius {
		if val <= v {
			return v * 1000
		}
	}
	return Radius[len(Radius)-1] * 1000
}

// scrape.Matcher helper funcs
func or(m scrape.Matcher, matchers ...scrape.Matcher) scrape.Matcher {
	return func(n *html.Node) bool {
		result := m(n)
		for _, v := range matchers {
			if result {
				return true
			}
			result = result || v(n)
		}
		return result
	}
}

func and(m scrape.Matcher, matchers ...scrape.Matcher) scrape.Matcher {
	return func(n *html.Node) bool {
		result := m(n)
		for _, v := range matchers {
			if !result {
				return false
			}
			result = result && v(n)
		}
		return result
	}
}

func elemMatcher(elem string) scrape.Matcher {
	return func(n *html.Node) bool {
		return n.Data == elem
	}
}

func attrMatcher(key, val string) scrape.Matcher {
	return func(n *html.Node) bool {
		attr := html.Attribute{Key: key, Val: val}
		for _, v := range n.Attr {
			if v.Key == attr.Key && strings.Index(v.Val, attr.Val) >= 0 {
				return true
			}
		}
		return false
	}
}
