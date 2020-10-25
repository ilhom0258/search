package main

import (
	"context"
	"log"

	"github.com/ilhom0258/search/pkg/search"
)

func main() {
	ch := search.All(context.Background(), "Ilhom", []string{"text.txt"})
	r, ok := <-ch
	log.Print(r, ok)
}
