package search

import (
	"context"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
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
	ch := make(chan []Result)
	ctx, cancel := context.WithCancel(ctx)
	wg := sync.WaitGroup{}

	for _, file := range files {
		wg.Add(1)
		log.Print(file)
		go func(ctx context.Context, ch chan<- []Result, fileName string, wg *sync.WaitGroup) {
			defer wg.Done()
			results, err := findAll(phrase, fileName)
			if err != nil {
				log.Print(err)
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
	ctx, cancel := context.WithCancel(ctx)
	wg := sync.WaitGroup{}

	for i := 0; i < len(files); i++ {
		select {
		case <-ctx.Done():
			cancel()
			break
		case <-time.After(time.Second):
			wg.Add(1)
			go findAnyConcurrent(ctx, ch, files[i], phrase, &wg, cancel)
		}
	}
	go func() {
		defer close(ch)
		wg.Wait()
	}()
	cancel()
	return ch
}

func findAnyConcurrent(ctx context.Context, ch chan<- Result, file string, phrase string, wg *sync.WaitGroup, cancel func()) {
	defer wg.Done()
	select{
	case <-ctx.Done():
		return
	default:
		result, err := findAny(phrase, file)
		if err != nil {
			log.Printf("%v error in go", err)
			return
		}
		if (result == Result{}) {
			return
		}
		<-ctx.Done()
		ch <- result
		cancel()
	}
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
