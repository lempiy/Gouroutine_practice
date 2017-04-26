package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"sync"
)

const maxGoroutines = 100

func main() {
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("Usage: $ %s <img_url>, ...\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	imgUrls := os.Args[1:]
	infoChan := make(chan string, maxGoroutines*2)

	go findImages(infoChan, imgUrls)
	htmlTags := mergeResults(infoChan)
	for _, tag := range htmlTags {
		fmt.Println(tag)
	}
}

func findImages(infoChan chan string, imgUrls []string) {
	waiter := &sync.WaitGroup{}
	for _, url := range imgUrls {
		waiter.Add(1)
		go proccessImage(infoChan, url, waiter)
	}
	waiter.Wait()
	close(infoChan)
}

func proccessImage(infoChan chan string, imgURL string, waiter *sync.WaitGroup) {
	file, err := os.Open(imgURL)
	defer file.Close()
	if err != nil {
		waiter.Done()
		return //ignore errors
	}
	config, _, err := image.DecodeConfig(file)
	if err != nil {
		waiter.Done()
		return //ignore errors
	}
	infoChan <- fmt.Sprintf(
		`<img src="%s" width="%dpx" height="%dpx" />`,
		filepath.Base(imgURL), config.Width, config.Height)
	waiter.Done()
}

func mergeResults(infoChan <-chan string) []string {
	var results []string
	for info := range infoChan {
		results = append(results, info)
	}
	return results
}
