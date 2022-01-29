/*
 * @Description: redis config parse
 * @Autor: HTmonster
 * @Date: 2022-01-28 21:44:33
 */

package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/HTmonster/redissgo/internal/logger"
)

//TODO: add more item
// redis server configuration structure
type ConfigProperties struct {
	Bind      string `json:"bind"`      //e.g. bind 127.0.0.1 ::1
	Port      int    `json:"port"`      //e.g. port 6379
	Timeout   int    `json:"timeout"`   //e.g. timeout 0
	Daemonize bool   `json:"daemonize"` //e.g. daemonize yes
	Logfile   string `json:"logfile"`   //e.g. logfile /var/log/redis/redis-server.log
	Database  int    `json:"database"`  //e.g. databases 16
}

// global vars
var Properties *ConfigProperties
var Version string = "version 1.0 Beta"
var usage string = `Usage: ./redis-server [/path/to/redis.conf] [options]
	./redis-server - (read config from stdin)
	./redis-server -v or --version
	./redis-server -h or --help

Examples:
	./redis-server (run the server with default conf)
	./redis-server /etc/redis/6379.conf
	./redis-server --port 7777
`

// default properties
func init() {
	Properties = &ConfigProperties{
		Bind:     "127.0.0.1",
		Port:     6379,
		Timeout:  0,
		Database: 16,
	}
}

/**
 * @description: parse config from file
 * @param {string} file
 * @return {*}
 */
func parseConfigFile(file string) (*ConfigProperties, error) {
	logger.Log.Info("# loading configuration from file " + file)

	// open configuration file
	f, err := os.Open(file)
	if err != nil {
		logger.Log.Fatal("Could not open config file: \n\t", err)
	}
	defer f.Close()

	// read configuration file
	reader := bufio.NewReader(f)

	// read line by line
	configMap := make(map[string]string)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.Log.Warn("* read line error ", err)
			continue
		}
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		// key and value
		fileds := strings.Fields(string(line))
		key, value := fileds[0], fileds[1]

		configMap[key] = value
	}

	// parse config
	refType := reflect.TypeOf(Properties)
	refValue := reflect.ValueOf(Properties)

	for i := 0; i < refType.Elem().NumField(); i++ {
		filed := refType.Elem().Field(i)
		filedValue := refValue.Elem().Field(i)

		key, ok := filed.Tag.Lookup("json")
		if !ok {
			key = filed.Name
		}
		value, ok := configMap[key]
		if ok {
			switch filed.Type.Kind() {
			case reflect.String:
				filedValue.SetString(value)

			case reflect.Int:
				if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
					filedValue.SetInt(intValue)
				}

			case reflect.Bool:
				boolValue := value == "yes"
				filedValue.SetBool(boolValue)
				//TODO: more type
			}
		}
	}

	return Properties, nil
}

/**
 * @description: parse config from stdin args
 * @param {*} key
 * @param {string} value
 * @return {*}
 */
func parseConfigArg(key, value string) (*ConfigProperties, error) {

	refType := reflect.TypeOf(Properties)
	refValue := reflect.ValueOf(Properties)

	for i := 0; i < refType.Elem().NumField(); i++ {
		filed := refType.Elem().Field(i)
		filedValue := refValue.Elem().Field(i)

		refKey, ok := filed.Tag.Lookup("json")
		if !ok {
			refKey = filed.Name
		}

		if refKey == key {
			switch filed.Type.Kind() {
			case reflect.String:
				filedValue.SetString(value)

			case reflect.Int:
				if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
					filedValue.SetInt(intValue)
				}

			case reflect.Bool:
				boolValue := value == "yes"
				filedValue.SetBool(boolValue)
				//TODO: more type
			}
			break
		}
	}

	return Properties, nil
}

/**
 * @description: pares configuration from stdin or file
 * @param {[]string} args
 * @return {*}
 */
func ParseConfig(args []string) (*ConfigProperties, error) {
	// default configuration
	if len(args) == 0 {
		logger.Log.Warn("* no config file specified, using the default config.")
		return Properties, nil
	}

	// specified config file
	if len(args) == 1 && args[0][0] != '-' {
		return parseConfigFile(args[0])
	}

	// parse arg
	i := 0
	for {
		if i >= len(args) {
			break
		}

		arg := args[i]

		if arg == "-v" || arg == "--version" {
			//version
			fmt.Println("redissgo " + Version)
			os.Exit(0)
		} else if arg == "-h" || arg == "--help" {
			//usage
			fmt.Println(usage)
			os.Exit(0)
		}

		if strings.HasPrefix(arg, "--") {
			//Todo: slice
			if i+1 >= len(args) {
				break
			}
			key, value := args[i][2:], args[i+1]

			parseConfigArg(key, value)
			i += 2
		} else {
			i++
		}
	}
	return Properties, nil
}
