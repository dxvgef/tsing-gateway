package etcd

import (
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/rs/zerolog/log"

	"local/global"
)

type Etcd struct {
	client               *clientv3.Client
	ClientID             string   `json:"-"`
	KeyPrefix            string   `json:"key_prefix"`
	Endpoints            []string `json:"endpoints"`
	DialTimeout          uint     `json:"dial_timeout,omitempty"`
	Username             string   `json:"username,omitempty"`
	Password             string   `json:"password,omitempty"`
	AutoSyncInterval     uint     `json:"auto_sync_interval,omitempty"`
	DialKeepAliveTime    uint     `json:"dial_keep_alive_time,omitempty"`
	DialKeepAliveTimeout uint     `json:"dial_keep_alive_timeout,omitempty"`
	MaxCallSendMsgSize   uint     `json:"max_call_send_msg_size,omitempty"`
	MaxCallRecvMsgSize   uint     `json:"max_call_recv_msg_size,omitempty"`
	RejectOldCluster     bool     `json:"reject_old_cluster,omitempty"`
	PermitWithoutStream  bool     `json:"permit_without_stream,omitempty"`
}

func New(config string) (*Etcd, error) {
	var instance Etcd
	instance.ClientID = global.SnowflakeNode.Generate().String()

	err := instance.UnmarshalJSON(global.StrToBytes(config))
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}

	instance.client, err = clientv3.New(clientv3.Config{
		Endpoints:            instance.Endpoints,
		DialTimeout:          time.Duration(instance.DialTimeout) * time.Second,
		Username:             instance.Username,
		Password:             instance.Password,
		AutoSyncInterval:     time.Duration(instance.AutoSyncInterval) * time.Second,
		DialKeepAliveTime:    time.Duration(instance.DialKeepAliveTime) * time.Second,
		DialKeepAliveTimeout: time.Duration(instance.DialKeepAliveTimeout) * time.Second,
		MaxCallSendMsgSize:   int(instance.MaxCallSendMsgSize),
		MaxCallRecvMsgSize:   int(instance.MaxCallRecvMsgSize),
		RejectOldCluster:     instance.RejectOldCluster,
		PermitWithoutStream:  instance.PermitWithoutStream,
	})
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}
	return &instance, nil
}
