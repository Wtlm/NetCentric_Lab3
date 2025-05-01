package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var genreURLs = map[string]string{
	"Drama":     "https://www.webtoons.com/en/drama",
	"Fantasy":   "https://www.webtoons.com/en/fantasy",
	"Comedy":    "https://www.webtoons.com/en/comedy",
	"Action":    "https://www.webtoons.com/en/action",
	"Romance":   "https://www.webtoons.com/en/romance",
	"Superhero": "https://www.webtoons.com/en/genres/super_hero",
	"Sci-fi":    "https://www.webtoons.com/en/genres/sf",
	"Horror":    "https://www.webtoons.com/en/horror",
	"Thriller":  "https://www.webtoons.com/en/thriller",
	"Sports":    "https://www.webtoons.com/en/genres/sports",
}

func main() {
	result := make(map[string][]string)

	for genre, url := range genreURLs {
		fmt.Println("Fetching:", genre, "from", url)
		titles := fetchTitles(url)
		if len(titles) > 10 {
			titles = titles[:10]
		}
		result[genre] = titles
	}

	// Write to JSON
	file, err := os.Create("webtoons.json")
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		fmt.Println("Error writing JSON:", err)
	}
	fmt.Println("Data saved")
}

func fetchTitles(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Request error:", err)
		return nil
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("HTML parse error:", err)
		return nil
	}

	var titles []string
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "p" {
			for _, a := range n.Attr {
				if a.Key == "class" && strings.Contains(a.Val, "subj") {
					if n.FirstChild != nil {
						titles = append(titles, strings.TrimSpace(n.FirstChild.Data))
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return titles
}
