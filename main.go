package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// docket represents the data parsed from the Court Docket.
type docket struct {
	Data    map[string]map[string][][]string // format of... { [time]: [{ [name]: [[charge, count]] }] }.
	Current string                           // current time slot.
}

// newDocket inits a docket struct.
func newDocket() *docket {
	return &docket{
		Data:    make(map[string]map[string][][]string),
		Current: "",
	}
}

// AddTime creates a new map for the time supplied, if not existing.
func (d *docket) AddTime(time string) {
	if _, exists := d.Data[time]; !exists {
		d.Data[time] = make(map[string][][]string)
	}
	d.Current = time
}

// AddCase parses the case title and append it to the current time.
// Should a case contain multiple people via "; ", then each will be
// appended to the current time.
func (d *docket) AddCase(cas string) {
	for _, c := range strings.Split(d.clean(cas), "; ") {
		if _, exists := d.Data[d.Current][c]; !exists {
			d.Data[d.Current][c] = [][]string{}
		}
	}
}

// AddCharge will either append a charge to the case or if the charge
// already exists for the case, it will increase it's count.
func (d *docket) AddCharge(cas string, charge string) {
	ncas := d.clean(cas)
	ncrg := d.clean(charge)
	if ncrg == "" {
		// No charge listed, use a placeholder.
		ncrg = "---"
	} else {
		// Parse the charge, removing the article.
		ncrg = strings.Split(ncrg, "] ")[1]
	}
	// Check if the charge already exists.
	idx := -1
	for i, v := range d.Data[d.Current][ncas] {
		if v[0] == ncrg {
			idx = i
			break
		}
	}
	if idx == -1 {
		// No previous charge for this case, append with a count of 1.
		d.Data[d.Current][ncas] = append(d.Data[d.Current][ncas], []string{ncrg, "1"})
	} else {
		// Previous charge exists for this case, increase the count.
		curcnt, _ := strconv.Atoi(d.Data[d.Current][ncas][idx][1])
		d.Data[d.Current][ncas][idx][1] = strconv.Itoa(curcnt + 1)
	}
}

// clean removes linebreaks, non-breaking spaces, and multiple spaces
// that appear from parsing the HTML... we just want the text.
func (d *docket) clean(s string) string {
	cnd := strings.Replace(strings.TrimSuffix(s, "\n"), "\u00A0", "", 1)
	return strings.Replace(cnd, "  ", " ", 1)
}

// fetchDocument runs an HTTP call to the court docket page for the
// current date and supplied office. Returning an HTTP response.
func fetchDocument(office string) (*http.Response, error) {
	qs := url.Values{}
	qs.Add("date", time.Now().Format("2006-01-02"))
	qs.Add("days_to_display", "1")
	qs.Add("office[]", office)

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://docket.court.nl.ca/", nil)
	req.URL.RawQuery = qs.Encode()
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("expected 200 status, got %d status", res.StatusCode)
	}
	return res, nil
}

func main() {
	res, err := fetchDocument(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	dkt := newDocket()
	tr := doc.Find("table tr")
	tr.Each(func(i int, s *goquery.Selection) {
		cr := s.Children()
		if cr.Length() == 2 {
			// This is a case row.
			dkt.AddTime(cr.Eq(0).Find("span").Text())
			dkt.AddCase(strings.Replace(cr.Eq(0).Text(), dkt.Current, "", 1)) // Remove the time to get the names.
		} else if cr.Length() == 3 {
			// This is a charge row.
			name, crg := cr.Eq(0).Text(), cr.Eq(1).Text()
			dkt.AddCharge(name, crg)
		}
	})

	// Pretty output of JSON.
	x, _ := json.MarshalIndent(dkt.Data, "", "    ")
	fmt.Println(string(x))
}
