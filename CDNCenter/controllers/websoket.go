package controllers

import (
	// "encoding/json"
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/gorilla/websocket"
	// "net/url"
	"models"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	wsclient *websocket.Conn
	sendc    = make(chan map[string]interface{}, 10)
)

const (
	ACTION_LOGIN       = "login"
	ACTION_RECEIVE     = "receive"
	ACTION_PEER_UPDATE = "peer_update"
	ACTION_RESTART     = "restart"
	ACTION_STOP        = "stop"
)

var Key []string

type Date struct {
	Name string
	size int
}
type DateSub struct {
	IP       string
	dateList []Date
}

const (
	f_date = "2006-01-02" //长日期格式
)

var wmutex sync.Mutex // 写操作需要用到的互斥锁。
var rmutex sync.Mutex // 读操作需要用到的互斥锁。
var Bm cache.Cache
var Pointer *DateSub

func NewWsHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(conn.RemoteAddr())
		defer conn.Close()

		var ip string
		remoteHost, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
		for {
			var err error
			var msg = make(map[string]interface{})
			if err = conn.ReadJSON(&msg); err != nil {
				break
			}

			action, _ := msg["action"].(string)
			beego.Debug(msg)
			switch action {
			case ACTION_LOGIN:
				ip = "http://" + remoteHost + ":8080"
				currKeys := peerGroup.Keys()
				peerGroup.AddPeer(ip, conn)
				err = conn.WriteJSON(map[string]string{
					"self":  ip,
					"peers": strings.Join(peerGroup.Keys(), ","),
				})

				peerGroup.RLock()
				for _, key := range currKeys {
					if s, exists := peerGroup.m[key]; exists {
						s.Connection.WriteJSON(map[string]string{
							"action": ACTION_PEER_UPDATE,
							"peers":  strings.Join(peerGroup.Keys(), ","),
						})
					}
				}
				peerGroup.RUnlock()
				fmt.Printf("Peer: %s JOIN", ip)
			case ACTION_RECEIVE:
				name, _ := msg["name"].(string)
				size, _ := msg["size"].(string)
				Size, _ := strconv.Atoi(size)
				beego.Debug(name, size)
				mac := Bm.Get(ip)
				istrue := CheckKey(Key, ip)
				if istrue != true {
					Key = append(Key, ip)
				}
				if mac == nil || mac == "" {
					AddCache(name, ip, Size)
				} else {
					Pointer = mac.(*DateSub)
					ChangeCache(Size, name, ip)
				}
				// delete(msg, "action")
				// msgb, _ := json.Marshal(map[string]interface{}{
				// 	"timestamp": time.Now().Unix(),
				// 	"data":      msg,
				// 	"peer":      ip,
				// })
				// wslog.Println(string(msgb))
			default:
				fmt.Println("UNKNOWN:", msg)
			}
			if err != nil {
				break
			}
		}

		peerGroup.Delete(ip)
		peerGroup.RLock()
		for _, key := range peerGroup.Keys() {
			if s, exists := peerGroup.m[key]; exists {
				s.Connection.WriteJSON(map[string]string{
					"action": ACTION_PEER_UPDATE,
					"peers":  strings.Join(peerGroup.Keys(), ","),
				})
			}
		}
		peerGroup.RUnlock()
		fmt.Printf("Peer: %s QUIT", ip)
	}
}
func AddCache(name, ip string, size int) {
	// var tm DateTime
	wmutex.Lock()
	defer wmutex.Unlock()
	var date Date
	data := DateSub{}
	data.IP = ip
	date.Name = name
	date.size = size
	data.dateList = append(data.dateList, date)
	Pointer = &data
	Bm.Put(ip, Pointer, 84600)
}
func ChangeCache(size int, name, ip string) {
	rmutex.Lock()
	defer rmutex.Unlock()
	// timeNow := time.Now()
	var date Date
	date.Name = name
	date.size = size
	Pointer.dateList = append(Pointer.dateList, date)
	// istrue := CheckTime(Pointer.OperationTime, timeNow)
	// if istrue == true {
	// 	AddLogger(ip)
	// 	Bm.Delete(ip)
	// }
}

//检查是否到时间
// func CheckTime(startTime, endTime time.Time) (istrue bool) {
// 	d := endTime.Sub(startTime)
// 	beego.Debug(d)
// 	if d > 1*time.Hour {
// 		return true
// 	}
// 	return false
// }

//添加数据到库
func AddLogger(name string) {
	wmutex.Lock()
	defer wmutex.Unlock()
	timeNow := time.Now()
	var logs []models.TDownload
	beego.Debug(Pointer.dateList)
	for i := 0; i < len(Pointer.dateList); i++ {
		var log models.TDownload
		log.Ip = name
		log.Download = Pointer.dateList[i].Name
		log.DownSize = Pointer.dateList[i].size / 1024
		log.DownloadTime = timeNow
		logs = append(logs, log)
	}
	err := models.AddAllDownload(logs)
	if err != nil {
		beego.Debug(err)
	}
}
