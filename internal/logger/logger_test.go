/*
 * @Description:
 * @Autor: HTmonster
 * @Date: 2022-01-29 10:48:41
 */
package logger

import "testing"

func TestSetLogFile(t *testing.T) {
	if err := SetLogFile("/tmp/test.log"); err != nil {
		t.Errorf("set log file failed. %s", err)
	}
}
