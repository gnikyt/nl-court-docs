package ncd

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Alias type for the docket data container.
type DocketMapping = map[string]map[string][]Charge

// Charge represents a charge for a case.
type Charge struct {
	Description string // description of the charge.
	Count       int    // number of counts for the charge.
}

// Increase will increase the count by 1.
func (c *Charge) Increase() {
	c.Count += 1
}

// HasMultiple checks if there are multiple counts for the charge.
func (c Charge) HasMultiple() bool {
	return c.Count > 1
}

// Docket represents the data parsed from the Court Docket.
type Docket struct {
	Data        DocketMapping // format of... { [time]: [{ [name]: [{ Description: [charge], Count: [count] }] } }.
	Office      string        // which office for the docket.
	Date        string        // which date for the docket.
	Client      *http.Client  // HTTP client handler.
	currentTime string        // current time slot.
}

// NewDocket inits a Docket struct.
func NewDocket(date string, office string, client *http.Client) *Docket {
	d := &Docket{
		currentTime: "",
		Data:        make(map[string]map[string][]Charge),
		Office:      office,
		Date:        date,
	}
	if client == nil {
		d.Client = &http.Client{}
	}
	return d
}

// AddTime creates a new map for the time supplied, if not existing.
func (d *Docket) AddTime(time string) {
	if _, exists := d.Data[time]; !exists {
		d.Data[time] = make(map[string][]Charge)
	}
	d.currentTime = time
}

// AddCase parses the case title and append it to the current time.
// Should a case contain multiple people via "; ", then each will be
// appended to the current time.
func (d *Docket) AddCase(cas string) {
	for _, c := range strings.Split(d.cleanString(cas), "; ") {
		if _, exists := d.Data[d.currentTime][c]; !exists {
			d.Data[d.currentTime][c] = make([]Charge, 0)
		}
	}
}

// AddCharge will either append a charge to the case or if the charge
// already exists for the case, it will increase it's count.
func (d *Docket) AddCharge(cas string, charge string) {
	ncas := d.cleanString(cas)
	ncrg := d.cleanString(charge)
	if ncrg == "" {
		// No charge listed, use a placeholder.
		ncrg = "(empty)"
	} else {
		// Parse the charge, removing the article.
		ncrg = strings.Split(ncrg, "] ")[1]
	}
	// Check if the charge already exists.
	idx := -1
	for i, v := range d.Data[d.currentTime][ncas] {
		if v.Description == ncrg {
			idx = i
			break
		}
	}
	if idx == -1 {
		// No previous charge for this case, append with a count of 1.
		d.Data[d.currentTime][ncas] = append(
			d.Data[d.currentTime][ncas],
			Charge{Description: ncrg, Count: 1},
		)
	} else {
		// Previous charge exists for this case, increase the count.
		d.Data[d.currentTime][ncas][idx].Increase()
	}
}

// Fetch runs an HTTP call to the court docket page for the
// current date and supplied office. Returning an HTTP response.
func (d *Docket) Fetch() (*http.Response, error) {
	qs := url.Values{}
	qs.Add("date", d.Date)
	qs.Add("days_to_display", "1")
	qs.Add("office[]", d.Office)

	req, err := http.NewRequest("GET", "https://docket.court.nl.ca/", nil)
	req.URL.RawQuery = qs.Encode()
	if err != nil {
		return nil, err
	}
	res, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("expected 200 status, got %d status", res.StatusCode)
	}
	return res, nil
}

// Parse takes the HTTP response and parses it, saving into the data field.
func (d *Docket) Parse(res *http.Response) error {
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}
	tr := doc.Find("table tr")
	tr.Each(func(i int, s *goquery.Selection) {
		cr := s.Children()
		if cr.Length() == 2 {
			// This is a case row.
			d.AddTime(cr.Eq(0).Find("span").Text())
			d.AddCase(strings.Replace(cr.Eq(0).Text(), d.currentTime, "", 1)) // Remove the time to get the names.
		} else if cr.Length() == 3 {
			// This is a charge row.
			name, crg := cr.Eq(0).Text(), cr.Eq(1).Text()
			d.AddCharge(name, crg)
		}
	})
	return nil
}

// Output accepts an Outputter and returns a string.
func (d *Docket) Output(out Outputter) (string, error) {
	res, err := out.Format()
	if err != nil {
		return "", err
	}
	return res.String(), nil
}

// clean removes linebreaks, non-breaking spaces, and multiple spaces
// that appear from parsing the HTML... we just want the text.
func (d *Docket) cleanString(s string) string {
	cnd := strings.Replace(strings.TrimSuffix(s, "\n"), "\u00A0", "", 1)
	return strings.Replace(cnd, "  ", " ", 1)
}
