package tailf

import (
	"context"
	"github.com/hpcloud/tail"
	"log"
	"time"

	kafkaclient "logagent/kafkaClient"
)

var config tail.Config

func init() {
	config = tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 //是否follow
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件哪里开始读
		MustExist: false,                                // 文件不存在报错
		Poll:      true,
	}
}

// Task 一个日志任务实例
type Task struct {
	Path     string
	Topic    string
	instance *tail.Tail
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewTailTask(Path, Topic string) (task *Task) {
	// one task match on tail.Tail
	task = &Task{
		Path:  Path,
		Topic: Topic,
	}
	task.InitTail() // 打开要追踪文件
	return
}

// InitTail one task match on tail.Tail
func (t *Task) InitTail() {
	log.Printf("Tail: Begining to watch %s ", t.Path)
	var err error
	t.instance, err = tail.TailFile(t.Path, config)
	if err != nil {
		log.Fatalf("watch %s failed, err: %v", t.Path, err)
		return
	}
	t.ctx, t.cancel = context.WithCancel(context.Background())
	go t.Run()
}

func (t *Task) Run() {
	t.TaskForCatch()
}

// TaskForCatch single instance to catch file
// ctx.Done should free tailf
func (t *Task) TaskForCatch() {
	// send message to kafka topic
	for {
		select {
		case line := <-t.instance.Lines:
			//  send to kafka
			kafkaclient.SentToChan(line.Text, t.Topic)

		case <-t.ctx.Done():
			log.Printf("Task [%s  %s] shutdown", t.Path, t.Topic)
			return
		default:
			time.Sleep(time.Second)
		}
	}
}
