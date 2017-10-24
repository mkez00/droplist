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
	"io"
)

const DropListUrl = "https://www.spamhaus.org/drop/drop.txt";
const AppFileLoc = "/var/droplist/"
const CheckFileName = ".checkfile"
const CurrFileName = ".curr"

// Function that completes the task of fetching the drop list and applying it to the firewall
func FetchAndApply(){

	// Get drop list from URL as HTTP response
	response := GetDropList()
	defer response.Body.Close()

	// Copy the HTTP response body into a temp file
	file := CopyResponseIntoFile(response)
	defer os.Remove(file.Name())

	//Remove old list
	//RemoveOldList()

	// Take the temp file and create .curr file for processing
	CleanFileSpamHaus(file)

	// Iterate through .curr file adding deny entries to firewall
	IterateAndDeny()

	// copy new list into old
	CopyNewFileIntoOld()
}

// Fetches the droplist from the URL provided via DROPLISTURL env variable or from default location
func GetDropList() *http.Response {
	// if user is not overriding drop list url use default
	dropListUrl := GetDropListUrl()

	// store response from request
	log.Println("Using drop list URL: " + dropListUrl)
	response, err := http.Get(dropListUrl)
	check(err)
	return response
}

// Takes IP addresses written to .curr file and applies firewall deny rule
func IterateAndDeny(){
	file, err := os.Open(GetAppFileLoc() +  CurrFileName)
	check(err)

	// iterate through lines of file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip := strings.TrimSpace(scanner.Text())
		log.Println("Blocking IP: '" + ip + "'")
		cmd := exec.Command("ufw","deny", "from", ip)
		executeCommand(cmd)
		log.Println("Rule created for IP: " + ip)
	}
}

// Spamhaus specific list.  Take this list and translate into proper format for IterateAndDeny() function
func CleanFileSpamHaus(file *os.File) {
	currFile, err := os.Create(GetAppFileLoc() +  CurrFileName)
	check(err)

	// iterate through lines of file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		curr := scanner.Text()
		//ignore description text which all begins with ;
		if curr[0:1]!=";" {
			//split string with ; first entry is IP to block
			splits := strings.Split(curr,";")
			ip := strings.TrimSpace(splits[0])
			currFile.WriteString(ip + "\n")
		}
	}
}

// This does not work
func RemoveOldList(){
	log.Println("Removing old block list")
	checkFile, err := os.Open(GetAppFileLoc() +  CheckFileName)
	if err!=nil {
		log.Println("First run....no previous runs")
		_, err := os.Create(GetAppFileLoc() +  CheckFileName)
		check(err)
		return
	}
	scanner := bufio.NewScanner(checkFile)
	for scanner.Scan() {
		ip := strings.TrimSpace(scanner.Text())
		cmd := exec.Command("ufw","delete", "deny", ip)
		log.Println("Removing Deny Rule on IP: " + ip)
		executeCommand(cmd)

	}
	log.Println("Removing old block list complete")
}

// Takes the current list of IP addresses appended during the most recent run and copy to the check file for next run
func CopyNewFileIntoOld(){
	in , err := os.Open(GetAppFileLoc() + CurrFileName)
	check(err)
	defer in.Close()

	out, err := os.Create(GetAppFileLoc() + CheckFileName)
	check(err)
	defer out.Close()

	_, err2 := io.Copy(out,in)
	check(err2)

	os.Remove(GetAppFileLoc() + CurrFileName)
}

// Take the HTTP response which fetched the file and add to temp file
func CopyResponseIntoFile(response *http.Response) *os.File{
	// create temp file to store drop list
	file, err := ioutil.TempFile(os.TempDir(), "droplist")
	check(err)
	// write response body to temp file
	body, err := ioutil.ReadAll(response.Body)
	ioutil.WriteFile(file.Name(), body, 0644)
	return file
}

// Gets previously run file
func GetPreviousRunFile() (string,[]byte) {
	checkFileLoc := GetAppFileLoc()
	path := checkFileLoc + CheckFileName
	file, err := ioutil.ReadFile(path)
	check(err)
	return checkFileLoc,file
}

func GetDropListUrl() string {
	return GetEnvOverriddenVal(DropListUrl, "DROPLISTURL")
}

func GetAppFileLoc() string{
	return GetEnvOverriddenVal(AppFileLoc, "APPFILELOC")
}

func GetEnvOverriddenVal(base string, key string) string{
	returnString := os.Getenv(key)
	if len(returnString)==0 {
		returnString = base
	}
	return returnString
}

func check(err error){
	if err!=nil{
		log.Fatal(err)
		panic(err)
	}
}

func executeCommand(cmd *exec.Cmd) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		panic(err)
	}
}