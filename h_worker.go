package main

import (
	"sync"
	"errors"
	"time"
	"log"
	"fmt"
	)

// worker func
type workerFunc func(queueName string, args ...interface{}) error

type Workers struct {
	workers map[string]workerFunc
}

type HWorker struct {
	Workers
	ConsumeQueue []QueueList
	mt           sync.Mutex
}

func NewHWorker() *HWorker {
	o := &HWorker{}
	o.workers = make(map[string]workerFunc, 0)
	return o
}

// 注册可执行方法
func (this *HWorker) Register(funcName string, f workerFunc) error {
	this.mt.Lock()
	defer this.mt.Unlock()
	// 判断是否已经注册
	if _, ok := this.workers[funcName]; ok {
		return errors.New("已经注册该方法")
	}
	this.workers[funcName] = f
	return nil
}

// 添加任务队列
func (this *HWorker) Add(queueList QueueList) error {
	this.ConsumeQueue = append(this.ConsumeQueue, queueList)
	return nil
}

// 读取任务,并调用注册函数
func (this *HWorker) Worker() error {
	wg := &sync.WaitGroup{}
	for {
		// 监听队列
		if (len(this.ConsumeQueue) == 0) {
			time.Sleep(time.Nanosecond * 500)
			continue
		}

		runQueue := this.ConsumeQueue[0]
		this.ConsumeQueue = this.ConsumeQueue[1:len(this.ConsumeQueue)]

		if _, ok := this.workers[runQueue.FuncName]; !ok {
			log.Println("队列方法未定义")
			continue
		}
		// 执行
		this.Run(this.workers[runQueue.FuncName], runQueue, wg)
	}
	wg.Wait()
	return nil
}

func (this *HWorker) Run(f workerFunc, queueTask QueueList, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("HWorker Run Panicking %s\n", fmt.Sprint(r))
			}
		}()
		err := f(queueTask.QueueName, queueTask.Args)
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}()
}
