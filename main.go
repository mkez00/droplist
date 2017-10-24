package main

import (
	"log"
	"os"
	"net/http"
	"io/ioutil"
	"bufio"
	"strings"
	"os/exec"
	"bytes"
	"fmt"
)

const DropListUrl = "https://www.spamhaus.org/drop/drop.txt";

func main() {
	log.Println("Starting Blacklist import")

	// if user is not overriding drop list url use default
	dropListUrl := os.Getenv("DROPLISTURL")
	if len(dropListUrl)==0 {
		dropListUrl = DropListUrl
	}

	// store response from request
	log.Println("Using drop list URL: " + dropListUrl)
	response, err := http.Get(dropListUrl)
	check(err)
	defer response.Body.Close()

	// create temp file to store drop list
	file, err := ioutil.TempFile(os.TempDir(), "droplist")
	check(err)
	defer os.Remove(file.Name())

	// write response body to temp file
	body, err := ioutil.ReadAll(response.Body)
	ioutil.WriteFile(file.Name(), body, 0644)

	// iterate through lines of file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		curr := scanner.Text()
		//ignore description text which all begins with ;
		if curr[0:1]!=";" {
			//split string with ; first entry is IP to block
			splits := strings.Split(curr,";")
			ip := strings.TrimSpace(splits[0])

			// configure deny command for ufw and execute
			log.Println("Blocking IP: '" + ip + "'")
			cmd := exec.Command("ufw","deny", "from", ip)
			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			err := cmd.Run()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
				return
			}
			log.Println("Rule created for IP: " + ip)
		}
	}
	log.Println("Finished Blacklist import")
}

func check(err error){
	if err!=nil{
		log.Fatal(err)
		panic(err)
	}
}