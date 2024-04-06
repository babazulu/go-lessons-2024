package main

import (
	"balun-school-task-2/scheduler"
	"fmt"
	"time"
)

type OneJob struct {
}

func (j *OneJob) Do() {
	fmt.Println("DO")
}

func (j *OneJob) Cancel() {
	fmt.Println("STOP")
}

func main() {
	sch := scheduler.NewScheduler()

	job1 := &OneJob{}
	job2 := &OneJob{}

	sch.SetTimeout(job1, 3*time.Second)

	id2 := sch.SetTimeout(job2, 10*time.Second)

	time.Sleep(2 * time.Second)
	sch.CancelTimeout(id2)

	time.Sleep(3 * time.Second)

}
