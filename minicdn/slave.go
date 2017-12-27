package main

import (
	"container/list"
	"sync"
)

var (
	slaveGroup = SlaveGroup{
		m: list.New(),
	}
)

type Slave struct {
	Name string
	Size int
}

type SlaveGroup struct {
	sync.RWMutex
	m *list.List
}

func (sm *SlaveGroup) AddSlave(name string, size int) {
	slave := Slave{
		Name: name,
		Size: size,
	}
	sm.Lock()
	defer sm.Unlock()
	sm.m.PushBack(slave)
}
func (sm *SlaveGroup) Size() (num int) {
	i := sm.m.Len()
	return i
}
func (sm *SlaveGroup) First() (slave Slave) {
	e := sm.m.Front()
	s := e.Value.(Slave)
	return s
}
func (sm *SlaveGroup) Delete(slave Slave) {
	sm.Lock()
	for e := sm.m.Front(); e != nil; e = e.Next() {
		if e.Value == slave {
			sm.m.Remove(e)
		}
	}
	sm.Unlock()
}
