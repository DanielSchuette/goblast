package main

import (
	"fmt"
	"log"

	"github.com/DanielSchuette/goblast"
)

func main() {
	// query parameters
	params := &goblast.BlastParams{
		Query:   "gagtctcctttggaactctgcaggttctatttgctttttcccagatgagctctttttctggtgtttgtct",
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
