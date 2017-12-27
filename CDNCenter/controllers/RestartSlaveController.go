package controllers

import (
	"github.com/astaxie/beego"
)

// WebSocketController handles WebSocket requests.
type RestartSlaveController struct {
	beego.Controller
}

func (c *RestartSlaveController) Get() {
	sess := c.StartSession()
	username := sess.Get("userName")
	if username == nil || username == "" {
		c.TplNames = "login.html"
	} else {
		ip := c.Input().Get("ip")
		currkey := peerGroup.Keys()
		if CheckState(currkey, ip) == true {
			v := make(map[string]interface{})
			v = map[string]interface{}{"action": ACTION_RESTART}
			peerGroup.PutMessage(ip, v)
		} else {

		}
		c.Redirect("/slavemanager", 302)
	}

}
