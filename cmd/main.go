package main

import (
	"context"
	"log"

	"github.com/ilhom0258/search/pkg/search"
)

func main() {
	ch := search.Any(context.Background(), "Ilhom", []string{"txt.txt"})
	r, ok := <-ch
	log.Print(r, ok)
}
