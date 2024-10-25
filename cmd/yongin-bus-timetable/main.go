package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/affirm-bats-yodel/yongin-bus-timetable"
	"github.com/urfave/cli/v2"
)

var CLI = &cli.App{
	Name:        "yongin-bus-timetable",
	Description: "yongin bus timetable explorer",
	Commands: []*cli.Command{
		{
			Name:  "sync",
			Usage: "Extract a Bus List and save it to DB",
			Action: func(c *cli.Context) error {
				ctx, cancel := signal.NotifyContext(c.Context, syscall.SIGTERM, syscall.SIGINT)
				defer cancel()

				url := c.String("url")

				log.Println("creating request with context", "url", url)

				req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
				if err != nil {
					return cli.Exit(fmt.Errorf("error: creating request context with %q: %v", url, err), 1)
				}

				res, err := http.DefaultClient.Do(req)
				if err != nil {
					return cli.Exit(fmt.Errorf("error: error while request to %q: %v", url, err), 1)
				}
				defer res.Body.Close()

				if res.StatusCode != http.StatusOK {
					return cli.Exit(fmt.Errorf("error: response from %q: %d", url, res.StatusCode), 1)
				}

				log.Println("extract url from body")
				busLinkExt, err := yonginbustimetable.NewBusListExtractor(res.Body, url)
				if err != nil {
					return cli.Exit(fmt.Errorf("error: creating bus list extractor: %v", err), 1)
				}

				busLinks, err := busLinkExt.Extract(ctx)
				if err != nil {
					return cli.Exit(fmt.Errorf("error: extracting bus links: %v", err), 1)
				}

				if err := res.Body.Close(); err != nil {
					return cli.Exit(fmt.Errorf("error: closing response body: %v", err), 1)
				}

				for _, elem := range busLinks {
					log.Println("name", elem.ExtractBusNumber(), "route", elem.Route, "link", elem.WindowOpenLink)
				}

				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "url",
					Required: true,
					Usage:    "url of the yongin-bus-terminal",
					EnvVars: []string{
						"YONGIN_BUS_TIMETABLE_EXTRACT_URL",
					},
				},
			},
		},
	},
}

func main() {
	if err := CLI.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
