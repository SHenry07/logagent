package kafkaclient

import (
	"github.com/Shopify/sarama"
	"log"
)

var (
	// Client kafka client
	Client  sarama.SyncProducer
	MsgChan chan *Row
)

type Row struct {
	Data string
	Topic string
}

// make a cache channel for network latency

// Init initial kafka producer
func Init(addr []string, MaxSize int) (err error) {
	config := sarama.NewConfig()
	// WaitForAll 需要所有的broker返回1
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 成功交付的信息将在success channel 返回
	config.Producer.Return.Successes = true
	Client, err = sarama.NewSyncProducer(addr, config)
	if err != nil {
		return err
	}
	log.Println("kafka: connect success ")
	MsgChan = make(chan *Row, MaxSize)
	// 开启后台的goroutine等待获取通道中的数据
	go SendToKafka()
	return

}

func SentToChan(data string, Topic string) {
	MsgChan <- &Row{Data:data,
		Topic: Topic,}
}

// SendToKafka send Data to kafka
func SendToKafka() {
	for {
		Row := <-MsgChan
		XXX := &sarama.ProducerMessage{Topic: Row.Topic,
			Value: sarama.StringEncoder(Row.Data)}
		partition, offset, err := Client.SendMessage(XXX)
		if err != nil {
			log.Printf("send message failed,err: %v", err)
			return
		}
		log.Printf("partition is %d, offset is %d", partition, offset)
	}
}
