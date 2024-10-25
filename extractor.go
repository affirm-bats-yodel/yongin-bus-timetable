package yonginbustimetable

import (
	"context"
	"errors"
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// NewBusListExtractor Create new BusListExtractor from Reader
func NewBusListExtractor(r io.Reader, urls ...string) (*BusLinkExtractor, error) {
	var url string
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	if urls != nil && urls[0] != "" {
		url = urls[0]
	}
	return &BusLinkExtractor{
		Doc: doc,
		URL: url,
	}, nil
}

// BusLinkExtractor Bus List Extractor
type BusLinkExtractor struct {
	// Doc Goquery Document
	Doc *goquery.Document
	// URL a Request URL to make full bus timetable url
	URL string
}

// Extract the button elements
func (b *BusLinkExtractor) Extract(ctx context.Context) ([]*BusLink, error) {
	var (
		links  []*BusLink
		absURL string
	)

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if b.URL != "" {
		u, err := url.Parse(b.URL)
		if err != nil {
			return nil, errors.Join(errors.New("extract url"), err)
		} else if u.Scheme != "" {
			absURL = u.Scheme + "://" + u.Host
		} else {
			absURL = "https://" + u.Host
		}
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
				WindowOpenLink: absURL + popUpEndpoint,
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
