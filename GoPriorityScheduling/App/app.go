package app

import (
	"sync"
	"time"
)

type AppStruct struct {
	schedulerManager *scheduler
}

func NewApp() *AppStruct {
	s := &scheduler{}
	s.condition = sync.NewCond(&s.mutex)
	app := &AppStruct{
		schedulerManager: s,
	}
	go app.worker(1, s)
	return app
}

func (a *AppStruct) HandleRequest(priority int) interface{} {
	var taskName string
	if priority >= 1 {
		taskName = "high"
	} else {
		taskName = "low"
	}
	t := &Task{
		priority:          priority,
		Arrival:           time.Now(),
		PriorityUpdatesAt: time.Now(),
		name:              taskName,
		responseChan:      make(chan interface{}),
	}

	a.schedulerManager.addTask(t)

	response := <-t.responseChan

	return response

}
