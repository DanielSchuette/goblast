package main

import (
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
	err = goblast.ParseResponse(resp)
	if err != nil {
		log.Fatalf("error while parsing the API response: %v\n", err)
	}
}
