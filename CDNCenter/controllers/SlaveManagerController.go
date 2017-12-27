package controllers

import (
	// "encoding/json"
	// "net"
	// "net/http"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils/pagination"
	// "models"
	"strconv"
	"time"
	// "github.com/gorilla/websocket"
	// "github.com/beego/samples/WebIM/models"
)

// WebSocketController handles WebSocket requests.
type SlaveManagerController struct {
	beego.Controller
}
type SlaveList struct {
	IP       string
	State    bool
	IsRead   bool
	DownTime string
	Master   string
	Line     string
}
type Lists struct {
	Value      string
	IsSelected bool
}

var ProMsg string
var IsShow bool

func (c *SlaveManagerController) Get() {
	sess := c.StartSession()
	username := sess.Get("userName")
	if username == nil || username == "" {
		c.TplNames = "login.html"
	} else {
		pageNum, err := strconv.Atoi(c.Input().Get("p"))
		if err != nil {
			pageNum = 1
		}
		beego.Debug(SlaveKey)
		var list []string
		var lists []Lists
		pagesize := 10
		var start, end, totalCount int
		var cates []SlaveList
		master := c.Input().Get("master")
		if len(master) == 0 || master == "all" {
			totalCount = len(SlaveKey)
			beego.Debug(totalCount)
			i := pageNum * pagesize
			if i < totalCount {
				start = (pageNum - 1) * pagesize
				end = pageNum * pagesize
			} else {
				start = (pageNum - 1) * pagesize
				end = totalCount
			}
			cates = GetSlaveList(start, end)
		} else {
			for i := 0; i < len(SlaveKey); i++ {
				if SlaveKey[i].master == master {
					list = append(list, SlaveKey[i].key)
				}
			}
			totalCount = len(list)
			beego.Debug(totalCount)
			i := pageNum * pagesize
			if i < totalCount {
				start = (pageNum - 1) * pagesize
				end = pageNum * pagesize
			} else {
				start = (pageNum - 1) * pagesize
				end = totalCount
			}
			cates = GetCateslist(start, end, master)
		}
		// currKeys := slaveGroup.Keys()
		//beego.Debug(currKeys)

		for i := 0; i < len(Master); i++ {
			var list Lists
			list.Value = Master[i]
			list.IsSelected = false
			lists = append(lists, list)
		}
		c.Data["Select"] = master
		c.Data["Lists"] = lists
		c.Data["CateList"] = cates
		if totalCount > 0 {
			p := pagination.NewPaginator(c.Ctx.Request, pagesize, totalCount)
			c.Data["paginator"] = p
		}
		c.TplNames = "slavemanager.html"
	}

}
func GetCateslist(start, end int, master string) (lists []SlaveList) {
	var cates []SlaveList
	var list []string
	currkeys := slaveGroup.Keys()
	cKeys := peerGroup.Keys()
	for i := 0; i < len(SlaveKey); i++ {
		if SlaveKey[i].master == master {
			list = append(list, SlaveKey[i].key)
		}
	}
	for k := start; k < end; k++ {
		var cate SlaveList
		cate.IP = list[k]
		if len(cKeys) == 0 {
			cate.State = false
		} else {
			if CheckState(cKeys, list[k]) == true {
				cate.State = true
			} else {
				cate.State = false
			}
		}
		if CheckState(currkeys, list[k]) == true {
			cate.IsRead = true
		}
		if len(slaveGroup.m[cate.IP].Name) == 0 {
			cate.DownTime = "-"
		} else {
			times := time.Now().Sub(slaveGroup.m[cate.IP].ConnectTime)
			beego.Debug(times)
			cate.DownTime = times.String()
		}

		cate.Master = slaveGroup.m[cate.IP].Master
		line := slaveGroup.m[cate.IP].Line
		if line == "d" {
			cate.Line = "电信"
		} else {
			cate.Line = "联通"
		}
		cates = append(cates, cate)
	}
	beego.Debug(cates)
	return cates
}
func GetSlaveList(start, end int) (list []SlaveList) {
	var cates []SlaveList
	currKeys := slaveGroup.Keys()
	cKeys := peerGroup.Keys()
	beego.Debug(start, end)
	for k := start; k < end; k++ {
		var cate SlaveList
		cate.IP = SlaveKey[k].key
		if len(cKeys) == 0 {
			cate.State = false
		} else {
			if CheckState(cKeys, SlaveKey[k].key) == true {
				cate.State = true
			} else {
				cate.State = false
			}
		}
		if CheckState(currKeys, SlaveKey[k].key) == true {
			cate.IsRead = true
		}
		if len(slaveGroup.m[cate.IP].Name) == 0 {
			cate.DownTime = "-"
		} else {
			times := time.Now().Sub(slaveGroup.m[cate.IP].ConnectTime)
			beego.Debug(times)
			cate.DownTime = times.String()
		}
		// cate.DownTime = slaveGroup.m[cate.IP].ConnectTime.Format("2006-01-02 15:04:05")
		cate.Master = slaveGroup.m[cate.IP].Master
		line := slaveGroup.m[cate.IP].Line
		if line == "d" {
			cate.Line = "电信"
		} else {
			cate.Line = "联通"
		}
		cates = append(cates, cate)
	}
	beego.Debug(cates)
	return cates
}
func CheckState(lists []string, list string) (istrue bool) {
	for i := 0; i < len(lists); i++ {
		if lists[i] == list {
			return true
		}
	}
	return false
}
