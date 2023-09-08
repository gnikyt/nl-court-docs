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

	if os.Args[3] == "json" {
		out, err := ncd.OutputJSON(d, true)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(out))
	} else {
		out := ncd.OutputText(d)
		fmt.Print(out)
	}
}
