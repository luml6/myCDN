package main

import (
	//"crypto/tls"
	//"flag"
	"fmt"
	//"github.com/alecthomas/kingpin"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"log"

	"github.com/widuu/goini"
)

var (
	VERSION = "0.1"

	wslog     *log.Logger
	mOpts     MasterOptions
	sOpts     SlaveOptions
	SlaveAddr string
	Cmirror   string
	Line      string
)

type MasterOptions struct {
	LogFile    string
	MirrorAddr *url.URL
	CacheDir   string
	ListenAddr string
	Secret     string
}

type SlaveOptions struct {
	Secret     string
	CacheDir   string
	ListenAddr string
	MasterAddr *url.URL
}
type UploadMessage struct {
	UploadIp   string
	UploadSize int
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func InitSignal() {
	sig := make(chan os.Signal, 10)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for {
			s := <-sig
			fmt.Println("Got signal:", s)
			if state.IsClosed() {
				fmt.Println("Cold close !!!")
				os.Exit(1)
			}
			fmt.Println("Warm close, waiting ...")

			go func() {
				state.Close()
				os.Exit(0)
			}()
		}
	}()
}

func checkErr(er error) {
	if er != nil {
		log.Fatal(er)
	}
}

func runMaster(opts MasterOptions) {
	logfile := opts.LogFile
	if logfile == "-" || logfile == "" {
		wslog = log.New(os.Stderr, "", log.LstdFlags)
	} else {
		fd, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			log.Fatal(err)
		}
		wslog = log.New(fd, "", 0)
	}

	http.HandleFunc(defaultWSURL, NewWsHandler(opts.MirrorAddr.String(), wslog))
	http.HandleFunc("/", NewFileHandler(true, opts.MirrorAddr.String(), opts.CacheDir))
	http.HandleFunc("/_log", func(w http.ResponseWriter, r *http.Request) {
		if logfile == "" || logfile == "-" {
			http.Error(w, "Log file not found", 404)
			return
		}
		http.ServeFile(w, r, logfile)
	})
	if err1 := InitCenter(Cmirror, "master"); err1 != nil {
		log.Fatal(err1)
	}
	log.Printf("Listening on %s", opts.ListenAddr)
	InitSignal()
	//s := &http.Server{Addr: opts.ListenAddr, TLSConfig: &tls.Config{InsecureSkipVerify: true}}
	log.Fatal(http.ListenAndServe(opts.ListenAddr, nil))
	//log.Fatal(s.ListenAndServeTLS("server.crt", "server.key"))
	//log.Fatal(http.ListenAndServeTLS(opts.ListenAddr, "server.crt", "server.key", nil))

}

func runSlave(opts SlaveOptions) {
	if err := InitPeer(opts.MasterAddr.String(), opts.ListenAddr, opts.CacheDir, opts.Secret); err != nil {
		log.Fatal(err)
	}
	//if err1 := InitCenter(Cmirror, "master"); err1 != nil {
	//	log.Fatal(err1)
	//}
	log.Printf("Listening on %s", opts.ListenAddr)
	InitSignal()
	log.Fatal(http.ListenAndServe(opts.ListenAddr, nil))
	log.Println("success")
	//log.Fatal(http.ListenAndServeTLS(opts.ListenAddr, "server.crt", "server.key", nil))
}

func main() {
	conf := goini.SetConfig("./conf/config.ini")
	mOpts.CacheDir = conf.GetValue("Master", "cachedir")
	mOpts.Secret = conf.GetValue("Master", "secret")
	mirror := conf.GetValue("Master", "mirror")
	mOpts.LogFile = conf.GetValue("Master", "log")
	mOpts.ListenAddr = conf.GetValue("Master", "addr")
	mOpts.MirrorAddr, _ = url.Parse(mirror)
	sOpts.CacheDir = mOpts.CacheDir
	sOpts.Secret = mOpts.Secret
	smirror := conf.GetValue("Slave", "mirror")
	//sOpts.LogFile = mOpts.LogFile
	sOpts.ListenAddr = conf.GetValue("Slave", "addr")
	sOpts.MasterAddr, _ = url.Parse(smirror)
	typeStyle := conf.GetValue("Select", "type")
	SlaveAddr = conf.GetValue("Listen", "addr")
	Cmirror = conf.GetValue("Center", "mirror")
	Line = conf.GetValue("Line", "line")
	switch typeStyle {
	case "Master":
		runMaster(mOpts)
	case "Slave":
		runSlave(sOpts)
	default:
		log.Fatalf("Unknown command: %s", typeStyle)
	}
	// app := kingpin.New("minicdn", "Master node which manage slaves")

	// app.Flag("cachedir", "Cache file directory").Short('d').Default("cache").ExistingDirVar(&mOpts.CacheDir) //.StringVar(&mOpts.CacheDir)
	// app.Flag("secret", "Secret key for server").Short('s').Default("sandy mandy").StringVar(&mOpts.Secret)
	// ma := app.Command("master", "CDN Master")

	// ma.Flag("addr", "Listen address").Default(":7010").StringVar(&mOpts.ListenAddr)
	// ma.Flag("log", "Log file, - for stdout").Short('l').Default("-").StringVar(&mOpts.LogFile)

	// ma.Flag("mirror", "Mirror http address, ex: http://t.co/").Required().URLVar(&mOpts.MirrorAddr)

	// sa := app.Command("slave", "Slave node")
	// sa.Flag("addr", "Listen address").Default(":7020").StringVar(&sOpts.ListenAddr)
	// sa.Flag("maddr", "Master server address, ex: localhost:7010").Short('m').Required().URLVar(&sOpts.MasterAddr)

	// app.Version(VERSION).Author("codeskyblue")
	// app.HelpFlag.Short('h')
	// app.VersionFlag.Short('v')
	// kingpin.CommandLine.Help = "The very simple and mini CDN"
	// // parse command line
	// cmdName := kingpin.MustParse(app.Parse(os.Args[1:]))

	//sOpts.CacheDir = mOpts.CacheDir
	//sOpts.Secret = mOpts.Secret
	//fmt.Println(mOpts)
	//fmt.Println(sOpts)
	//switch cmdName {
	//case ma.FullCommand():
	//	runMaster(mOpts)
	//case sa.FullCommand():
	//	runSlave(sOpts)
	//default:
	//	log.Fatalf("  command: %s", cmdName)
	//}
}
