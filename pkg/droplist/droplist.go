package droplist

import (
	"os"
	"log"
	"net/http"
	"bufio"
	"strings"
	"os/exec"
	"bytes"
	"fmt"
	"io/ioutil"
)

const DropListUrl = "https://www.spamhaus.org/drop/drop.txt";

func FetchAndApply(){
	response := GetDropList()
	defer response.Body.Close()

	file := CopyResponseIntoFile(response)
	defer os.Remove(file.Name())

	// iterate through lines of file
	IterateAndDeny(file)
}

func GetDropList() *http.Response {
	// if user is not overriding drop list url use default
	dropListUrl := os.Getenv("DROPLISTURL")
	if len(dropListUrl)==0 {
		dropListUrl = DropListUrl
	}

	// store response from request
	log.Println("Using drop list URL: " + dropListUrl)
	response, err := http.Get(dropListUrl)
	check(err)
	return response
}

func IterateAndDeny(file *os.File){
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
}

func CopyResponseIntoFile(response *http.Response) *os.File{
	// create temp file to store drop list
	file, err := ioutil.TempFile(os.TempDir(), "droplist")
	check(err)
	// write response body to temp file
	body, err := ioutil.ReadAll(response.Body)
	ioutil.WriteFile(file.Name(), body, 0644)
	return file
}

func check(err error){
	if err!=nil{
		log.Fatal(err)
		panic(err)
	}
}