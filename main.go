package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/briandowns/spinner"
)

var artist string
var song string
var Green = "\033[32m"
var Red = "\033[31m"

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	//take input from user
	fmt.Print("Artist: ")
	scanner.Scan()
	artist = scanner.Text()
	fmt.Print("Song: ")
	scanner.Scan()
	song = scanner.Text()

	//Remove spaces and convert to lower case
	lowerCaseSongNoSpaces := strings.ToLower(strings.Replace(song, " ", "", -1))
	lowerCaseArtistNoSpaces := strings.ToLower(strings.Replace(artist, " ", "", -1))

	url := "https://azlyrics.com/lyrics/" + lowerCaseArtistNoSpaces + "/" + lowerCaseSongNoSpaces + ".html"

	//spinner
	s := spinner.New(spinner.CharSets[36], 100*time.Millisecond)
	s.Start()
	time.Sleep(1 * time.Second)

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
		x := getX()
		url = fmt.Sprintf("https://search.azlyrics.com/?q=%s+%s&w=songs&x=%s", artist, song, x)
		s.Stop()
		searchSong(url)
	} else {
		clearTerminal()
		s.Stop()
		getSong(url)
	}
}

// gets the value of the 'x' parameter from the geo.js file
func getX() string {
	url := "https://www.azlyrics.com/geo.js"

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

	bodyContent, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(`"value", "([^"]+)"`)

	matches := re.FindStringSubmatch(string(bodyContent))
	if len(matches) < 2 {
		log.Fatal("Value not found")
	}

	return matches[1]
}

func searchSong(url string) {

	// Get the search results
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the alert message for no results is present
	noResults := doc.Find(".alert-warning").Text()
	if strings.Contains(noResults, "Sorry, your search returned no results") {
		fmt.Println(Red + "Sorry no results found :(")
		return
	}

	clearTerminal()
	fmt.Println("Search Results:")
	fmt.Println()

	// Find all song links
	doc.Find(".visitedlyr a").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("span").First().Text())
		artist := strings.TrimSpace(s.Find("b").Last().Text())

		fmt.Printf("%d. \"%s\" by %s\n", i+1, title, artist)
	})

	fmt.Println()

	// Prompt the user to enter a number to get the URL of the corresponding song
	var choice int
	fmt.Print("Enter the number of the song: ")
	_, err = fmt.Scan(&choice)
	if err != nil {
		fmt.Println(Red + "Invalid input")
		return
	}

	// Find the URL corresponding to the user's choice
	selectedSongURL := doc.Find(".visitedlyr a").Eq(choice-1).AttrOr("href", "")
	getSong(selectedSongURL)

}

// Get the lyrics of the selected song
func getSong(url string) {
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

	// Clear the terminal
	clearTerminal()

	// Get the lyrics of the selected song
	fmt.Println("Lyrics:")
	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(doc.Text(), "\n")

	for i, line := range lines {
		for n := 175; n < 3000; n++ {
			if i == n && line != "" {
				fmt.Println(Green + line)
			}
			if strings.Contains(line, "Android|webOS|iPhone|iPod|iPad|BlackBerry") {
				os.Exit(0)
			}
		}
	}
}

func clearTerminal() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
