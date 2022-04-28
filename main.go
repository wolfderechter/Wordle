package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	guesses := make([]string, 0, 4)

	//Get list of words
	request, err := http.NewRequest("GET", "https://gist.githubusercontent.com/cfreshman/a03ef2cba789d8cf00c08f767e0fad7b/raw/28804271b5a226628d36ee831b0e36adef9cf449/wordle-answers-alphabetical.txt", nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal("Something went wrong when making the request to fetch the words:", err)
	}
	respBytes, _ := io.ReadAll(resp.Body)
	wordsString := string(respBytes)

	wordsList := strings.Split(wordsString, "\n")
	randomNumber := rand.New(rand.NewSource(time.Now().UnixNano()))
	correctWord := wordsList[randomNumber.Intn(len(wordsList))]
	//uncomment to see what the word will be
	//fmt.Println(correctWord)

	//uncomment to use a predefined word
	//correctWord = "smash"

	var input string
	for i := 0; i < 5; i++ {
		for {
			fmt.Println("Enter your guess:")
			//will read input until \n (enter) character
			input, _ = reader.ReadString('\n')
			//need to remove whitespace from input
			input = strings.TrimSpace(input)
			input = strings.ToLower(input)
			if len(input) != 5 {
				fmt.Printf("Only 5 characters allowed!\n")
			} else if !contains(wordsList, input) {
				fmt.Printf("Only valid words allowed!\n")
			} else {
				break
			}
		}
		fmt.Printf("_______________\n")

		guesses = append(guesses, input)
		for _, g := range guesses {
			printResult(correctWord, g)
		}
	}
	fmt.Printf("You didn't get it :(\n")
	fmt.Printf("It was: %v\n", correctWord)
}

func printResult(correctWord string, guess string) {
	correctChars := 0
	for index, charRune := range guess {
		char := fmt.Sprintf("%c", charRune)

		//char bestaat in het woord EN zit op de juiste plaats
		// if strings.Contains(correctWord, char) && (strings.Index(correctWord, char) == index) {
		if strings.Contains(correctWord, char) && strings.Split(correctWord, "")[index] == char {
			green := color.New(color.FgGreen).SprintFunc()
			fmt.Printf("%s", green(char))
			correctChars++
			//char bestaat in het woord maar zit op de verkeerde plaats
			//het aantal keer dat een char in de correctword zit moet groter zijn dan aantal keer in mijn guess anders is alle chars goed geraden
		} else if strings.Contains(correctWord, char) && strings.Count(correctWord, char) >= strings.Count(guess[:index+1], char) {
			yellow := color.New(color.FgYellow).SprintFunc()
			fmt.Printf("%s", yellow(char))

			//char is fout
		} else {
			fmt.Printf("%s", char)
		}
	}
	fmt.Printf("\n")
	if correctChars == 5 {
		fmt.Println("CONGRATS :)")
		os.Exit(1)
	}
}

func contains(wordsList []string, item string) bool {
	for _, s := range wordsList {
		if s == item {
			return true
		}
	}
	return false
}
