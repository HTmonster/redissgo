/*
 * @Description:
 * @Autor: HTmonster
 * @Date: 2022-01-29 11:08:28
 */
package config

import (
	"path"
	"runtime"
	"testing"
)

func TestParseConfigFile(t *testing.T) {
	_, fullFileName, _, _ := runtime.Caller(0)
	confFile := path.Clean(path.Join(fullFileName, "../../../redis.conf"))

	if _, err := parseConfigFile(confFile); err != nil {
		t.Errorf("error parsing config file: %s", err)
	}
}

func TestParseConfigArg(t *testing.T) {
	if _, err := parseConfigArg("port", "9999"); err != nil {
		t.Errorf("error parsing config arg: %s", err)
	}
	if _, err := parseConfigArg("port", "ssss"); err != nil {
		t.Errorf("error parsing config arg: %s", err)
	}
	if _, err := parseConfigArg("none", "9999"); err != nil {
		t.Errorf("error parsing config arg: %s", err)
	}
}
