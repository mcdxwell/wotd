package wotd

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const wotdBaseURL = "https://www.merriam-webster.com/word-of-the-day/"

var seededRand *rand.Rand

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	seededRand = rand.New(source)
}

type Word struct {
	Date       string `json:"date"`
	Word       string `json:"word"`
	Class      string `json:"class,omitempty"`
	Meaning    string `json:"meaning,omitempty"`
	Definition string `json:"definition,omitempty"`
	Example    string `json:"example,omitempty"`
}

func Wotd() string {
	today := time.Now().UTC().Format("2006-01-02")
	word := fetchWord(today, getDatedURL(today))
	return word
}

func RandomWord() string {
	date := getRandomDate()
	rndURL := getDatedURL(date)
	randomWord := fetchWord(date, rndURL)
	return randomWord
}

func Link(word string) string {
	return fmt.Sprintf("https://www.merriam-webster.com/dictionary/%v", word)
}

// fetchWord first checks the cache, then fetches if necessary.
func fetchWord(date, url string) string {
	// 1. Check cache (~/.cache/wotd/wotd.json)
	words, err := searchCache()
	if err != nil {
		// Log the cache read failure but continue, as we can still fetch the word.
		fmt.Fprintf(os.Stderr, "Warning: Failed to read cache, proceeding with fetch: %v\n", err)
		words = []Word{} // Treat it as an empty cache
	}

	for _, word := range words {
		if word.Date == date {
			// Found in cache, success!
			return word.Word
		}
	}

	// 2. Fetch from Merriam-Webster
	title, err := getWordTitle(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch word title for %s: %v\n", date, err)
		return ""
	}

	// 3. Clean up the title to get the word
	word := formatTitleToWord(title)

	// 4. Save to cache (Only save if word list was successfully loaded or created)
	if err := saveInfo(date, word); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to save word info to cache: %v\n", err)
	}

	return word
}

// getRandomDate creates a random date between 2006-10-10 and the current date.
func getRandomDate() string {
	const minDateStr = "2006-10-10"
	const dateLayout = "2006-01-02"

	start, err := time.Parse(dateLayout, minDateStr)
	if err != nil {
		// Should only happen if constant is mistyped
		panic(fmt.Sprintf("Invalid hardcoded start date: %v", err))
	}

	today := time.Now().UTC()

	// Calculate the difference in days
	diff := today.Sub(start)
	days := int(diff.Hours() / 24)

	// Generate a random number of days back from today
	randomDayOffset := seededRand.Intn(days + 1) // +1 to include today

	// Subtract the offset from today
	randomDate := today.AddDate(0, 0, -randomDayOffset)

	return randomDate.Format(dateLayout)
}

// getDatedURL constructs the URL for a specific date.
func getDatedURL(date string) string {
	// Today's WOTD page is just the base URL.
	if date == time.Now().UTC().Format("2006-01-02") {
		return wotdBaseURL
	}
	// Historical words use /word-of-the-day/YYYY-MM-DD
	return wotdBaseURL + date
}

// formatTitleToWord cleans the HTML <title> tag text to extract the word.
func formatTitleToWord(title string) string {
	res := strings.TrimPrefix(title, "Word of the Day: ")
	res = strings.TrimSuffix(res, " | Merriam-Webster")
	return res
}

// getWordTitle fetches the URL and extracts the content of the <title> tag.
func getWordTitle(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to make HTTP GET request to %s: %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status code (%d) from %s", res.StatusCode, url)
	}

	return extractTitleFromHTML(res.Body)
}

// extractTitleFromHTML processes the HTML body to find the title.
func extractTitleFromHTML(r io.Reader) (string, error) {
	token := html.NewTokenizer(r)
	var isTitle bool

	for {
		tt := token.Next()
		switch tt {
		case html.ErrorToken:
			err := token.Err()
			if err == io.EOF {
				return "", fmt.Errorf("title tag not found before EOF")
			}
			return "", fmt.Errorf("HTML tokenizing error: %w", err)
		case html.StartTagToken:
			t := token.Token()
			isTitle = t.Data == "title"
		case html.TextToken:
			t := token.Token()
			if isTitle {
				return t.Data, nil
			}
		}
	}
}
