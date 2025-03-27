package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/term"
)

const (
	banner = `
⠀⠀⠀⠀⠀⠀⢰⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡄⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⢸⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡇⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⢸⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡇⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⣾⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡇⠀⠀⠀⠀⠀⠀
⡇⠀⠀⠀⠀⠀⢸⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡇⠀⠀⠀⠀⠀⢀
⡇⠀⠀⠀⠀⠀⢨⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠠⡃⠀⠀⠀⠀⠀⠘
⢰⠀⠀⠀⠀⠀⢰⣇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣸⡆⠀⠀⠀⠀⠀⡇
⢸⡄⠀⠀⠀⠀⠀⣿⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣿⠀⠀⠀⠀⠀⢠⠇
⠘⣧⠀⠀⠀⠀⠀⢸⣇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⡇⠀⠀⠀⠀⠀⣼⠀
⠀⠹⣆⠀⠀⠀⠀⠀⣿⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣿⠀⠀⠀⠀⠀⣰⠏⠀
⠀⠀⠹⣧⠀⠀⠀⠀⠸⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⡏⠀⠀⠀⠀⣰⠏⠀⠀
⠀⠀⠀⠹⣧⠀⠀⠀⠀⠹⣷⡀⠀⠀⠀⠀⠀⠀⢀⣾⠍⠀⠀⠀⠀⣴⠏⠀⠀⠀
⠀⠀⠀⠀⠙⡧⣀⠀⠀⠀⠘⣿⡄⠀⠀⠀⠀⢠⣾⠏⠀⠀⠀⣀⣼⠏⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠈⠙⠻⣶⣤⡀⠘⢿⡄⣀⣀⢠⣿⠃⠀⣠⣴⡾⠛⠁⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⠻⢷⣜⣿⣿⣿⣿⣣⣶⠿⠋⠁⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣠⣤⣽⣿⣿⣿⣿⣯⣅⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢀⣤⣴⠾⠿⠛⢋⣥⣿⣿⣿⣿⣿⣿⣍⠛⠻⠿⢶⣤⣄⡀⠀⠀⠀⠀
⠀⠀⠀⢰⡟⠉⠀⠀⠀⣠⡾⣻⢟⣥⣶⣿⣿⣿⡿⣷⣄⠀⠀⠈⠀⢿⡄⠀⠀⠀
⠀⠀⢠⡟⠀⠀⠀⣠⡾⠋⢰⣯⣾⣿⣿⣿⣿⣿⣿⡈⠻⣷⣄⠀⠀⠈⢷⡀⠀⠀
⠀⢀⡾⠁⠀⠀⣼⠋⠀⠀⢸⢸⣿⡿⠿⣿⠿⣿⣿⡇⠀⠈⢫⣧⠀⠀⠘⣷⠀⠀
⠀⣼⠃⠀⠀⢠⣿⠀⠀⠀⠸⣿⣿⣿⡆⠀⣼⡟⣹⠀⠀⠀⠀⣿⠀⠀⠀⠸⣧⠀
⠀⡟⠀⠀⠀⢸⡏⠀⠀⠀⠀⠙⢿⣯⣶⣶⣮⡿⠃⠀⠀⠀⠀⢹⡇⠀⠀⠀⣿⠀
⠀⡇⠀⠀⠀⣼⠇⠀⠀⠀⠀⠀⠀⠉⠛⠋⠉⠀⠀⠀⠀⠀⠀⢸⣇⠀⠀⠀⢸⠀
⠀⡇⠀⠀⠀⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⣿⠀⠀⠀⢸⠀
⠀⡇⠀⠀⠀⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⠀⠀⠀⢸⠀
⠀⡇⠀⠀⠀⢸⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⡏⠀⠀⠀⢸⠀
⠀⠁⠀⠀⠀⠀⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣸⠁⠀⠀⠀⠈⠀
⠀⠀⠀⠀⠀⠀⠸⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡇⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⢧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡸⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠈⠄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⠁⠀⠀⠀⠀⠀⠀⠀
`

	CYAN          = "\033[36m"
	RED           = "\033[91m"
	GREEN         = "\033[92m"
	YELLOW        = "\033[93m"
	RESET         = "\033[0m"
	author        = "Loca Martin"
	baseGoogleURL = "https://www.google.com/search?q="
	userAgent     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
	resultSelector = "div.g"
	nextPageSelector = "a#pnnext"
)

var (
	site      = flag.String("site", "", "Single domain to search")
	siteList  = flag.String("site-list", "", "File containing list of domains")
	dorkList  = flag.String("dork-list", "", "File containing list of dorks")
	pages     = flag.Int("p", 10, "Number of pages to fetch")
	threads   = flag.Int("t", 10, "Number of concurrent threads")
	output    = flag.String("o", "", "Output file name")
	silent    = flag.Bool("silent", false, "Suppress banner output")
	version   = flag.Bool("version", false, "Show version information")
	showHelp  = flag.Bool("h", false, "Show help message")
)

type SearchTask struct {
	Domain string
	Dork   string
}

type GoogleSearcher struct {
	client *http.Client
	delay  time.Duration
}

func NewGoogleSearcher() *GoogleSearcher {
	return &GoogleSearcher{
		client: &http.Client{Timeout: 15 * time.Second},
		delay:  3 * time.Second,
	}
}

