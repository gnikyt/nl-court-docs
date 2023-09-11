package main

import (
	"fmt"
	"log"
	"os"

	ncd "github.com/gnikyt/nl-court-docs"
)

func main() {
	d := ncd.NewDocket(os.Args[1], os.Args[2], nil)
	res, err := d.Fetch()
	if err != nil {
		log.Fatal(err)
	}
	if err := d.Parse(res); err != nil {
		log.Fatal(err)
	}

	var out ncd.Outputter
	switch os.Args[3] {
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
