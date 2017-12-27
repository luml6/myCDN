package main

import (
	"controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	_ "github.com/go-sql-driver/mysql"
	// "net/http"
	//"os"
	//"os/signal"
	_ "routers"
)

func init() {
	dbdsn := beego.AppConfig.String("demo::dbdsn")
	dbdriver := beego.AppConfig.String("demo::dbdriver")
	orm.RegisterDataBase("default", dbdriver, dbdsn)
	var err error
	controllers.Bm, err = cache.NewCache("memory", `{"interval":84600}`)
	CheckErr(err)
}

func CheckErr(err error) {
	if err != nil {
		beego.Debug(err)
	}
}
func main() {
	beego.SessionOn = true
	beego.SetLogger("file", `{"filename":"logs/cdntest.log","maxsize":1048576}`)
	tk1 := toolbox.NewTask("tk1", "0 */10 * * * *", Doset)
	// orm.Debug = true
	// beego.BeeLogger.DelLogger("console")
	//go sysSignhandleDemo()
	// http.HandleFunc(controllers.DefaultWSURL, controllers.NewWsHandler())
	toolbox.AddTask("tk1", tk1)
	toolbox.StartTask()
	defer toolbox.StopTask()
	beego.Run()

}
func Doset() (err error) {
	beego.Debug(controllers.Key)
	for i := 0; i < len(controllers.Key); i++ {
		if controllers.Bm.Get(controllers.Key[i]) != nil {
			controllers.Pointer = controllers.Bm.Get(controllers.Key[i]).(*controllers.DateSub)
			controllers.AddLogger(controllers.Key[i])
		}
	}
	controllers.Bm.ClearAll()
	return
}
