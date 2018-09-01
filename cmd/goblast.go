package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/DanielSchuette/goblast"
)

var (
	querySeq = flag.String("query", "", "query sequence which can be an accession or GI identifier or FASTA formatted\nthis argument cannot be empty and must be a valid sequence for the API call to be successful")
)

func main() {
	// parse command line arguments
	flag.Parse()

	// check validity of input
	if *querySeq == "" {
		fmt.Println("the argument 'query' cannot be empty (see --help)")
		os.Exit(1)
	}

	// query parameters
	params := &goblast.BlastParams{
		Query:   *querySeq,
		DB:      "nt",
		Program: "blastn",
		CMD:     "Put",
		Format:  "Text",
	}

	// make a request and parse the results
	resp, err := goblast.Request(params)
	if err != nil {
		log.Fatalf("error while making an API request: %v\n", err)
	}
	rid, err := goblast.ParseResponse(resp)
	if err != nil {
		log.Fatalf("error while parsing the API response: %v\n", err)
	}
	fmt.Printf("parsed RID: %v\n", rid)

	// print the results web page for the user
	goblast.GetResultsByRID(rid)
}
