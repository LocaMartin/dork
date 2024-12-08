package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/LocaMartin/dork"
)

func main() {
	// Define flags
	domainsFile := flag.String("domains", "domains.txt", "Path to the file containing domains")
	inurlsFile := flag.String("inurls", "inurls.txt", "Path to the file containing inurls")
	outputFile := flag.String("output", "output.txt", "Path to the output file")

	flag.Parse()

	// Read domains from file
	domainsFileHandle, err := os.Open(*domainsFile)
	if err != nil {
		log.Fatalf("Error opening domains file: %v", err)
	}
	defer domainsFileHandle.Close()

	domains := []string{}
	scanner := bufio.NewScanner(domainsFileHandle)
	for scanner.Scan() {
		domains = append(domains, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading domains file: %v", err)
	}

	// Read inurls from file
	inurlsFileHandle, err := os.Open(*inurlsFile)
	if err != nil {
		log.Fatalf("Error opening inurls file: %v", err)
	}
	defer inurlsFileHandle.Close()

	inurls := []string{}
	scanner = bufio.NewScanner(inurlsFileHandle)
	for scanner.Scan() {
		inurls = append(inurls, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading inurls file: %v", err)
	}

	// Create output file
	outputFileHandle, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outputFileHandle.Close()

	// Generate URLs using Dork method and write to file
	for _, domain := range domains {
		for _, inurl := range inurls {
			url := fmt.Sprintf("site:%s %s\n", domain, inurl)
			_, err := outputFileHandle.WriteString(url)
			if err != nil {
				log.Fatalf("Error writing to output file: %v", err)
			}
		}
	}

	fmt.Printf("URLs generated successfully and saved to %s\n", *outputFile)
}
