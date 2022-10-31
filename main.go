package main

import (
	"fmt"
	"sync"
	"time"
)

type hoge struct {
	c           *sync.Cond
	sharedValue map[string]int
}

var vals = []string{"hoge", "fuga", "piyo"}

func (h *hoge) set(str string, i int, wg *sync.WaitGroup) {
	h.c.L.Lock()
	defer h.c.L.Unlock()
	h.sharedValue[str] = i
	fmt.Printf("waiting\n")
	h.c.Wait()
	fmt.Printf("%s end!\n", str)
	wg.Done()
}

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU())
	rwmux := new(sync.RWMutex)
	c := sync.NewCond(rwmux)
	h := &hoge{
		c:           c,
		sharedValue: make(map[string]int),
	}
	wg := new(sync.WaitGroup)

	wg.Add(len(vals))
	for i, val := range vals {
		go h.set(val, i, wg)
	}

	time.Sleep(3 * time.Second)

	fmt.Println(h.sharedValue)
	h.c.Broadcast()
	wg.Wait()
}
