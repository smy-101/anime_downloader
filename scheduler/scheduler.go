package scheduler

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type Task struct {
	Name     string
	Interval time.Duration
	Execute  func() error
}

type Scheduler struct {
	tasks []*Task
}

// NewScheduler 创建一个新的 Scheduler 实例
func NewScheduler() *Scheduler {
	return &Scheduler{
		tasks: make([]*Task, 0),
	}
}

// AddTask 添加一个新任务到调度器
func (s *Scheduler) AddTask(name string, interval time.Duration, execute func() error) {
	task := &Task{
		Name:     name,
		Interval: interval,
		Execute:  execute,
	}
	s.tasks = append(s.tasks, task)
}

// Start 开始运行所有任务
func (s *Scheduler) Start() {
	for _, task := range s.tasks {
		go s.runTask(task)
	}
}

// runTask 运行单个任务
func (s *Scheduler) runTask(task *Task) {
	ticker := time.NewTicker(task.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("执行任务: %s\n", task.Name)
			if err := task.Execute(); err != nil {
				fmt.Printf("任务 %s 执行出错: %v\n", task.Name, err)
			}
		}
	}
}

// RunAsDaemon 将调度器作为守护进程运行
func (s *Scheduler) RunAsDaemon() {
	if !isParent() {
		cmd := exec.Command(os.Args[0], "-daemon")
		cmd.Start()
		fmt.Printf("守护进程已启动,PID: %d\n", cmd.Process.Pid)
		os.Exit(0)
	}

	// logFile,err := os.OpenFile("scheduler.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	logFile, err := os.OpenFile("scheduler.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.Println("守护进程已启动")
	s.Start()

	//select{} 会一直阻塞当前 goroutine，使得程序不会退出
	select {}
}

func isParent() bool {
	return os.Getppid() != 1
}
