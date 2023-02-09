package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	var artist string
	var song string
	scanner := bufio.NewScanner(os.Stdin)
	var Green = "\033[32m"
	var Red = "\033[31m"

	fmt.Println(`

	 _            _       _
	| |_   _ _ __(_)_ __ | |_
	| | | | | '__| | '_ \| __|
	| | |_| | |  | | |_) | |_
	|_|\__, |_|  |_| .__/ \__|
	   |___/       |_|

	`)
	//input
	fmt.Println("Enter the artist name: ")
	scanner.Scan()
	artist = scanner.Text()

	fmt.Println("Enter the song name: ")
	scanner.Scan()
	song = scanner.Text()

	//Remove spaces
	song = strings.Replace(song, " ", "", -1)
	artist = strings.Replace(artist, " ", "", -1)

	//lowercase
	artist = strings.ToLower(artist)
	song = strings.ToLower(song)

	url := "https://azlyrics.com/lyrics/" + artist + "/" + song + ".html"

	fmt.Println()
	fmt.Println("Please wait...")

	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)

		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Println(Red + "Sorry, the song was not found. Please try again :(")
	} else {

		doc, err := goquery.NewDocumentFromReader(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		lines := strings.Split(doc.Text(), "\n")
		for i, line := range lines {
			for n := 174; n < 3000; n++ {
				if i == n && line != "" {
					fmt.Println(Green + line)
				}
				if strings.Contains(line, "Android|webOS|iPhone|iPod|iPad|BlackBerry") {
					os.Exit(0)
				}
			}
		}
	}
}
