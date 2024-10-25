package yonginbustimetable

import (
	"context"
	"io"
	"regexp"
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

		busName := s.Find("b").First().Text()
		busRoute := s.Find("span").First().Text()

		// Trim `window.open(`, `'` and split using `,` quote
		onClickSplitRes := strings.SplitN(strings.ReplaceAll(strings.TrimPrefix(onClickValue, "window.open("), "'", ""), ",", 3)
		if len(onClickSplitRes) != 3 {
			return
		} else if popUpEndpoint := onClickSplitRes[0]; !strings.HasPrefix(popUpEndpoint, "/board") {
			return
		} else {
			links = append(links, &BusLink{
				Name:           busName,
				Route:          busRoute,
				WindowOpenLink: popUpEndpoint,
			})
		}
	})

	return links, nil
}

// BusLink Link information from `button` element
type BusLink struct {
	// Name Bus name
	//
	// extracted from button's b tag
	Name string `json:"name"`
	// Route Bus Route
	Route string `json:"route,omitempty"`
	// WindowOpenLink a Extracted URL from button's onclick attribute
	WindowOpenLink string `json:"windowOpenLink"`
}

// busNumberRegexp Regexp for Extract Bus number from Name
//
// If Name is "시내2번", it will extracted as "2"
var busNumberRegexp = regexp.MustCompile(`(\d+[\d-]*)`)

// ExtractBusNumber Extract Bus number using regexp and return exact bus number
func (b *BusLink) ExtractBusNumber() string {
	if b.Name == "" {
		return ""
	} else if v := busNumberRegexp.FindString(b.Name); v != "" {
		return v
	} else {
		return b.Name
	}
}
