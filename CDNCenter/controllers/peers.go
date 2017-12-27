package controllers

import (
	"errors"
	"math/rand"
	"sync"

	"github.com/codeskyblue/groupcache"
	"github.com/gorilla/websocket"
)

var (
	peerGroup = PeerGroup{
		m: make(map[string]Peer, 10),
	}

	pool *groupcache.HTTPPool
)

type Peer struct {
	Name       string
	Connection *websocket.Conn
	// IsIn       bool
}

type PeerGroup struct {
	sync.RWMutex
	m map[string]Peer
}

func (sm *PeerGroup) AddPeer(name string, conn *websocket.Conn) {
	sm.Lock()
	defer sm.Unlock()
	sm.m[name] = Peer{
		Name:       name,
		Connection: conn,
		// IsIn:       true,
	}
}

func (sm *PeerGroup) Delete(name string) {
	sm.Lock()
	// sm.m[name].IsIn = false
	delete(sm.m, name)
	sm.Unlock()
}

func (sm *PeerGroup) Keys() []string {
	sm.RLock()
	defer sm.RUnlock()
	keys := []string{}
	for key, _ := range sm.m {
		keys = append(keys, key)
	}
	return keys
}

func (sm *PeerGroup) PeekPeer() (string, error) {
	// FIXME(ssx): need to order by active download count
	sm.RLock()
	defer sm.RUnlock()
	ridx := rand.Int()
	keys := []string{}
	for key, _ := range sm.m {
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		return "", errors.New("Peer count zero")
	}
	return keys[ridx%len(keys)], nil
}

func (sm *PeerGroup) BroadcastJSON(v interface{}) error {
	var err error
	for _, s := range sm.m {
		if err = s.Connection.WriteJSON(v); err != nil {
			return err
		}
	}
	return nil
}
func (sm *PeerGroup) PutMessage(name string, v interface{}) error {
	var err error
	if err = sm.m[name].Connection.WriteJSON(v); err != nil {
		return err
	}
	return nil
}
