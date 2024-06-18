package handler

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mosteligible/go-brrrr/pkg/loads"
	"github.com/mosteligible/go-brrrr/pkg/types"
)

type ConcurrencyHandler struct {
	Concurrency  int     // num of concurrent requests to execute
	rate         float64 // concurrency per rate seconds
	lock         *sync.Mutex
	mq           *types.MetricsQueue
	testDuration float64
	loadQueue    chan loads.Loader
	Count        int
	userTracker  chan bool
}

func NewConcurrencyHandler(concurrency int, rate float64, testDuration float64) *ConcurrencyHandler {
	mq := types.NewMetricsQueue(concurrency)
	ch := ConcurrencyHandler{
		Concurrency:  concurrency,
		lock:         new(sync.Mutex),
		rate:         rate,
		mq:           mq,
		Count:        0,
		testDuration: testDuration,
		loadQueue:    make(chan loads.Loader, concurrency),
		userTracker:  make(chan bool, concurrency),
	}

	return &ch
}

func (ch *ConcurrencyHandler) Add(item loads.Loader) {
	ch.Count++
	ch.loadQueue <- item
}

func (ch *ConcurrencyHandler) Consume() loads.Loader {
	return <-ch.loadQueue
}

func (ch *ConcurrencyHandler) AddUser() {
	ch.userTracker <- true
}

func (ch *ConcurrencyHandler) RemoveUser() {
	<-ch.userTracker
}

func (ch *ConcurrencyHandler) IoLoopNreqps(parameters *types.Parameters) {
	testDuration := time.Now()
	var waitGrp sync.WaitGroup
	go MetricsBuilder(ch.mq)

	for {
		batchStart := time.Now()
		go func() {
			for range ch.Concurrency {
				item := <-ch.loadQueue
				waitGrp.Add(1)
				go item.Load(
					parameters,
					ch.mq,
					&waitGrp,
					batchStart.Format(time.DateTime),
				)
			}
		}()
		if time.Since(testDuration).Seconds() >= ch.testDuration {
			log.Printf("\nStopping test, time taken: %f", time.Since(testDuration).Seconds())
			break
		}
		nextStart := batchStart.Add(time.Duration(ch.rate * float64(time.Second)))
		fmt.Printf("start: %s\n next: %s\n", batchStart, nextStart)
		if nextStart.After(time.Now()) {
			sleepTime := nextStart.Sub(time.Now())
			fmt.Printf("-- sleeping for %s - rate: %fs\n\n", sleepTime, ch.rate)
			time.Sleep(sleepTime)
		}
	}
	fmt.Println("-- waitgroup wait n users/s")
	waitGrp.Wait()
	ShowResults(ch.mq)
}

func (ch *ConcurrencyHandler) IoLoopNUsers(parameters *types.Parameters) {
	testDuration := time.Now()
	var waitGrp sync.WaitGroup
	go MetricsBuilder(ch.mq)
	for {
		item := <-ch.loadQueue
		waitGrp.Add(1)
		ch.AddUser()
		nowStr := time.Now().Format(time.DateTime)
		go func() {
			item.Load(
				parameters,
				ch.mq,
				&waitGrp,
				nowStr,
			)
			ch.RemoveUser()
		}()

		if time.Since(testDuration).Seconds() >= ch.testDuration {
			log.Printf("\nStopping test, time taken: %f", time.Since(testDuration).Seconds())
			break
		}
	}
	fmt.Println("-- waitgroup wait n users/s")
	waitGrp.Wait()
	ShowResults(ch.mq)
}
