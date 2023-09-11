package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	ncd "github.com/gnikyt/nl-court-docs"
)

func main() {
	var date, office, frmt string
	flag.StringVar(&date, "date", time.Now().Format("2006-01-02"), "date in YYYY-MM-DD format, no past values")
	flag.StringVar(&office, "office", "", "office ID")
	flag.StringVar(&frmt, "format", "json", "format: json, text, or csv")
	flag.Parse()
	if office == "" {
		log.Fatal("office ID flag required")
	}

	d := ncd.NewDocket(date, office, nil)
	res, err := d.Fetch()
	if err != nil {
		log.Fatal(err)
	}
	if err := d.Parse(res); err != nil {
		log.Fatal(err)
	}

	var out ncd.Outputter
	switch frmt {
	default:
	case "json":
		out = ncd.NewPrettyJsonOutput(d.Data)
	case "text":
		out = ncd.NewTextOutput(d.Data)
	case "csv":
		out = ncd.NewCsvOutput(d.Data)
	}
	fmr, err := d.Output(out)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(fmr)
}
