package controllers

import (
	"github.com/astaxie/beego"
)

// WebSocketController handles WebSocket requests.
type StopSlaveController struct {
	beego.Controller
}

func (c *StopSlaveController) Get() {
	sess := c.StartSession()
	username := sess.Get("userName")
	if username == nil || username == "" {
		c.TplNames = "login.html"
	} else {
		ip := c.Input().Get("ip")
		currkey := peerGroup.Keys()
		beego.Debug(currkey, ip)
		if CheckState(currkey, ip) == true {
			v := make(map[string]interface{})
			v = map[string]interface{}{"action": ACTION_STOP}
			err := peerGroup.PutMessage(ip, v)
			if err != nil {
				beego.Debug(err)
			}
		} else {
			beego.Debug("sdd")
		}
		c.Redirect("/slavemanager", 302)
	}

}
