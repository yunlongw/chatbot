package bot

import (
	"fmt"
	"log"
)

type GroupVerifyChanMap struct {
	m map[int64]map[int]chan bool
}

func NewGroupVerifyChanMap() *GroupVerifyChanMap {
	return &GroupVerifyChanMap{m: make(map[int64]map[int]chan bool)}
}

func (cm *GroupVerifyChanMap) SetChan(GroupID int64, key int) {
	log.Printf("初始化:%d,%d", GroupID, key)
	if _, ok := cm.m[GroupID]; ok {
		if _, o := cm.m[GroupID][key]; o != true {
			cm.m[GroupID][key] = make(chan bool)
		}
	} else {
		boils := make(map[int]chan bool)
		boils[key] = make(chan bool)
		cm.m[GroupID] = boils
	}
	fmt.Println(cm)
}

func (cm *GroupVerifyChanMap) Chan(GroupID int64, key int, b bool) {
	cm.m[GroupID][key] <- b
}

func (cm *GroupVerifyChanMap) DeleteChan(GroupID int64, key int) {
	delete(cm.m[GroupID], key)
	fmt.Println(cm)
}
