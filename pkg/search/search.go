package search

import (
	"context"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

//Result returns
type Result struct {
	Phrase  string
	Line    string
	LineNum int64
	ColNum  int64
}

//All finds all occurence of text
func All(ctx context.Context, phrase string, files []string) <-chan []Result {
	ch := make(chan []Result, 5)
	ctx, cancel := context.WithCancel(ctx)
	wg := sync.WaitGroup{}

	for _, file := range files {
		wg.Add(1)
		go func(ctx context.Context, ch chan<- []Result, file string, wg *sync.WaitGroup) {
			defer wg.Done()
			results, err := findAll(phrase, file)
			if err != nil {
				return
			}
			if len(results) > 0 {
				ch <- results
			}
		}(ctx, ch, file, &wg)
	}
	go func() {
		defer close(ch)
		wg.Wait()
	}()
	cancel()
	return ch
}

// Any finds first one occurence of text
func Any(ctx context.Context, phrase string, files []string) <-chan Result {
	ch := make(chan Result)
	wg := sync.WaitGroup{}
	for _, file := range files {
		wg.Add(1)
		go func(ctx context.Context, ch chan<- Result, file string, wg *sync.WaitGroup) {
			defer wg.Done()
			result, err := findAny(phrase, file)
			if err != nil {
				return
			}
			if (result == Result{}) {
				return
			}
			ch <- result
		}(ctx, ch, file, &wg)
	}
	go func() {
		defer close(ch)
		wg.Wait()
	}()
	return ch
}

//Helper methods
func findAll(phrase string, file string) ([]Result, error) {
	dataRaw, err := readFile(file)
	if err != nil {
		log.Printf("%v error in findAll", err)
		return nil, err
	}
	results := []Result{}
	data := strings.Split(dataRaw, "\n")
	for i, item := range data {
		if strings.Contains(item, phrase) {
			line := int64(i) + 1
			col := int64(strings.Index(item, phrase)) + 1
			result := Result{
				Phrase:  phrase,
				Line:    item,
				LineNum: line,
				ColNum:  col,
			}
			results = append(results, result)
		}
	}
	return results, nil
}

func findAny(phrase string, file string) (result Result, err error) {
	dataRaw, err := readFile(file)
	if err != nil {
		log.Printf("%v error in findAny", err)
		return result, err
	}
	data := strings.Split(dataRaw, "\n")
	for i, item := range data {
		if strings.Contains(item, phrase) {
			line := int64(i) + 1
			col := int64(strings.Index(item, phrase)) + 1
			result = Result{
				Line:    item,
				Phrase:  phrase,
				LineNum: line,
				ColNum:  col,
			}
			break
		}
	}
	return result, err
}

func readFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Printf("error in reading file = %v", err)
		return "", err
	}
	return string(data), nil
}
