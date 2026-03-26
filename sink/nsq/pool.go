package nsq

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	goNsq "github.com/nsqio/go-nsq"
)

// NsdPool 管理 NSQ Producer 连接池，支持多 pool 和按 topic 路由
type NsdPool struct {
	locker     sync.Mutex
	pools      map[string][]*goNsq.Producer // addr → producers
	names      map[string][]string          // poolName → addr 列表
	topicNames map[string]string            // topicName → poolName
	indexs     map[string]int               // "table.topic" → 当前轮询 index
	debug      bool
}

func NewNsdPool(nsqAddrs map[string][]string, topics map[string]string, debug bool) *NsdPool {
	rand.Seed(time.Now().UnixNano())
	p := &NsdPool{
		pools:      make(map[string][]*goNsq.Producer),
		names:      make(map[string][]string),
		topicNames: make(map[string]string),
		indexs:     make(map[string]int),
		debug:      debug,
	}
	for name, addrs := range nsqAddrs {
		p.names[name] = addrs
	}
	for topicName, poolName := range topics {
		p.topicNames[topicName] = poolName
	}
	go p.keepAlive()
	return p
}

func (p *NsdPool) keepAlive() {
	t := time.NewTicker(time.Second * 5)
	for range t.C {
		p.locker.Lock()
		for _, producers := range p.pools {
			for _, prod := range producers {
				prod.Ping()
			}
		}
		p.locker.Unlock()
	}
}

func (p *NsdPool) add(addr string) (*goNsq.Producer, error) {
	p.locker.Lock()
	defer p.locker.Unlock()

	if _, ok := p.pools[addr]; !ok {
		p.pools[addr] = make([]*goNsq.Producer, 0, 5)
		for i := 0; i < 5; i++ {
			config := goNsq.NewConfig()
			config.HeartbeatInterval = time.Second * 5
			producer, err := goNsq.NewProducer(addr, config)
			if err != nil {
				return nil, err
			}
			p.pools[addr] = append(p.pools[addr], producer)
		}
	}
	x := rand.Intn(len(p.pools[addr]))
	return p.pools[addr][x], nil
}

func (p *NsdPool) reconnect(addr string) {
	p.locker.Lock()
	if ps, ok := p.pools[addr]; ok {
		for _, prod := range ps {
			prod.Stop()
		}
		delete(p.pools, addr)
	}
	p.locker.Unlock()
	p.add(addr)
}

func (p *NsdPool) get(tableName string, topic string) (*goNsq.Producer, string, error) {
	poolName, ok := p.topicNames[topic]
	if !ok {
		return nil, "", fmt.Errorf("no pool configured for topic %s (table %s)", topic, tableName)
	}
	addrs, ok := p.names[poolName]
	if !ok {
		return nil, "", fmt.Errorf("no addrs for pool %s", poolName)
	}

	key := fmt.Sprintf("%s.%s", tableName, topic)
	p.locker.Lock()
	if _, ok := p.indexs[key]; ok {
		p.indexs[key] = (p.indexs[key] + 1) % len(addrs)
	} else {
		p.indexs[key] = 0
	}
	idx := p.indexs[key]
	p.locker.Unlock()

	addr := addrs[idx]
	producer, err := p.add(addr)
	if err != nil {
		return nil, "", err
	}
	return producer, addr, nil
}

// Send 发布消息到指定 topic，失败自动重连重试
func (p *NsdPool) Send(tableName string, topic string, data []byte) error {
	producer, addr, err := p.get(tableName, topic)
	if err != nil {
		log.Printf("[NsdPool] get producer error: %v", err)
		return err
	}
	for i := 0; ; i++ {
		err = producer.Publish(topic, data)
		if err == nil {
			if p.debug {
				log.Printf("[NsdPool] %s → %s (%s) %d bytes", tableName, topic, addr, len(data))
			}
			return nil
		}
		p.reconnect(addr)
		log.Printf("[NsdPool] send retry %d topic=%s err=%v", i, topic, err)
		if strings.HasPrefix(err.Error(), "E_BAD_MESSAGE") {
			log.Printf("[NsdPool] E_BAD_MESSAGE, drop msg topic=%s", topic)
			return nil
		}
		time.Sleep(time.Millisecond * 1000)
	}
}

func (p *NsdPool) Stop() {
	p.locker.Lock()
	defer p.locker.Unlock()
	for addr, ps := range p.pools {
		for _, prod := range ps {
			prod.Stop()
		}
		log.Printf("[NsdPool] stopped addr=%s", addr)
	}
}
