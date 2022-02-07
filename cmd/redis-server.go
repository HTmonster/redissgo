/*
 * @Description: redis server
 * @Autor: HTmonster
 * @Date: 2022-01-28 10:33:09
 */

package main

import (
	"fmt"
	"os"

	"github.com/HTmonster/redissgo/internal/config"
	"github.com/HTmonster/redissgo/internal/logger"
	"github.com/HTmonster/redissgo/server"
	"github.com/sevlyar/go-daemon"
)

var banner = `                  ___                     
   ________  ____/ (_)_____________ _____ 
  / ___/ _ \/ __  / / ___/ ___/ __  / __ \  PID:%d
 / /  /  __/ /_/ / (__  |__  ) /_/ / /_/ /  Port:%d
/_/   \___/\__,_/_/____/____/\__, /\____/
                            /____/          %s
 `

var serverProperties = config.Properties
var pid int

func init() {
	// process ID
	pid = os.Getpid()
	// parse configuration file or arg
	serverProperties, _ := config.ParseConfig(os.Args[1:])
	// log output set
	if serverProperties.Logfile != "" {
		logger.SetLogFile(serverProperties.Logfile)
	} else {
		// banner
		fmt.Printf(banner, pid, serverProperties.Port, config.Version)
	}
	// log
	logger.Log.Info("# oO0OoO0OoO0Oo redissgo is starting oO0OoO0OoO0Oo")
}

/**
 * @description: run server in daemon mode
 */
func initDaemon() {
	context := new(daemon.Context)
	child, _ := context.Reborn()
	if child != nil {
		logger.Log.Info("= init daemon succeeded PID=", os.Getpid(), " Exit")
		os.Exit(0)
		return
	} else {
		defer func() {
			if err := context.Release(); err != nil {
				logger.Log.Error("Unable to release pid-file: %s", err.Error())
			}
		}()

		logger.Log.Info("= init daemon succeeded PID=", os.Getpid(), " Born")
	}
}

func main() {

	// Daemonize?
	if serverProperties.Daemonize {
		initDaemon()
	}

	// setup server
	if err := server.SetupAndListen(serverProperties.Bind, serverProperties.Port); err != nil {
		logger.Log.Error(err)
	}

}
