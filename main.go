package main

import (
	"log"
	"github.com/mkez00/droplist/pkg/droplist"
)

func main() {
	log.Println("Starting Droplist import")
	droplist.FetchAndApply()
	log.Println("Finished Droplist import")
}