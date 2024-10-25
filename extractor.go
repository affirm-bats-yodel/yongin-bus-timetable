package yonginbustimetable

import (
	"context"
	"io"
	"log"

	"github.com/PuerkitoBio/goquery"
)

// NewBusListExtractor Create new BusListExtractor from Reader
func NewBusListExtractor(r io.Reader) (*BusLinkExtractor, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	return &BusLinkExtractor{
		Doc: doc,
	}, nil
}

// BusLinkExtractor Bus List Extractor
type BusLinkExtractor struct {
	// Doc Goquery Document
	Doc *goquery.Document
}

// Extract the button elements
func (b *BusLinkExtractor) Extract(ctx context.Context) ([]*BusLink, error) {
	var links []*BusLink

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	b.Doc.Find("button").Each(func(i int, s *goquery.Selection) {
		log.Println("index", i, "element", s.Text())
	})

	return links, nil
}

// BusLink Link information from `button` element
type BusLink struct {
	// Name Bus name
	Name string `json:"name"`
	// WindowOpenLink a Extracted URL from button's onclick attribute
	WindowOpenLink string `json:"windowOpenLink"`
}
