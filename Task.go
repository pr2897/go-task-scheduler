package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TaskStatus uint8

const (
	Pending TaskStatus = iota
	Running
	Failed
	Completed
)

type Task struct {
	Id            uuid.UUID
	Name          string
	Status        TaskStatus
	ScheduledFor  time.Time
	ExecutionFunc func() error
	CompletedAt   time.Time
}

type TaskScheduler struct {
	mu        sync.Mutex
	tasks     map[uuid.UUID]*Task
	taskCount int
}

func NewTaskScheduler() *TaskScheduler {
	return &TaskScheduler{
		tasks: make(map[uuid.UUID]*Task),
	}
}

func (ts *TaskScheduler) AddTask(name string, scheduledFor time.Time, executionFunc func() error) uuid.UUID {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.taskCount++
	task := &Task{
		Id:            uuid.New(),
		Name:          name,
		Status:        Pending,
		ScheduledFor:  scheduledFor,
		ExecutionFunc: executionFunc,
	}
	ts.tasks[task.Id] = task
	fmt.Printf("Task '%v' scheduled for '%v'\n", name, scheduledFor)
	return task.Id
}

func (ts *TaskScheduler) RemoveTask(id uuid.UUID) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	delete(ts.tasks, id)
	fmt.Printf("Task Id %v removed \n", id.String())
}

func (ts *TaskScheduler) RunTask(id uuid.UUID) {
	ts.mu.Lock()

	task, exists := ts.tasks[id]
	if !exists {
		fmt.Printf("Task Id %v doesn't exists", id.String())
		ts.mu.Unlock()
		return
	}

	task.Status = Running
	ts.mu.Unlock()

	// execute the task function
	err := task.ExecutionFunc()
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if err != nil {
		task.Status = Failed
		fmt.Printf("Task '%s' failed: %v\n", task.Name, err)
	} else {
		task.Status = Completed
		task.CompletedAt = time.Now()
		fmt.Printf("Task '%s' completed at %v\n", task.Name, task.CompletedAt)
	}
}

// schedule workers
func (ts *TaskScheduler) ScheduleWorker() {
	for {
		ts.mu.Lock()
		now := time.Now()

		for _, task := range ts.tasks {
			if task.Status == Pending && now.After(task.ScheduledFor) {
				go ts.RunTask(task.Id)
			}
		}

		ts.mu.Unlock()
		time.Sleep(time.Second * 1)
	}
}
