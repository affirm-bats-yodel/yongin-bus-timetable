package yonginbustimetable_test

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	yonginbustimetable "github.com/affirm-bats-yodel/yongin-bus-timetable"
	"github.com/stretchr/testify/assert"
)

var url = flag.String("url", "", "url to test")

// doGet handle Get Request and return body
func doGet(ctx context.Context, reqURL string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status is not OK: %d", res.StatusCode)
	}

	return res.Body, nil
}

func TestBusLinkExtractor_Extract(t *testing.T) {
	if !assert.NotEmpty(t, url, "url should not be empty") {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	body, err := doGet(ctx, *url)
	if err != nil {
		t.Error("while getting body", "url", *url)
		return
	}
	defer body.Close()

	busLinkExt, err := yonginbustimetable.NewBusListExtractor(body, *url)
	if err != nil {
		t.Error("while creating busLinkExtractor", "error", err)
		return
	}

	busLinks, err := busLinkExt.Extract(ctx)
	if assert.NoError(t, err) && assert.NotEmpty(t, busLinks) {
		for _, elem := range busLinks {
			assert.NotEmpty(t, elem.Name)
			assert.NotEmpty(t, elem.WindowOpenLink)
			assert.Equal(t, true, strings.HasPrefix(elem.WindowOpenLink, "http://"), elem.WindowOpenLink)
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

func TestBusTimetableExtractor_Extract(t *testing.T) {
	if !assert.NotEmpty(t, url) {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("getting bus list", "url", *url)

	body, err := doGet(ctx, *url)
	if err != nil {
		t.Error("getting bus list", "error", err)
		return
	}
	defer body.Close()

	busListExt, err := yonginbustimetable.NewBusListExtractor(body, *url)
	if err != nil {
		t.Error("creating bus list extractor", "error", err)
		return
	}
	busList, err := busListExt.Extract(ctx)
	if err != nil {
		t.Error("extracting bus list", "error", err)
		return
	} else if !assert.NotEmpty(t, busList) {
		return
	}

	t.Log("closing body")

	if err := body.Close(); err != nil {
		t.Error("closing body", "error", err)
		return
	}

	firstBus := busList[0]

	t.Log("extract bus", "firstBus", firstBus)

	timetableExtractor := yonginbustimetable.NewBusTimetableExtractor()
	timetable, err := timetableExtractor.Extract(ctx, firstBus)
	if err != nil {
		t.Error("getting bus timetable", "error", err)
		return
	}

	t.Log("got bus timetable", "timetable", timetable)

	assert.NotEmpty(t, timetable.Stops, "there are more than one bus stops required")
	assert.NotEmpty(t, timetable.Timetables, "there are more than one bus timetable required")
}

func TestTimetable_ExtractTime(t *testing.T) {
	tcs := []struct {
		tt           yonginbustimetable.Timetable
		value        string
		commentFound bool
	}{
		{
			tt: yonginbustimetable.Timetable{
				Stop:     "용인출발",
				DepartAt: "6:05",
			},
			value:        "6:05",
			commentFound: false,
		},
		{
			tt: yonginbustimetable.Timetable{
				Stop:     "용인출발",
				DepartAt: "7:00(전세)평일운행",
			},
			value:        "7:00",
			commentFound: true,
		},
		{
			tt: yonginbustimetable.Timetable{
				Stop:     "용인출발",
				DepartAt: "10:35(전세)평일운행",
			},
			value:        "10:35",
			commentFound: true,
		},
	}

	for _, tc := range tcs {
		s, b, err := tc.tt.ExtractTime()
		assert.NoError(t, err)
		assert.Equal(t, tc.value, s)
		assert.Equal(t, tc.commentFound, b)
	}
}
