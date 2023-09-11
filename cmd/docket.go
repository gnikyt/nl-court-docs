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
	d.Parse(res)

	var out ncd.Outputter
	if os.Args[3] == "json" {
		out = ncd.JSONOutputter{Data: d.Data, Pretty: true}
	} else {
		out = ncd.TextOutputter{Data: d.Data}
	}
	fmr, err := out.Format()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(fmr)
}
