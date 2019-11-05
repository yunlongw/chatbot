package bot

import "log"

type ChanMap struct {
	m map[int]chan bool
}

func NewChanMap() *ChanMap {
	return &ChanMap{m: make(map[int]chan bool)}
}

func (cm *ChanMap) SetChan(key int)  {
	log.Printf("初始化:%d", key)
	cm.m[key] = make(chan bool)
}

func (cm *ChanMap) Chan(key int, b bool) {
	log.Printf("赋值:%d, %v", key, b)
	cm.m[key] <- b
}