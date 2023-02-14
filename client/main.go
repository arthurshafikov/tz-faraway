package main

import (
	"log"
	"sync"
	"time"

	"github.com/arthurshafikov/tz-faraway/client/internal/app"
)

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 250; i++ {
		wg.Add(1)
		go func(i int) {
			start := time.Now()

			app.Run()

			elapsed := time.Since(start)
			log.Printf("Query %v took %s", i, elapsed)
			wg.Done()
		}(i)
	}

	wg.Wait()
}
