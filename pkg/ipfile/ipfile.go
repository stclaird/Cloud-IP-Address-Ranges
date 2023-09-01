package ipfile

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var Downloaddir string

type DownloadFile struct {
	Url              string `json:"url"`
	DownloadFileName string `json:"DownloadFileName"`
	Cloudplatform    string `json:"Cloudplatform"`
}

type Common struct {
}

func (i *Common) Download(DownloadFileName string, Url string) (err error) {
	//function to download the cloud ip
	log.Printf("Downloading %s to %s", Url, DownloadFileName)
	//Download the IP Address file``
	// Create the file
	fileOut, err := os.Create(DownloadFileName)
	if err != nil {
		return err
	}
	defer fileOut.Close()

	resp, err := http.Get(Url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		fmt.Println(fmt.Errorf("bad status: %s", resp.Status))
	}

	// Write the body to file
	_, err = io.Copy(fileOut, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

type IpfileJson struct {
	SyncToken    string `json:"syncToken"`
	CreationTime string `json:"creationTime"`
}

type IpfileTXT struct {
	Common
	Prefixes []string
}

func IPtoCidr(str_in string)(string) {
	//Ensure an IP address has a slash (cidr Notation)
	contains := strings.Contains(str_in, "/")

	if contains {
		return str_in
	} else {
		return fmt.Sprintf("%s/32", str_in)
	}
}

func MatchIp(pattern string) []string {
	//match ip addresses from string pattern and return slice of ips as string
	re := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(?:/\d{1,2}|)`)
	result := re.FindAllString(pattern, -1)

	return result
}

func StrInSlice(str string, slice []string) bool {
	//find a string in slice return boolean
	for _, val := range slice {
		if val == str {
			return true
		}
	}
	return false
}
func AsJson[T any](DownloadFileName string) (fileOut T) {
	// Open downloaded file and return as json
	jsonFile, err := os.Open(DownloadFileName)
	if err != nil {
		log.Println("Error", err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &fileOut)

	return fileOut
}

func AsCSV[T any](DownloadFileName string, column int) []string {
	// Open a CSV and retrieve CIDR
	var cidrs []string
	csvfile, err := os.Open(DownloadFileName)
	if err != nil {
		log.Println("Error", err)
	}

	r := csv.NewReader(csvfile)
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		cidrs = append(cidrs, record[column])
	}

	return cidrs
}

func AsText[T any](DownloadFileName string) []string {
	file, err := os.Open(DownloadFileName)
	if err != nil {
		log.Println("Error", err)
	}
	defer file.Close()

	var cidrs []string
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
    scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		txt := scanner.Text()
		matched := MatchIp(txt)
		for _, cidr := range matched {
			cidrs = append(cidrs, cidr)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Print(err)
	}

	return cidrs
}

func Process(cidrs_in []string) (cidrs_out []string) {
	for _, val := range cidrs_in {
		val = IPtoCidr(val)
		exists := StrInSlice(val, cidrs_out)
		if exists == false {
			cidrs_out = append(cidrs_out, val)
		}
	}

	return cidrs_out
}

func ResolveAzureDownloadUrl() string {
	//Extract the dynamic download URL from the service tag published page

	downloadPage := "https://www.microsoft.com/en-us/download/confirmation.aspx?id=56519"
	resp, err := http.Get(downloadPage)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	item := doc.Find(".mscom-link.failoverLink").First()
	link, _ := item.Attr("href")

	return link
}
