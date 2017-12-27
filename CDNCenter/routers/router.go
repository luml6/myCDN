package routers

import (
	"controllers"
	"github.com/astaxie/beego"
)

func init() {
	// beego.Router("/Data_HeartBeat", &controllers.HeartBeatController{})
	beego.Router("/ws", &controllers.WebSocketController{})
	beego.Router("/wss", &controllers.DeamonSocketController{})
	beego.Router("/download", &controllers.DownLoadMsgController{})
	beego.Router("/", &controllers.MainController{})
	beego.Router("/slavemanager", &controllers.SlaveManagerController{})
	beego.Router("/stopSlave", &controllers.StopSlaveController{})
	beego.Router("/restartSlave", &controllers.RestartSlaveController{})
	//beego.Router("/V2/Data_HeartBeat", &controllers.V2HeartBeatController{})
	//beego.Router("/Stop", &controllers.HeartBeatController{}, "get:StopHeartBeat")
	//beego.Router("/", &controllers.HeartBeatController{}, "get:StopHeartBeat")
}
