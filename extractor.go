package yonginbustimetable

import (
	"context"
	"io"
	"strings"

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
		onClickValue, exist := s.Attr("onclick")
		if !exist {
			return
		} else if !strings.HasPrefix(onClickValue, "window.open") {
			return
		}

		// Trim `window.open(`, `'` and split using `,` quote
		onClickSplitRes := strings.SplitN(strings.ReplaceAll(strings.TrimPrefix(onClickValue, "window.open("), "'", ""), ",", 3)
		if len(onClickSplitRes) != 3 {
			return
		} else if popUpEndpoint := onClickSplitRes[0]; !strings.HasPrefix(popUpEndpoint, "/board") {
			return
		} else {
			links = append(links, &BusLink{
				Name:           s.Text(),
				WindowOpenLink: popUpEndpoint,
			})
		}
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
