package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const wotdURL = "https://www.merriam-webster.com/word-of-the-day/"

// Use cases for wotd
// 1. Get the word of the day
// 2. Generating a random (previous) word of the day
// 3. Generating a random word of the day by providing a date (optional - not implemented)

// Events that must take place to present a word:
// 1. The user picks between 3 commands - wotd `get`,  wotd `random``
// 2. The program then generates the current date.
// 3. If the current date or the random date are present in the json file => return word info.
// 4. Else, fetch the word given the date and its information.
// 5. Save the word and its information to wotd.json.

func main() {

	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	rndCmd := flag.NewFlagSet("random", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected `get` subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "get":
		HandleGet(getCmd)
	case "random":
		HandleRnd(rndCmd)
	default:
	}

}

func HandleGet(getCmd *flag.FlagSet) {

	getCmd.Parse(os.Args[2:])
	today := time.Now().UTC().Format("2006-01-02")
	checkDate(today, wotdURL)
}

func HandleRnd(rndCmd *flag.FlagSet) {

	rndCmd.Parse(os.Args[2:])
	date := getRandomDate()
	rndURL := getDateURL(date)
	checkDate(date, rndURL)
}

// Creates a random date between 2006-10-10 to the current date
// When you pass a negative integer to AddDate(), you get the date for (x) days ago.
// FYI: This technique could be used to fetch a date for (x) months or (x) years ago.
func getRandomDate() string {

	start := time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC)
	today := time.Now().UTC()
	diff := today.Sub(start)
	d := int(diff.Hours() / 24)
	rand.Seed(time.Now().UnixNano())
	randomDay := rand.Intn(d - 0)
	randomDate := today.AddDate(0, 0, -randomDay).Format("2006-01-02")

	return randomDate
}

// Use the date of a word to reference the json file.
// This eliminates the possibility of having more than 1 word
// per day, but it does not eliminate reused words.
// FYI: Merriam-Webster tends to reuse words.
func checkDate(d, url string) string {

	words := getWord()

	for _, word := range words {
		if word.Date == d {
			fmt.Println(word.Word)
			return word.Word // return the word - or the word struct with all the word information
		}
	}

	w := formatter(getWordTitle(url))

	fmt.Println(w)
	saveInfo(d, w)
	return w
}

func formatter(title string) string {

	res := strings.ReplaceAll(title, "Word of the Day: ", "")
	res = strings.ReplaceAll(res, " | Merriam-Webster", "")

	return res
}

// Params:
// date, word, class, meaning, defintion, and example
func saveInfo(d, w string) {

	wInfo := word{
		Date: d,
		Word: w,
	}

	// read and write word information to existing information in wotd.json
	info := getWord()
	info = append(info, wInfo)
	saveWord(info)
}

// Concatenates the wotd URL with a random date
func getDateURL(date string) string {

	var url strings.Builder
	url.WriteString(wotdURL)
	url.WriteString(date)

	return url.String()
}

// Credit for this function: https://zetcode.com/golang/net-html/
// Jan Bodnar - https://twitter.com/janbodnar
func getWordTitle(url string) string {

	res, err := http.Get(url)

	if err != nil {
		fmt.Println("Failed to get word title")
	}

	defer res.Body.Close()

	token := html.NewTokenizer(res.Body)

	var isTitle bool

	for {
		tt := token.Next()
		switch {
		case tt == html.ErrorToken:
			return "error"
		case tt == html.StartTagToken:

			t := token.Token()
			isTitle = t.Data == "title"

		case tt == html.TextToken:

			t := token.Token()

			if isTitle {
				isTitle = false
				return t.Data
			}
		}
	}
}