func (gs *GoogleSearcher) Search(ctx context.Context, query string, maxPages int) ([]string, error) {
	var results []string
	currentPage := 0

	for currentPage < maxPages {
		select {
		case <-ctx.Done():
			return results, nil
		default:
			pageResults, nextPage, err := gs.fetchPage(query, currentPage*10)
			if err != nil {
				return results, err
			}
			results = append(results, pageResults...)
			
			if nextPage == "" {
				return results, nil
			}
			query = nextPage
			currentPage++
			time.Sleep(gs.delay)
		}
	}
	return results, nil
}

func (gs *GoogleSearcher) fetchPage(query string, start int) ([]string, string, error) {
	searchURL := fmt.Sprintf("%s%s&start=%d", baseGoogleURL, url.QueryEscape(query), start)
	
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := gs.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, "", fmt.Errorf("google returned status: %d", resp.StatusCode)
	}

	return parseSearchResults(resp.Body)
}

func parseSearchResults(body io.Reader) ([]string, string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, "", err
	}

	var results []string
	doc.Find(resultSelector).Each(func(i int, s *goquery.Selection) {
		link := s.Find("a[href]")
		if href, exists := link.Attr("href"); exists {
			if strings.HasPrefix(href, "/url?q=") {
				cleanURL := strings.Split(href, "&sa=U&")[0]
				cleanURL = strings.TrimPrefix(cleanURL, "/url?q=")
				results = append(results, cleanURL)
			}
		}
	})

	nextPage, _ := doc.Find(nextPageSelector).Attr("href")
	return results, nextPage, nil
}

func printCenteredBanner() {
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))
	if width < 80 {
		width = 80
	}

	lines := strings.Split(strings.Trim(banner, "\n"), "\n")
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	padding := (width - maxLen) / 2
	if padding < 0 {
		padding = 0
	}

	fmt.Printf("%s", CYAN)
	for _, line := range lines {
		fmt.Printf("%s%s\n", strings.Repeat(" ", padding), line)
	}
	fmt.Printf("%sAuthor: %s\n\n%s", strings.Repeat(" ", padding+10), author, RESET)
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	return lines, scanner.Err()
}

func worker(id int, tasks <-chan SearchTask, wg *sync.WaitGroup, outputFile *os.File, searcher *GoogleSearcher) {
	defer wg.Done()
	for task := range tasks {
		query := fmt.Sprintf("%s site:%s", task.Dork, task.Domain)
		results, err := searcher.Search(context.Background(), query, *pages)
		if err != nil {
			fmt.Printf("%s[ERROR]%s %v\n", RED, RESET, err)
			continue
		}

		for _, result := range results {
			outputLine := fmt.Sprintf("%s[]%s %s\n", GREEN, RESET, result)
			fmt.Print(outputLine)
			if outputFile != nil {
				outputFile.WriteString(result + "\n")
			}
		}
	}
}

func main() {
	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *version {
		fmt.Println("Dork v1.0")
		os.Exit(0)
	}

	if !*silent {
		printCenteredBanner()
	}

	// Handle Ctrl+C
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Printf("\n%s[!] Exiting...%s\n", RED, RESET)
		os.Exit(1)
	}()

	// Read domains
	var domains []string
	if *site != "" {
		domains = []string{*site}
	} else if *siteList != "" {
		var err error
		domains, err = readLines(*siteList)
		if err != nil {
			fmt.Printf("%s[ERROR]%s Failed to read site list: %v\n", RED, RESET, err)
			os.Exit(1)
		}
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				domains = append(domains, strings.TrimSpace(scanner.Text()))
			}
		}
	}

	if len(domains) == 0 {
		fmt.Printf("%s[ERROR]%s No domains specified\n", RED, RESET)
		os.Exit(1)
	}

	// Read dorks
	var dorks []string
	if *dorkList != "" {
		var err error
		dorks, err = readLines(*dorkList)
		if err != nil {
			fmt.Printf("%s[ERROR]%s Failed to read dork list: %v\n", RED, RESET, err)
			os.Exit(1)
		}
	} else {
		if term.IsTerminal(int(os.Stdin.Fd())) {
			fmt.Printf("%sdork-shell%%%s ", YELLOW, RESET)
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				dork := strings.TrimSpace(scanner.Text())
				if dork != "" {
					dorks = append(dorks, dork)
				}
				fmt.Printf("%sdork-shell%%%s ", YELLOW, RESET)
			}
		} else {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				dorks = append(dorks, strings.TrimSpace(scanner.Text()))
			}
		}
	}

	if len(dorks) == 0 {
		fmt.Printf("%s[ERROR]%s No dorks specified\n", RED, RESET)
		os.Exit(1)
	}

	// Prepare output file
	var outputFile *os.File
	if *output != "" {
		var err error
		outputFile, err = os.Create(*output)
		if err != nil {
			fmt.Printf("%s[ERROR]%s Failed to create output file: %v\n", RED, RESET, err)
			os.Exit(1)
		}
		defer outputFile.Close()
	}

	// Initialize search engine
	searcher := NewGoogleSearcher()

	// Create worker pool
	tasks := make(chan SearchTask, *threads*2)
	var wg sync.WaitGroup

	for i := 0; i < *threads; i++ {
		wg.Add(1)
		go worker(i, tasks, &wg, outputFile, searcher)
	}

	// Generate tasks
	for _, domain := range domains {
		for _, dork := range dorks {
			tasks <- SearchTask{Domain: domain, Dork: dork}
		}
	}

	close(tasks)
	wg.Wait()
}