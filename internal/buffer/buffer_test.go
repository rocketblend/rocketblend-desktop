package buffer

import (
	"sync"
	"testing"
)

func TestAddDataDropsOldest(t *testing.T) {
	bm := New(WithMaxBufferSize(2))
	bm.AddData(1)
	bm.AddData(2)
	bm.AddData(3)
	bm.Close()

	var got []Data
	for data := range bm.Channel() {
		got = append(got, data)
	}

	if len(got) != 2 || got[0] != 2 || got[1] != 3 {
		t.Fatalf("got %v, want [2 3]", got)
	}
}

func TestAddDataDuringClose(t *testing.T) {
	bm := New(WithMaxBufferSize(4))

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				bm.AddData(j)
			}
		}()
	}

	go func() {
		for range bm.Channel() {
		}
	}()

	bm.Close()
	wg.Wait()
}
