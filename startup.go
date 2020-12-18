package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/eviltomorrow/aphrodite-web/cache"

	//
	"net/http"
	_ "net/http/pprof"

	"go.uber.org/zap"

	"github.com/eviltomorrow/aphrodite-base/tools"
	"github.com/eviltomorrow/aphrodite-base/zlog"
	"github.com/eviltomorrow/aphrodite-web/app"
	"github.com/eviltomorrow/aphrodite-web/config"
	"github.com/eviltomorrow/aphrodite-web/db"
)

const (
	nmConfig  = "config"
	nmVersion = "v"
	nmModel   = "model"
)

//
var (
	GitTag    = ""
	BuildTime = ""
)

var (
	path     = flag.String(nmConfig, "config.toml", "配置文件路径")
	version  = flag.Bool(nmVersion, false, "版本信息")
	runModel = flag.String(nmModel, "release", "启动模式")
)

var cfg = config.DefaultGlobalConfig
var cpf []func() error
var globlFileLock tools.FileLock
var pidpath = "aphrodite-web.pid"

func main() {
	defer func() {
		if err := recover(); err != nil {
			zlog.Error("Panic: Unknown reason", zap.Error(fmt.Errorf("%v", err)))
			debug.PrintStack()
			zlog.Error("Stack", zap.String("stack", string(debug.Stack())))
		}
		zlog.Sync()
	}()
	flag.Parse()

	if *version {
		printVersion()
		os.Exit(0)
	}

	if err := cfg.Load(*path, overrideFlags); err != nil {
		log.Printf("Warning: Load config file failure, use default config, nest error: %v\r\n", err)
	}

	setupLogConfig()
	setupGlobalVars()
	printInfo()
	checkpid()
	startupPProfService()
	registerCleanupFunc()
	startupApplication()
	blockingUntilTermination()

}

func setupLogConfig() {
	global, prop, err := zlog.InitLogger(&zlog.Config{
		Level:            cfg.Log.Level,
		Format:           cfg.Log.Format,
		DisableTimestamp: cfg.Log.DisableTimestamp,
		File: zlog.FileLogConfig{
			Filename: cfg.Log.FileName,
			MaxSize:  cfg.Log.MaxSize,
		},
		DisableStacktrace: true,
	})
	if err != nil {
		log.Printf("Fatal: Setup log config failure, nest error: %v", err)
		os.Exit(1)
	}
	zlog.ReplaceGlobals(global, prop)
}

func setupGlobalVars() {
	cache.RedisDSN = cfg.Redis.DSN

	db.MySQLDSN = cfg.MySQL.DSN
	db.MySQLMinOpen = cfg.MySQL.MinOpen
	db.MySQLMaxOpen = cfg.MySQL.MaxOpen

	app.Port = cfg.System.HTTPServerPort
	app.PathHTML = cfg.System.PathHTML
}

func printInfo() {
	zlog.Info("Config information", zap.String("data", cfg.String()))
}

func printVersion() {
	var format = `Git Tag: %s
Build Time: %s
`
	fmt.Printf(format, GitTag, BuildTime)
}

func overrideFlags(cfg *config.Config) {

}

func checkpid() {
	var dir = filepath.Dir(os.Args[0])
	filelock, err := tools.NewFileLock(filepath.Join(dir, pidpath), false)
	if err != nil {
		zlog.Fatal("Check runtime's pid file failure", zap.Error(err))
	}
	globlFileLock = filelock
}

func registerCleanupFunc() {
	db.CloseMySQL()
	cache.CloseRedis()
}

func startupPProfService() {
	var runModel = strings.ToLower(*runModel)
	switch runModel {
	case "debug":
	default:
		return
	}

	go func() {
		if cfg.System.PProfListenPort == 0 {
			zlog.Fatal("PProf service port not configured.")
		}
		zlog.Info("Startup service pprof success.", zap.Int("pprof-port", cfg.System.PProfListenPort), zap.String("visit-page", fmt.Sprintf("http://localhost:%d/debug/pprof/", cfg.System.PProfListenPort)))

		err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", cfg.System.PProfListenPort), nil)
		if err != nil {
			zlog.Fatal("Start service pprof failure", zap.Error(err))
		}
	}()
}

func startupApplication() {
	app.Startup()
	zlog.Info("Start application success.", zap.String("name", "aphrodite-web"), zap.String("model", *runModel))
}

func blockingUntilTermination() {
	var ch = make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	switch <-ch {
	case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
	case syscall.SIGUSR1:
	case syscall.SIGUSR2:
	default:
	}
	for _, f := range cpf {
		f()
	}
	zlog.Info("Termination main programming, cleanup function is executed complete")
}
