package app

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Task struct {
	Arrival           time.Time
	PriorityUpdatesAt time.Time
	priority          int
	name              string
	responseChan      chan interface{}
}

type scheduler struct {
	tasks     []*Task
	mutex     sync.Mutex
	condition *sync.Cond
}

func (a *AppStruct) worker(id int, scheduler *scheduler) {
	for {
		task := scheduler.getNextTask()
		fmt.Printf("Worker %d executing task: %s with priority %d\n", id, task.name, task.priority)
		// Simulate task execution time
		a.MakeHTTPRequest(task)
	}
}

func (s *scheduler) addTask(task *Task) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	task.Arrival = time.Now()
	s.tasks = append(s.tasks, task)
	s.condition.Signal()
}

func (s *scheduler) getNextTask() *Task {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for len(s.tasks) == 0 {
		s.condition.Wait()
	}
	s.ageTasks() // Apply aging
	var highestPriorityTask *Task
	for _, task := range s.tasks {
		if highestPriorityTask == nil || task.priority > highestPriorityTask.priority {
			highestPriorityTask = task
		}
	}
	s.tasks = removeTask(s.tasks, highestPriorityTask)
	return highestPriorityTask
}

func (s *scheduler) ageTasks() {
	for _, task := range s.tasks {
		if int(time.Since(task.PriorityUpdatesAt).Seconds()) > 5 {
			// elapsedTicks := int(time.Since(task.Arrival).Seconds()) / 5
			log.Println("\n increaseing ", task.priority, " to: ", (task.priority + 1), "\n")
			task.priority++
			task.PriorityUpdatesAt = time.Now()
		}
	}
}

func removeTask(tasks []*Task, target *Task) []*Task {
	var result []*Task
	for _, task := range tasks {
		if task != target {
			result = append(result, task)
		}
	}
	return result
}

func (a *AppStruct) MakeHTTPRequest(t *Task) {
	resp, err := http.Get("http://localhost:8001/")
	fmt.Println("Response Received")
	if err != nil {
		log.Println("Failed to execute the http request")
		t.responseChan <- nil
		return
	}
	// log.Println("repsonse body:", resp.Body)
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	t.responseChan <- string(resBody)
}
