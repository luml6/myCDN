package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/widuu/goini"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	wsclient *websocket.Conn
	IsStop   bool = false
	IsClosed bool = false
	IsStart  bool = false
)

const (
	ACTION_LOGIN   = "login"
	ACTION_RESTART = "restart"
	ACTION_STOP    = "stop"
)

var Cmd *exec.Cmd
var commond string

func main() {
	lf, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		os.Exit(1)
	}
	defer lf.Close()

	// 日志
	l := log.New(lf, "", os.O_APPEND)
	conf := goini.SetConfig("./conf/config.ini")
	cmirror := conf.GetValue("Center", "mirror")
	//typeStyle := conf.GetValue("Select", "type")
	commond = conf.GetValue("Cmd", "cmd")
	//StartCmd(cmd, l)
	if err := InitCenter(cmirror, l); err != nil {
		fmt.Println(err)
	}
	http.HandleFunc("/", NewFileHandler)
	http.ListenAndServe(":8888", nil)
}
func NewFileHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
func StartCmd(log *log.Logger) {
	for {
		cmd := exec.Command(commond)
		Cmd = cmd
		if IsStop {
			continue
		}
		err := cmd.Start()
		if err != nil {
			log.Printf("%s 启动命令失败", time.Now().Format("2006-01-02 15:04:05"), err)

			time.Sleep(time.Second * 5)
			continue
		}
		log.Printf("%s 进程启动", time.Now().Format("2006-01-02 15:04:05"), err)
		err = cmd.Wait()
		log.Printf("%s 进程退出", time.Now().Format("2006-01-02 15:04:05"), err)
		time.Sleep(time.Second * 2)
	}

}
func InitCenter(mirror string, log *log.Logger) (err error) {
	u, err := url.Parse(mirror)
	if err != nil {
		log.Println(err)
		return
	}
	u.Path = "/wss"

	conn, err := net.Dial("tcp", u.Host)
	if err != nil {
		log.Println(err)
		return
	}
	wsclient, _, err = websocket.NewClient(conn, u, nil, 1024, 1024)
	if err != nil {
		return
	}
	fmt.Println("连接成功")
	var msg = make(map[string]interface{})
	wsclient.WriteJSON(map[string]string{
		"action": ACTION_LOGIN,
	})
	go func() {
		for {
			err = wsclient.ReadJSON(&msg)
			if err != nil {
				log.Println("Connection to Center closed !!!")
				for {
					log.Println("> retry in 5 seconds")
					time.Sleep(time.Second * 5)
					if err := InitCenter(mirror, log); err == nil {
						break
					}
				}
				break
			} else {
				log.Println(msg)
				action := msg["action"]
				switch action {
				case ACTION_RESTART:
					if !IsClosed {
						Cmd.Process.Kill()
					}
					IsClosed = false
					IsStop = false
				case ACTION_STOP:
					Cmd.Process.Kill()
					IsStop = true
					IsClosed = true
				default:
					log.Println("命令错误")
				}
			}
		}

	}()
	go func() {
		if !IsStart {
			IsStart = true
			StartCmd(log)
		}

	}()
	return err
}
