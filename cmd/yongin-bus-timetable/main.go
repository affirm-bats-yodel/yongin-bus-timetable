package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/affirm-bats-yodel/yongin-bus-timetable"
	_ "github.com/marcboeker/go-duckdb"
	"github.com/urfave/cli/v2"
)

const (
	// insertBusQuery Query to Insert Bus Information to DuckDB
	insertBusQuery = "INSERT OR IGNORE INTO bus_lists (name) VALUES (?)"
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
				duckdbPath := c.String("duckdb-path")

				log.Println("open duckdb", "duckdbPath", duckdbPath)
				db, err := sql.Open("duckdb", duckdbPath)
				if err != nil {
					return cli.Exit(fmt.Errorf("error: open duckdb: %q: %v", duckdbPath, err), 1)
				}
				defer db.Close()

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

				log.Println("begin transaction")

				tx, err := db.BeginTx(ctx, &sql.TxOptions{})
				if err != nil {
					return cli.Exit(fmt.Errorf("error: begin transaction: %v", err), 1)
				}

				for _, elem := range busLinks {
					_, err := tx.ExecContext(ctx, insertBusQuery, elem.ExtractBusNumber())
					if err != nil {
						return cli.Exit(fmt.Errorf("error: inserting bus: %v", err), 1)
					}
				}

				log.Println("commit transaction")
				if err := tx.Commit(); err != nil {
					return cli.Exit(fmt.Errorf("error: commit transaction: %v", err), 1)
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
				&cli.StringFlag{
					Name:     "duckdb-path",
					Required: true,
					Usage:    "path of duckdb db file",
					EnvVars: []string{
						"YONGIN_BUS_TIMETABLE_DUCKDB_PATH",
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
