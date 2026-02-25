package ipfile

import (
	"bytes"
	"bufio"
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"golang.org/x/net/http2"
)

var Downloaddir string

type DownloadFile struct {
	Url              string `json:"url"`
	DownloadFileName string `json:"DownloadFileName"`
	Cloudplatform    string `json:"Cloudplatform"`
}

type Common struct {
	Debug bool
}

type azureCandidateCollector struct {
	candidates *[]string
}

func (c azureCandidateCollector) Collect(_ int, s *goquery.Selection) {
	href, ok := s.Attr("href")
	if !ok || href == "" {
		return
	}
	lower := strings.ToLower(href)
	if strings.Contains(lower, "download.microsoft.com") && (strings.HasSuffix(lower, ".json") || strings.HasSuffix(lower, ".zip")) {
		*c.candidates = append(*c.candidates, href)
	}
}

func checkResponse( StatusCode int) bool{
	if StatusCode >= 200 && StatusCode < 300 {
		return true
	}
	return false
}

func (i *Common) Download(DownloadFileName string, Url string) (err error) {
	//function to download the cloud ip
	log.Printf("Downloading %s to %s", Url, DownloadFileName)
	start := time.Now()
	if i != nil && i.Debug {
		log.Printf("Download debug enabled")
	}
	//Download the IP Address file``
	// Create the file
	fileOut, err := os.Create(DownloadFileName)
	if err != nil {
		return err
	}
	defer fileOut.Close()

	resp, err := http.Get(Url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch error: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	if i != nil && i.Debug {
		log.Printf("Response status: %s", resp.Status)
		if resp.ContentLength >= 0 {
			log.Printf("Content length: %d", resp.ContentLength)
		}
	}

	//Check server response
	if !checkResponse(resp.StatusCode) {
		err = fmt.Errorf("bad respose status: %v", resp.StatusCode)
		return err
	}

	// Write the body to file
	bytesWritten, err := io.Copy(fileOut, resp.Body)
	if err != nil {
		return err
	}

	if i != nil && i.Debug {
		log.Printf("Wrote %d bytes in %s", bytesWritten, time.Since(start))
	}
	return nil
}

func WriteFile(DownloadFileName string, cidrs []string ) {
	f, err := os.Create(DownloadFileName)
    if err != nil {
        log.Fatal(err)
    }
    // remember to close the file
    defer f.Close()

    for _, cidr := range cidrs {
        _, err := f.WriteString(cidr + "\n")
        if err != nil {
            log.Fatal(err)
        }
    }
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
	}
	
	// Detect IP type and add appropriate suffix
	if strings.Contains(str_in, ":") {
		// IPv6 address
		return fmt.Sprintf("%s/128", str_in)
	}
	
	// IPv4 address
	return fmt.Sprintf("%s/32", str_in)
}

func MatchIp(pattern string) []string {
	//match ip addresses from string pattern and return slice of ips as string
	// Match IPv4 addresses only
	reIPv4 := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(?:/\d{1,2}|)`)
	result := reIPv4.FindAllString(pattern, -1)
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
	byteValue, _ := io.ReadAll(jsonFile)
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

func ResolveAzureDownloadUrl(debug bool) (string, error) {
	//Extract the dynamic download URL from the service tag published page
	var link string

	downloadPage := "https://www.microsoft.com/en-us/download/confirmation.aspx?id=56519"
	if debug {
		log.Printf("Azure download page: %s", downloadPage)
	}

	tlsTransport := &http2.Transport{
        TLSClientConfig: &tls.Config{
            MaxVersion: tls.VersionTLS12,
        },
    }

	client := http.Client{Transport: tlsTransport}
	req , err := http.NewRequest("GET", downloadPage, nil)
	if err != nil {
		return link, err
	}

	resp , err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Azure status code error: %d %s\n", resp.StatusCode, resp.Status)
	} else if debug {
		log.Printf("Azure status: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyBytes))

	if err != nil {
		log.Fatal(err)
		return "error", err
	}

	var candidates []string
	collector := azureCandidateCollector{candidates: &candidates}
	doc.Find("a").Each(collector.Collect)

	if debug {
		log.Printf("Azure download candidates: %d", len(candidates))
	}

	if len(candidates) > 0 {
		link = candidates[0]
	} else {
		re := regexp.MustCompile(`https?://download\.microsoft\.com/[^"'\s]+\.(json|zip)`)
		matches := re.FindAllString(string(bodyBytes), -1)
		if debug {
			log.Printf("Azure regex candidates: %d", len(matches))
		}
		if len(matches) > 0 {
			link = matches[0]
		}
	}

	if debug {
		if link == "" {
			log.Printf("Azure download link not found")
		} else {
			log.Printf("Azure download link: %s", link)
		}
	}

	return link, nil
}

