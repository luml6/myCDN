package controllers

import (
	// "encoding/json"
	"net"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

// WebSocketController handles WebSocket requests.
type WebSocketController struct {
	beego.Controller
}
type Slaves struct {
	key    string
	master string
}

var SlaveKey []Slaves
var Master []string

// Join method handles WebSocket requests for WebSocketController.
func (this *WebSocketController) Get() {

	// Upgrade from http request to WebSocket.
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)

	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}
	var ip string
	remoteHost, _, _ := net.SplitHostPort(ws.RemoteAddr().String())
	ip = remoteHost
	// Message receive loop.
	for {
		var msg = make(map[string]interface{})
		if err = ws.ReadJSON(&msg); err != nil {
			break
		}

		action, _ := msg["action"].(string)
		switch action {
		case ACTION_RECEIVE:
			beego.Debug(msg)
			name, _ := msg["name"].(string)
			size, _ := msg["size"].(string)
			if len(name) != 0 && len(size) != 0 {
				Size, _ := strconv.Atoi(size)
				mac := Bm.Get(ip)
				istrue := CheckKey(Key, ip)
				if istrue != true {
					Key = append(Key, ip)
				}
				beego.Debug(Key)
				if mac == nil || mac == "" {
					AddCache(name, ip, Size)
				} else {
					Pointer = mac.(*DateSub)
					ChangeCache(Size, name, ip)
				}
			}
		case ACTION_LOGIN:
			master, _ := msg["master"].(string)
			line, _ := msg["line"].(string)
			ip = remoteHost
			if master == "master" {
				istrue := CheckKey(Master, ip)
				if istrue != true {
					Master = append(Master, ip)
				}
				istrue = CheckSalve(ip, ip)
				if istrue != true {
					var slave Slaves
					slave.key = ip
					slave.master = ip
					SlaveKey = append(SlaveKey, slave)
				}
				slaveGroup.AddSlave(ip, ip, line, ws)
			} else {
				istrue := CheckKey(Master, master)
				if istrue != true {
					Master = append(Master, master)
				}
				istrue = CheckSalve(ip, master)
				if istrue != true {
					var slave Slaves
					slave.key = ip
					slave.master = master
					SlaveKey = append(SlaveKey, slave)
				}
				slaveGroup.AddSlave(ip, master, line, ws)
			}
			// currKeys := slaveGroup.Keys()

			beego.Debug("Slave JOIN:", ip)
		default:
			beego.Debug("UNKNOWN:", msg)
		}
		if err != nil {
			break
		}
		// publish <- newEvent(models.EVENT_MESSAGE, uname, string(p))
	}
	slaveGroup.Delete(ip)
	beego.Debug("Slave QUIT:", ip)
}
func CheckKey(key []string, name string) (istrue bool) {
	for i := 0; i < len(key); i++ {
		if key[i] == name {
			return true
		}
	}
	return false
}
func CheckSalve(key, master string) (istrue bool) {
	for i := 0; i < len(SlaveKey); i++ {
		if key == SlaveKey[i].key && master == SlaveKey[i].master {
			return true
		}
	}
	return false
}
