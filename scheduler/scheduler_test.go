package scheduler

import (
	"testing"
	"time"
)

func TestNewScheduler(t *testing.T) {
	s := NewScheduler()
	if s == nil {
		t.Fatal("NewScheduler 返回 nil")
	}
	if len(s.tasks) != 0 {
		t.Fatalf("预期 tasks 长度为 0,实际为 %d", len(s.tasks))
	}
}

// 测试 AddTask 函数
func TestAddTask(t *testing.T) {
	s := NewScheduler()
	taskName := "testTask"
	taskInterval := time.Second
	taskExecute := func() error {
		return nil
	}

	s.AddTask(taskName, taskInterval, taskExecute)

	if len(s.tasks) != 1 {
		t.Fatalf("预期 tasks 长度为 1,实际为 %d", len(s.tasks))
	}

	task := s.tasks[0]
	if task.Name != taskName {
		t.Fatalf("预期任务名称为 %s,实际为 %s", taskName, task.Name)
	}
	if task.Interval != taskInterval {
		t.Fatalf("预期任务间隔为 %v,实际为 %v", taskInterval, task.Interval)
	}
	if task.Execute == nil {
		t.Fatal("任务执行函数为 nil")
	}
}

// 测试 Start 函数
func TestStart(t *testing.T) {
	s := NewScheduler()
	taskExecuted := false
	taskName := "testTask"
	taskInterval := time.Millisecond * 100
	taskExecute := func() error {
		taskExecuted = true
		return nil
	}

	s.AddTask(taskName, taskInterval, taskExecute)
	s.Start()

	time.Sleep(taskInterval * 2)

	if !taskExecuted {
		t.Fatal("任务未执行")
	}
}

// 测试 RunAsDaemon 函数
func TestRunAsDaemon(t *testing.T) {
	// 由于 RunAsDaemon 涉及到进程管理和日志文件操作，测试起来比较复杂
	// 这里我们只测试 isParent 函数
	if !isParent() {
		t.Fatal("预期为父进程")
	}
}
