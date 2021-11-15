package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/charmap"
)

// den it is Czech days mapping
var den = []string{
	"Neděle",
	"Pondělí",
	"Úterý",
	"středa",
	"Čtvrtek",
	"Pátek",
	"Sobota",
}
var wg sync.WaitGroup

// Scraping the webPage with goquery query if it is need decode from windows-1250 to utf-8 set decode to true.
// channel is for return value to gorutines.
func Scraping(webPage *string, query *string, decode bool, channel chan<- []string) {
	defer wg.Done()
	var result []string
	//HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	request, err := http.NewRequest("GET", *webPage, nil)
	if err != nil {
		fmt.Println(err)
	}
	response, err := client.Do(request)

	if response.StatusCode == 200 {
		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		// Page title
		if decode {
			result = append(result, DecodeWindows1250([]byte(doc.Find("head title").Text())))
		} else {
			result = append(result, doc.Find("head title").Text())
		}

		// Find the items
		doc.Find(*query).Each(func(i int, item *goquery.Selection) {
			text := strings.TrimSpace(item.Children().Remove().End().Text()) // get only root element without children element
			if decode {
				text = DecodeWindows1250([]byte(text))
			}
			result = append(result, text)
		})
	}
	channel <- result
}

func DecodeWindows1250(enc []byte) string {
	dec := charmap.Windows1250.NewDecoder()
	out, _ := dec.Bytes(enc)
	return string(out)
}

// trimResult get today menu based on weekday from text.
// return slice of string with today menu and page title
func trimResult(text []string, weekday *string) []string {
	today := regexp.MustCompile("(?i)" + *weekday + "(\\d.\\d.\\d)?")
	day := regexp.MustCompile("(?i)(" + strings.Join(den, ")|(") + ")(\\d.\\d.\\d)?")
	start, end := 0, len(text)
	var result []string

	for i, t := range text {
		if today.MatchString(t) {
			start = i + 1
		} else if day.MatchString(t) {
			if start != 0 {
				end = i
				break
			}
		}
	}
	if start == 0 || end <= start {
		result = text[0:1] // page title
	} else {
		result = append(text[0:1], text[start:end]...)
	}

	return result
}

type VSlice []string

//VSlice String convert slice of string to a string with new line after each element
func (s VSlice) String() string {
	var str string
	for _, i := range s {
		str += fmt.Sprintf("%v\n", i)
	}
	return str
}

func main() {
	weekday := den[time.Now().Weekday()]
	menu := [...]string{
		"https://www.pivnice-ucapa.cz/denni-menu.php",
		"https://www.suzies.cz/poledni-menu",
		"https://www.menicka.cz/4921-veroni-coffee--chocolate.html",
	}
	query := [...]string{
		"div.listek div.day,div.listek div.polevka,div.listek div.food",
		"div.uk-card-body h2,div.uk-card-body div.uk-width-expand,div.uk-card-body h3",
		"div.menicka div.nadpis,div.menicka div.polozka",
	}
	decode := [...]bool{
		false,
		false,
		true}
	var result [][]string
	channel := make(chan []string, 20)

	for i := 0; i < len(menu); i++ {
		wg.Add(1)
		go Scraping(&menu[i], &query[i], decode[i], channel)
		tmp := <-channel
		result = append(result, tmp)
	}
	wg.Wait()

	fmt.Printf("Dneska je %v a menucka su taketo:\n\n", weekday)
	for _, r := range result {
		r = trimResult(r, &weekday)
		fmt.Println(VSlice(r))
	}
	fmt.Println("Dobru chut :)")
}
