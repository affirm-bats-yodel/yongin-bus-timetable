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

func TestBusLink_ExtractBusNumber(t *testing.T) {
	var tc = []struct {
		b       *yonginbustimetable.BusLink
		equalTo string
	}{
		{
			b: &yonginbustimetable.BusLink{
				Name: "시내66번",
			},
			equalTo: "66",
		},
		{
			b: &yonginbustimetable.BusLink{
				Name: "시내66-4번",
			},
			equalTo: "66-4",
		},
		{
			b: &yonginbustimetable.BusLink{
				Name: "시내5700번",
			},
			equalTo: "5700",
		},
		{
			b: &yonginbustimetable.BusLink{
				Name: "마을201번",
			},
			equalTo: "201",
		},
	}

	for _, elem := range tc {
		v := elem.b.ExtractBusNumber()
		t.Log(v)
		assert.Equal(t, elem.equalTo, v)
	}
}
