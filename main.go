package main

import (
	"strconv"
	"time"
	"fmt"
)

func main() {
	w := NewHWorker()
	w.Register("MyFunc", MyFunc)

	for i := 0; i <= 100; i++ {
		w.Add(QueueList{
			QueueName: "MyFunc" + strconv.Itoa(i),
			FuncName:  "MyFunc",
			Args:      []interface{}{"a", "b"},
		})
	}

	go func() {
		time.Sleep(time.Second*10)

		for i := 0; i <= 100; i++ {
			w.Add(QueueList{
				QueueName: "MyFunc" + strconv.Itoa(i),
				FuncName:  "MyFunc",
				Args:      []interface{}{"a", "b"},
			})
		}
	}()

	w.Worker()
}

func MyFunc(queueName string, t ...interface{}) error {
	fmt.Println(queueName, " ----------- ", t)
	return nil
}
