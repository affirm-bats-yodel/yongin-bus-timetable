# yongin-bus-timetable

Yongin Bus Timetable Scraper &amp; Viewer using Go

> See [README_ko.md](./README_ko.md) if you're a Korean.

## Why?

* In Yongin, There are several buses are having a timetable.
* However, most of the map service (such as, N\*V\*R, K\*K\*O, etc) does not supports
Specific Bus Timetable.

## How?

1. Scrape Bus Timetable from [Yongin Bus Terminal](http://knyongintr.co.kr)
2. Extract Links from Button (that usually trigger `window.open` to show timetables of a bus.)
3. Go to extracted address from `window.open` and extract `<table>...</table>`
4. Store it to DB (or `json` formatted code to easily retrieve from Static Site.)
5. Profit!
