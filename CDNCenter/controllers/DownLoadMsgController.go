package controllers

import (
	// "encoding/json"
	// "net"
	// "net/http"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils/pagination"
	"models"
	"strconv"
	"time"
	// "github.com/gorilla/websocket"
	// "github.com/beego/samples/WebIM/models"
)

// WebSocketController handles WebSocket requests.
type DownLoadMsgController struct {
	beego.Controller
}
type DownloadList struct {
	IP           string
	DownloadName string
	DownloadSize string
	// DownloadCount string
	DownloadTime string
}

// var totalCount int
// var totalSize int

func (c *DownLoadMsgController) Get() {
	sess := c.StartSession()
	username := sess.Get("userName")
	if username == nil || username == "" {
		c.TplNames = "login.html"
	} else {
		pageNum, err := strconv.Atoi(c.Input().Get("p"))
		if err != nil {
			pageNum = 1
		}
		var total []models.TDownload
		var cates []DownloadList
		pagesize := 10
		var totalCount int
		startTime := c.Input().Get("startTime")
		endTime := c.Input().Get("endTime")
		ip := c.Input().Get("ip")
		if len(startTime) != 0 && len(endTime) != 0 {
			c.Data["StartTime"] = startTime
			c.Data["EndTime"] = endTime
			startTime, endTime = ChangeTime(startTime, endTime)
			totalCount, total = models.FindTDownloadData(startTime, endTime)
			_, maps := models.FindTDownloadWithDatePage(pageNum, pagesize, totalCount, startTime, endTime)
			cates = GetTDownload(maps)
			c.GetTotal(total)
			// gamelist = GetGamelist(maps)
		} else if len(ip) != 0 {
			c.Data["IPname"] = ip
			totalCount, total = models.FindTDownloadIP(ip)
			_, maps := models.FindTDownloadWithIPPage(pageNum, pagesize, totalCount, ip)
			cates = GetTDownload(maps)
			c.GetTotal(total)
		} else {
			totalCount, total = models.FindTDownloadAll()
			_, ch := models.FindTDownloadWithPage(pageNum, pagesize, totalCount)
			cates = GetTDownload(ch)
			c.GetTotal(total)
		}
		c.Data["CateList"] = cates
		if totalCount > 0 {
			p := pagination.NewPaginator(c.Ctx.Request, pagesize, totalCount)
			c.Data["paginator"] = p
		}
		c.TplNames = "download.html"
	}

}
func ChangeTime(startTime, endTime string) (tm3, tm4 string) {
	loc, _ := time.LoadLocation("Local")                         //重要：获取时区
	tm1, _ := time.ParseInLocation("2006-01-02", startTime, loc) //使用模板在对应时区转化为time.time类型
	tm2, _ := time.ParseInLocation("2006-01-02", endTime, loc)
	d, _ := time.ParseDuration("24h")
	if tm1.Unix() > tm2.Unix() {
		tm3 := tm1.Add(d)
		tm1 = tm2
		tm2 = tm3
	} else {
		tm3 := tm2.Add(d)
		tm2 = tm3
	}
	startTime = tm1.Format("2006-01-02")
	endTime = tm2.Format("2006-01-02")
	return startTime, endTime
}
func (c *DownLoadMsgController) GetTotal(maps []models.TDownload) {
	var size int
	for i := 0; i < len(maps); i++ {
		size += maps[i].DownSize
	}
	c.Data["TotalCount"] = len(maps)
	c.Data["TotalSize"] = strconv.Itoa(size/1024) + "M"
}
func GetTDownload(maps []models.TDownload) (ch []DownloadList) {
	var cates []DownloadList
	for i := 0; i < len(maps); i++ {
		var cate DownloadList
		cate.IP = maps[i].Ip
		cate.DownloadName = maps[i].Download
		size := maps[i].DownSize
		cate.DownloadSize = strconv.Itoa(size)
		cate.DownloadTime = maps[i].DownloadTime.Format("2006-01-02")
		cates = append(cates, cate)
	}
	return cates
}
