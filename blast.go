package goblast

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// the following package level constants are used to enforce compliance with the NCBI Usage Guidelines
// at https://blast.ncbi.nlm.nih.gov/Blast.cgi?CMD=Web&PAGE_TYPE=BlastDocs&DOC_TYPE=DeveloperInfo
const (
	// TimeOutAfterRequest determines how long to wait after every API call
	TimeOutAfterRequest = 10
	// RequestFrequency determines the frequency with which potential results will be requested after every API call
	RequestFrequency = 10
)

// BlastParams TODO: documentation
type BlastParams struct {
	Query   string
	DB      string
	Program string
	CMD     string
	Format  string
}

// Request makes an API call to the NCBI BLAST API with a set of `params'
func Request(params *BlastParams) (*http.Response, error) {
	// check validity of input

	// construct a request URL from `params' and return an error if a parameter is missing
	url := fmt.Sprintf("https://blast.ncbi.nlm.nih.gov/Blast.cgi?QUERY=%s&DATABASE=%s&PROGRAM=%s&CMD=%s&FORMAT_TYPE=%s", params.Query, params.DB, params.Program, params.CMD, params.Format)

	// make an http request using the `url'
	log.Printf("requesting %v\nplease wait...\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error requesting %v: %v", url, err)
	}
	time.Sleep(TimeOutAfterRequest * time.Second)
	return resp, nil
}

// ParseResponse takes TODO: documentation
func ParseResponse(resp *http.Response) error {
	// if the server does not respond with a status 200 OK, return the status code
	if resp.StatusCode != 200 {
		return fmt.Errorf("server sent an invalid response back: %v", resp.Status)
	}

	// read the response body with a certain frequency and return the result to the user
	defer resp.Body.Close()
	for {
		time.Sleep(RequestFrequency * time.Second)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error while reading the response body: %v", err)
		}
		fmt.Printf("response:\n%v\n", resp)
		fmt.Printf("body:\n%v\n", string(body))
	}
	return nil
}
