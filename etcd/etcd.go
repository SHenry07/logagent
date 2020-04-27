package etcd

import (
	"context"
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

var Cli *clientv3.Client

type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

func Init(endpoints []string) (err error) {
	Cli, err = clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	return nil
}

// PutConf send the configure of kafaka  to etcd
//func PutConf(key, val string)( err error) {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
//	defer cancel()
//
//	_, err = Cli.Put(ctx, key, val)
//	if err != nil {
//		log.Printf("etcd: put key to etcd failed: %v", err)
//		return err
//	}
//	return
//}

// GetConf  get conf from etcd
func GetConf(key string) (LogEntry []*LogEntry) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3000)
	defer cancel()

	resp, err := Cli.Get(ctx, key)
	if err != nil {
		log.Panicf("etcd:[getconf]: Cant reslove configure from etcd. err:%v", err)
	}

	for _, event := range resp.Kvs {
		err = json.Unmarshal(event.Value, &LogEntry)
		if err != nil {
			log.Panicf("etcd:[getconf]: Cant unmarshal configure. err:%v", err)
		}
	}
	return LogEntry
}

// WatchConf stand on the etcd
func WatchConf(key string, conf chan<- []*LogEntry) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		watchman := Cli.Watch(ctx, key)

		var err error
		var LogEntry []*LogEntry
		for wresp := range watchman {
			for _, event := range wresp.Events {
				err = json.Unmarshal(event.Kv.Value, &LogEntry)
				if err != nil {
					log.Fatalf("etcd:[getconf]: Cant unmarshal configure. err:%v", err)
				}
				conf <- LogEntry
			}
		}
		cancel()
	}
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 派一个哨兵, 一直监视这 key 的变化(新增, 修改, 删除)
	//		fmt.Printf("Type:%v Key:%v value: %v \n",event.Type, string(event.Kv.Key), string(event.Kv.Value))
	//	}
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//cancel()
	//resp, err := cli.Put(ctx, "sunheng", "24")
	//if err != nil {
	//	switch err {
	//	case context.Canceled:
	//		log.Fatalf("ctx is canceled by another routine: %v", err)
	//	case context.DeadlineExceeded:
	//		log.Fatalf("ctx is attached with a deadline is exceeded: %v", err)
	//	case rpctypes.ErrEmptyKey:
	//		log.Fatalf("client-side error: %v", err)
	//	default:
	//		log.Fatalf("bad cluster endpoints, which are not etcd servers: %v", err)
	//	}
	//}
	//
}
