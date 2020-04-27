package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"logagent/conf"
	"logagent/etcd"
	kafkaclient "logagent/kafkaClient"
	tailf "logagent/tail"
	"logagent/utils"
	"sync"
)

var (
	cfg *conf.AppConf
)

func main() {
	// ini
	cfg = new(conf.AppConf)
	err := ini.MapTo(cfg, "./conf/conf.ini")
	if err != nil {
		log.Panicf("cant load config file, eeror %v\n", err)
	}
	log.Println("ini: Loading server configuration success ")

	// etcd
	endPoints := cfg.EtcdConf.Addr
	err = etcd.Init(endPoints)
	if err != nil {
		log.Panicf("cant connect etcd, eeror %v\n", err)
	}
	log.Println("etcd: connect success ")
	defer etcd.Cli.Close()

	// kafka client
	err = kafkaclient.Init(cfg.KafkaConf.Addr, cfg.KafkaConf.ChanMaxSize)
	if err != nil {
		log.Panicf("Error creating client: %v \n", err)
	}

	// 2.1 从etcd中获取日志收集项的配置信息
	ip := utils.Init()
	key := fmt.Sprintf(cfg.EtcdConf.Key, ip)
	logentries := etcd.GetConf(key)
	log.Println("get conf from etcd success\n")

	// 2.2 监控etcd key的变化 并通知tailf
	//etcd.WatchConf()
	tailf.Init(logentries)

	// 3. 哨兵监控配置变化
	NewConfChan := tailf.NewConfChan()
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConf(key, NewConfChan) // 哨兵发现最新的配置信息h会通知上面的通道
	wg.Wait()

	// TODO: 当第一运行时etcd为空的时候会卡住
}
