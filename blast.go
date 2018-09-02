package goblast

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
	// ErrNoRID indicates a missing response identifier in a function call or API response
	ErrNoRID = errors.New("no RID found or provided")
)

// RID is of base type string and represents a response ID
// a RID can be used to retrieve results via the API
type RID string

// BlastParams holds all possible API  parameters and is used in func `Request' to construct
// the request URL
type BlastParams struct {
	// query sequence which can be an accession or
	// GI identifier or FASTA formatted
	Query string

	// database to query, defaults to `nt'
	DB string

	// the BLAST program to use, one of `blastn', `megablast', `blastp',
	// `blastx', `tblastn', or `tblastx'
	Program string

	// the server request type, one of `PUT' or `GET'
	CMD string

	// the response format, one of `HTML', `Text', `XML', `XML2',
	// `JSON2', or `Tabular'
	Format string

	// whether to use filtering, use `F' to disable, `T' or `L' to enable
	// and prepend `m' for mask at lookup (e.g., `mL')
	Filter string

	// a number greater than zero indicating the expected value
	Expect float64

	// reward for matching bases (BLASTN and megaBLAST), an integer greater
	// than zero
	Reward int

	// cost mismatched bases (BLASTN and megaBLAST), an integer less than zero
	Penalty int

	// gap existence and extension costs, space-separated positive integers
	// (e.g., '11 1')
	GapCost string

	// scoring matrix to use, one of `BLOSUM45', `BLOSUM50', `BLOSUM62' (default),
	// `BLOSUM80', `BLOSUM90', `PAM250', `PAM30' or `PAM70'
	ScoringMatrix string

	// number of database sequences to keep, an integer greater than zero
	HitListSize int

	// number of descriptions to print, an integer greater than zero
	NumberOfDescriptions int

	// number of alignments to print, an integer greater than zero
	Alignments int

	// whether or not to show NCBI GIs in report string (`T' or `F')
	NCBIGI string

	// Neighboring score for initial words, an integer greater than zero (BLASTP default is 11,
	// this value does not apply to BLASTN or MegaBLAST)
	Threshold int

	// size of word for initial matches, an integer greater than zero
	WordSize int

	// Composition based statistics algorithm to use, one of `0', `1', `2', or `3'
	CompBasedStats int

	// report formatting, one of `SearchInfo' (status check) or `Alignment'
	FormatObject string

	// response ID that can be used in a GET request to retrieve query results
	RID RID
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

	// parse the `RID' from the response body if a `name="RID" value="<RID>"` form field exists
	phrase := `name="RID" value=`
	if strings.Contains(string(body), phrase) {
		idx := strings.Index(string(body), phrase)
		if idx == -1 {
			return "", fmt.Errorf("error while parsing the response body: %v", ErrNoRID)
		}
		idx += len(phrase) /* when parsing, the actual `phrase' is going to be skipped */
		var rid []byte
		var parse bool
	ParseLoop:
		for {
			if body[idx] == '"' {
				switch parse {
				case true:
					break ParseLoop
				case false:
					parse = true
					idx++
					continue ParseLoop
				}
			}
			if parse {
				rid = append(rid, body[idx])
			}
			idx++
		}
		return RID(string(rid)), nil
	}
	return "", fmt.Errorf("error while parsing the response body: %v", ErrNoRID)
}

// GetResultsByRID takes a response ID (`RID') and prints the URL that lists the
// BLAST results associated with that `RID'
func GetResultsByRID(rid RID) {
	resultURL := fmt.Sprintf(baseURL+"?CMD=Get&RID=%s", rid)
	fmt.Printf("visit %s to see the results of your NCBI BLAST query\n", resultURL)
}

// TODO: re-write GetResultsByRID --> GetResultsURL
// then, allow for different data retrieval formats (JSON, plain text, ...)
// and via command line flags in `goblast', allow user
// to specify how to get data (as JSON, ... download or via visiting the URL)
// importantly, create two independent data structures for get request parameters
// and POST request parameters
