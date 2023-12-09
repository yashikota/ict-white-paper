package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/schollz/progressbar"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main() {
	files, err := os.ReadDir("./url")
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.Create("./data/all.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	for _, file := range files {
		// skip index.txt
		if file.Name() == "index.txt" {
			continue
		}

		filePath := "./url/" + file.Name()
		fmt.Println("\n" + filePath)

		f, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		// create file
		f2, err := os.Create("./data/" + path.Base(filePath))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f2.Close()

		// read lines
		scanner := bufio.NewScanner(f)
		var urls []string
		for scanner.Scan() {
			urls = append(urls, scanner.Text())
		}

		downloadFile(path.Base(filePath), urls)
	}

	fmt.Println("\nDone!")
}

func downloadFile(filePath string, urls []string) {
	bar := progressbar.New(len(urls))
	for _, url := range urls {
		bar.Add(1)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		// if status code is 404, skip
		if resp.StatusCode == 404 {
			fmt.Println("\n404: " + url)
			continue
		}

		// convert encoding response body
		reader := transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())

		// write file
		f, err := os.OpenFile("./data/"+filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		f2, err := os.OpenFile("./data/all.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f2.Close()

		// write file
		w := io.MultiWriter(f, f2)
		if _, err := io.Copy(w, reader); err != nil {
			fmt.Println(err)
			return
		}
	}
}
