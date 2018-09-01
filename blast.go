package goblast

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// These package level constants are used to enforce compliance with the NCBI Usage Guidelines
// at https://blast.ncbi.nlm.nih.gov/Blast.cgi?CMD=Web&PAGE_TYPE=BlastDocs&DOC_TYPE=DeveloperInfo
const (
	// Timeout determines how long to wait for every API call before timing out
	Timeout = 10
	// MaxRequestFrequency determines the frequency with which API calls can be made
	MaxRequestFrequency = 10
)

var (
	client  = http.Client{Timeout: Timeout * time.Second}
	baseURL = "https://blast.ncbi.nlm.nih.gov/Blast.cgi"
)

// RID is a string representing a response ID
type RID string

// BlastParams holds all possible API  parameters and is used in func `Request' to construct the request URL
type BlastParams struct {
	Query                string  /* query sequence which can be an accession or GI identifier or FASTA formatted */
	DB                   string  /* database to query, defaults to `nt' */
	Program              string  /* one of blastn, megablast, blastp, blastx, tblastn, tblastx */
	CMD                  string  /* PUT or GET */
	Format               string  /* one of HTML, Text, XML, XML2, JSON2, or Tabular */
	Filter               string  /* 'F' to disable. 'T' or 'L' to enable. Prepend 'm' for mask at lookup (e.g., 'mL') */
	Expect               float64 /* a number greater than zero indicating the expected value */
	Reward               int     /* reward for matching bases (BLASTN and megaBLAST), an integer greater than zero */
	Penalty              int     /* cost mismatched bases (BLASTN and megaBLAST), an integer less than zero */
	GapCost              string  /* gap existence and extension costs, space-separated positive integers (e.g., '11 1') */
	ScoringMatrix        string  /* scoring matrix to use, one of BLOSUM45, BLOSUM50, BLOSUM62 (default), BLOSUM80, BLOSUM90, PAM250, PAM30 or PAM70*/
	HitListSize          int     /* number of databases sequences to keep, an integer greater than zero */
	NumberOfDescriptions int     /* number of descriptions to print, an integer greater than zero */
	Alignments           int     /* number of alignments to print, an integer greater than zero */
	NCBIGI               string  /* whether or not to show NCBI GIs in report string ('T' or 'F') */
	Threshold            int     /* Neighboring score for initial words, an integer greater than zero (BLASTP default is 11, does not apply to BLASTN or MegaBLAST) */
	WordSize             int     /* size of word for initial matches, an integer greater than zero */
	CompBasedStats       int     /* Composition based statistics algorithm to use, one of 0, 1, 2, or 3 */
	FormatObject         string  /* 'SearchInfo' (status check) or 'Alignment' (report formatting) */
	RID                  RID     /* response ID that can be used in a GET request to retrieve query results */
}

// Request makes an API call to the NCBI BLAST API with a set of `params'
func Request(params *BlastParams) (*http.Response, error) {
	// check validity of input

	// construct a request URL from `params' and return an error if a parameter is missing
	// TODO: use net/url package instead
	url := fmt.Sprintf(baseURL+"?QUERY=%s&DATABASE=%s&PROGRAM=%s&CMD=%s&FORMAT_TYPE=%s", params.Query, params.DB, params.Program, params.CMD, params.Format)

	// make an http request using the `url'
	fmt.Printf("requesting %v\nplease wait...\n", url)
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error requesting %v: %v", url, err)
	}
	time.Sleep(Timeout * time.Second)
	return resp, nil
}

// ParseResponse takes a `*http.Response' and parses and returns a `RID' (response ID) that
// can be used in an API GET request to retrieve the results of the initial query in `resp'
func ParseResponse(resp *http.Response) (RID, error) {
	// if the server does not respond with a status 200 OK, return the status code
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("server sent an invalid response back: %v", resp.Status)
	}
	time.Sleep(MaxRequestFrequency * time.Second) /* sleep for some time before reading response */

	// read the response body with a certain frequency and return the result to the user
	body, err := ioutil.ReadAll(resp.Body)
	defer func(*http.Response) {
		err = resp.Body.Close()
		if err != nil {
			log.Fatalf("error closing response body: %v\n", err)
		}
	}(resp)
	if err != nil {
		return "", fmt.Errorf("error while reading the response body: %v", err)
	}

	// parse the `RID' from the response body
	for i, n := 0, len(body); i < n; i++ {
		if string(body[i:(i+4)]) == "name" && (string(body[(i+6):(i+9)]) == "RID") && (string(body[(i+11):(i+16)]) == "value") {
			j := i + 18
			var rid []byte
			for {
				rid = append(rid, body[j])
				j++
				if string(body[j:(j+1)]) == "\"" {
					break
				}
			}
			return RID(string(rid)), nil
		}
	}
	return "", errors.New("error while parsing the response body: no 'RID' found")
}

// GetResultsByRID takes a response ID (`RID') and prints the URL that lists the
// BLAST results associated with that `RID'
func GetResultsByRID(rid RID) {
	resultURL := fmt.Sprintf(baseURL+"?CMD=Get&RID=%s", rid)
	fmt.Printf("visit %s to see the results of your NCBI BLAST query\n", resultURL)
}
