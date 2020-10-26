package search

import (
	"context"
	"testing"
)

func TestSearch_All_success(t *testing.T) {
	ch := All(context.Background(), "Ilhom", []string{"txt1.txt", "txt.txt"})
	res, ok := <-ch
	if !ok {
		t.Errorf("error in All success")
		return
	}
	for _, i := range res {
		t.Log(i)
	}
}

func BenchmarkSearch_Any_success(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ch := Any(context.Background(),"Ilhom",[]string{"txt1.txt","txt.txt"})
		val, ok := <-ch
		b.StopTimer()
		if !ok{
			b.Errorf("error %v",ok)
			return
		}
		b.Log(val)
		b.StartTimer()
	}
}