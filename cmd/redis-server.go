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
)

var banner = `                  ___                     
   ________  ____/ (_)_____________ _____ 
  / ___/ _ \/ __  / / ___/ ___/ __  / __ \  PID:%d
 / /  /  __/ /_/ / (__  |__  ) /_/ / /_/ /  Port:%d
/_/   \___/\__,_/_/____/____/\__, /\____/
                            /____/          %s
 `

// var serverProperties = config.Properties
var serverProperties *config.ConfigProperties
var pid int

func init() {
	// process ID
	pid = os.Getpid()
	// parse configuration file or arg
	serverProperties, _ := config.ParseConfig(os.Args[1:])
	// log
	logger.Log.Info("# oO0OoO0OoO0Oo redissgo is starting oO0OoO0OoO0Oo")
	// banner
	fmt.Printf(banner, pid, serverProperties.Port, config.Version)
}

func main() {

}
