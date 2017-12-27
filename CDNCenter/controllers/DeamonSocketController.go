package controllers

import (
	// "encoding/json"
	"net"
	"net/http"
	// "strconv"
	// "fmt"
	// "strings"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	// "github.com/beego/samples/WebIM/models"
)

// WebSocketController handles WebSocket requests.
type DeamonSocketController struct {
	beego.Controller
}

var DeamonKey []string

// Join method handles WebSocket requests for WebSocketController.
func (this *DeamonSocketController) Get() {

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
		case ACTION_LOGIN:
			// ip = remoteHost
			// currKeys := peerGroup.Keys()
			peerGroup.AddPeer(ip, ws)
			// istrue := CheckKey(Key, ip)
			// if istrue != true {
			// 	DeamonKey = append(DeamonKey, ip)
			// }
			beego.Debug("Slave Deamon JOIN:", ip)
		default:
			beego.Debug("UNKNOWN:", msg)
		}
		if err != nil {
			break
		}
	}
	peerGroup.Delete(ip)
	beego.Debug("SlaveDeamon QUIT:", ip)
}
