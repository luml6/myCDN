package controllers

import (
	"errors"
	"github.com/gorilla/websocket"
	"math/rand"
	"sync"
	"time"
)

var (
	slaveGroup = SlaveGroup{
		m: make(map[string]Slave, 10),
	}
)

type Slave struct {
	Name        string
	Connection  *websocket.Conn
	Master      string
	ConnectTime time.Time
	Line        string
}

type SlaveGroup struct {
	sync.RWMutex
	m map[string]Slave
}

func (sm *SlaveGroup) AddSlave(name, master, line string, conn *websocket.Conn) {
	sm.Lock()
	defer sm.Unlock()
	timenow := time.Now()
	sm.m[name] = Slave{
		Name:        name,
		Connection:  conn,
		Master:      master,
		ConnectTime: timenow,
		Line:        line,
	}
}

func (sm *SlaveGroup) Delete(name string) {
	sm.Lock()

	// sm.m[name].IsIn = false
	delete(sm.m, name)
	sm.Unlock()
}

func (sm *SlaveGroup) Keys() []string {
	sm.RLock()
	defer sm.RUnlock()
	keys := []string{}
	for key, _ := range sm.m {
		keys = append(keys, key)
	}
	return keys
}

func (sm *SlaveGroup) PeekSlave() (string, error) {
	// FIXME(ssx): need to order by active download count
	sm.RLock()
	defer sm.RUnlock()
	ridx := rand.Int()
	keys := []string{}
	for key, _ := range sm.m {
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		return "", errors.New("Slave count zero")
	}
	return keys[ridx%len(keys)], nil
}

func (sm *SlaveGroup) BroadcastJSON(v interface{}) error {
	var err error
	for _, s := range sm.m {
		if err = s.Connection.WriteJSON(v); err != nil {
			return err
		}
	}
	return nil
}
