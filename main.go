package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/briandowns/spinner"
)

func main() {
	var artist string
	var song string
	var Green = "\033[32m"
	var Red = "\033[31m"

	scanner := bufio.NewScanner(os.Stdin)

	//take input from user
	fmt.Print("Artist: ")
	scanner.Scan()
	artist = scanner.Text()
	fmt.Print("Song: ")
	scanner.Scan()
	song = scanner.Text()

	//Remove spaces
	song = strings.Replace(song, " ", "", -1)
	artist = strings.Replace(artist, " ", "", -1)

	//lowercase
	artist = strings.ToLower(artist)
	song = strings.ToLower(song)

	url := "https://azlyrics.com/lyrics/" + artist + "/" + song + ".html"

	//spinner
	s := spinner.New(spinner.CharSets[36], 100*time.Millisecond)
	s.Start()
	time.Sleep(1 * time.Second)
	s.Stop()

	//clear screen
	fmt.Print("\033[H\033[2J")
	fmt.Print("\033[H\033[2J")
	fmt.Print("\033[H\033[2J")

	//get lyrics
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
		fmt.Println("Lyrics:")
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
