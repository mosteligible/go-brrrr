package types

import (
	"fmt"
	"sync"
)

type MetricsQueue struct {
	limit        int
	failures     int
	success      int
	Tracker      chan bool
	lock         *sync.Mutex
	ReqSent      int
	ReqCompleted int
}

func NewMetricsQueue(limit int) *MetricsQueue {
	mq := MetricsQueue{
		limit:    limit,
		failures: 0,
		success:  0,
		Tracker:  make(chan bool, limit),
		lock:     new(sync.Mutex),
	}

	return &mq
}

func (mq *MetricsQueue) String() string {
	return fmt.Sprintf(
		"Limit: %d\nFailures: %d\nSuccess: %d\nCount: %d",
		mq.limit, mq.failures, mq.success, len(mq.Tracker),
	)
}

func (mq *MetricsQueue) Add() {
	mq.Tracker <- true
}

func (mq *MetricsQueue) AddSuccess() {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	mq.success++
	<-mq.Tracker
}

func (mq *MetricsQueue) AddFailure() {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	mq.failures++
	<-mq.Tracker
}

func (mq *MetricsQueue) NewReqSent() {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	mq.ReqSent++
}

func (mq *MetricsQueue) NewReqCompleted() {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	mq.ReqCompleted++
}

func (mq *MetricsQueue) GetCount() int {
	return len(mq.Tracker)
}
