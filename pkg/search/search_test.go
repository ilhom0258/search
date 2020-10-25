package search

import (
	"context"
	"log"
	"testing"
)

func TestSearch_All_success(t *testing.T) {
	ch := All(context.Background(), "Ilhom", []string{"txt.txt"})
	res, ok := <-ch
	if !ok {
		t.Errorf("error in All success")
		return
	}
	log.Printf("result = %v", res)
}

func TestSearch_Any_succes(t *testing.T) {
	ch := Any(context.Background(),"Ilhom",[]string{"txt.txt","txt1.txt"})
	res, ok := <-ch
	if !ok{
		t.Errorf("error in Any ok = %v", ok)
		return
	}	
	t.Logf("result = %v", res)
}
