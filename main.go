package main

import (
	"log"
	"github.com/mkez00/droplist/pkg/droplist"
)

func main() {
	log.Println("Starting Blacklist import")
	droplist.FetchAndApply()
	log.Println("Finished Blacklist import")
}