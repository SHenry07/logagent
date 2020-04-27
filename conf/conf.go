package conf

// AppConf project config set
type AppConf struct {
	EtcdConf `ini:"etcd"`
	KafkaConf `ini:"kafka"`
	TaillogConf `ini:"tailf"`
}

type EtcdConf struct {
	Addr []string `ini:"addr,omitempty,allowshadow"`
	Key string `ini:"key"`
}

// KafkaConf kafka client configure
type KafkaConf struct {
	Addr []string `ini:"addr,omitempty,allowshadow"`
	ChanMaxSize int `ini: "chanmaxsize"`
}

// TaillogConf The files are watched
type TaillogConf struct {
	FileName string `ini:"filename"`
}