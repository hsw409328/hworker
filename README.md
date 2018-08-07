# hworker Go简单的worker使用
1、本地使用
2、目前不支持redis

# 使用
```
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

```

# 版本
v0.1

# 计划
下一步：增加redis支持