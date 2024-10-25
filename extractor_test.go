package yonginbustimetable_test

import (
	"context"
	"flag"
	"net/http"
	"testing"

	yonginbustimetable "github.com/affirm-bats-yodel/yongin-bus-timetable"
	"github.com/stretchr/testify/assert"
)

var url = flag.String("url", "", "url to test")

func TestBusLinkExtractor_Extract(t *testing.T) {
	if !assert.NotEmpty(t, url, "url should not be empty") {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", *url, nil)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error("got error getting url: ", "url", url, "error", err)
		return
	}
	defer res.Body.Close()

	busLinkExt, err := yonginbustimetable.NewBusListExtractor(res.Body)
	if err != nil {
		t.Error("while creating busLinkExtractor", "error", err)
		return
	}

	busLinks, err := busLinkExt.Extract(ctx)
	if assert.NoError(t, err) && assert.NotEmpty(t, busLinks) {
		for _, elem := range busLinks {
			assert.NotEmpty(t, elem.Name)
			assert.NotEmpty(t, elem.WindowOpenLink)
		}
	}
}
