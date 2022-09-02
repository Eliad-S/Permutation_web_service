package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
)

func getPage(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return nil
}

func worker(urlCh chan string, sizeCh chan string, id int) {
	for {
		url := <-urlCh
		err := getPage(url)
		if err == nil {
			sizeCh <- fmt.Sprintf("[%d] %s succedded", id, url)
		} else {
			sizeCh <- fmt.Sprintf("Error getting %s :: %s", url, err)
		}
	}
}

func generator(url string, urlCh chan string) {
	urlCh <- url
}

func Create_url_list(file_name string) ([]string, error) {
	urls := []string{}
	file, err := os.Open(file_name)
	if err != nil {
		log.Fatalf("Can't open file %s. Err: %s", file_name, err)
		return []string{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		url := fmt.Sprintf("http://localhost:8080//api/v1/similar?word=%s", scanner.Text())
		urls = append(urls, url)
		// fmt.Println(scanner.Text())
	}
	return urls, nil
}

func main() {

	urlCh := make(chan string)
	sizeCh := make(chan string)

	urls, err := Create_url_list("../words_clean.txt")
	if err != nil {
		log.Fatalf(err.Error())
	}

	for i := 0; i < 10; i++ {
		go worker(urlCh, sizeCh, i)
	}

	for _, url := range urls {
		go generator(url, urlCh)
	}

	for i := 0; i < len(urls); i++ {
		fmt.Printf("%s\n", <-sizeCh)
	}
	fmt.Printf("done [%d] urls\n", len(urls))
}
