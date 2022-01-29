/*
 * @Description: logger
 * @Autor: HTmonster
 * @Date: 2022-01-29 10:11:21
 */
package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

// log initial
func init() {
	// logger level
	Log.SetLevel(logrus.DebugLevel)
	// format
	formater := new(logrus.TextFormatter)
	formater.FullTimestamp = true
	formater.TimestampFormat = "2006-01-01 00:00:00"
	formater.DisableTimestamp = false
	formater.DisableColors = false
	Log.SetFormatter(formater)
}

/**
 * @description: set log output file
 * @param {string} file
 * @return err
 */
func SetLogFile(file string) error {
	logfile, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		Log.Fatal(err)
		return err
	}
	logrus.SetOutput(logfile)
	return nil
}
