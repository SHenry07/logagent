package tailf

import (
	"fmt"
	"log"
	"logagent/etcd"
	"strings"
)

var tskMgr *tailMgr

type tailMgr struct {
	logEntries []*etcd.LogEntry
	tskMap     map[string]*Task
	ConfChan   chan []*etcd.LogEntry
}

func Init(logentries []*etcd.LogEntry) {
	tskMgr = &tailMgr{
		logEntries: logentries, // 将当前配置保存起来
		tskMap:     make(map[string]*Task, 16),
		ConfChan:   make(chan []*etcd.LogEntry), // 无缓冲的watch conf通道
	}
	for _, logentry := range logentries {
		// 记录 为了方便后续判断
		task := NewTailTask(logentry.Path, logentry.Topic)
		//info := logentry.Path + logentry.Topic
		info := fmt.Sprintf("%s^%s", logentry.Path, logentry.Topic)
		tskMgr.tskMap[info] = task
	}
	go Compare()
}

func Compare() {
	for {
		newlogentries := <-tskMgr.ConfChan
		var newtskMap map[string]int
		newtskMap = make(map[string]int)
		for _, newlogentry := range newlogentries {
			newinfo := fmt.Sprintf("%s^%s", newlogentry.Path, newlogentry.Topic)
			newtskMap[newinfo] = 1
		}
		//log.Printf("new configure coming, %#v",newtskMap)
		//log.Printf("new configure coming, %#v",tskMgr.tskMap)

		// 遍历旧的task map 是否出现在新的task map中
		for oldinfo, oldTask := range tskMgr.tskMap {
			if _, ok := newtskMap[oldinfo]; ok {
				delete(newtskMap, oldinfo)
				continue
			} else {
				// showdown old Task
				delete(tskMgr.tskMap, oldinfo)
				oldTask.cancel()
			}
		}

		// new Task
		for newinfo := range newtskMap {
			tmp := strings.Split(newinfo, "^")
			path := tmp[0]
			topic := tmp[1]
			log.Printf("new configure coming, path:%s topic:%s", path, topic)
			task := NewTailTask(path, topic)
			tskMgr.tskMap[newinfo] = task
		}
	}
}

//tskMgr.tskMap[oldlogentry.Path].cancel()
//delete(tskMgr.tskMap, oldlogentry.Path)
//// 重新初始化task
//
//// updata
//tskMgr.tskMap[oldlogentry.Path].Topic = newlogentry.Topic

func NewConfChan() chan []*etcd.LogEntry {
	return tskMgr.ConfChan
}
