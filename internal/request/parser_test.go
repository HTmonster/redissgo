/*
 * @Description:
 * @Autor: HTmonster
 * @Date: 2022-02-11 16:06:34
 */
package request

import (
	"bytes"
	"fmt"
	"testing"
)

func TestParseRequest(t *testing.T) {
	request := []byte("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")

	ch := ParseRequest(bytes.NewReader(request))
	if ch == nil {
		t.Errorf("error")
	}
	fmt.Print(<-ch)

	request = []byte("3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")

	ch = ParseRequest(bytes.NewReader(request))
	if ch == nil {
		t.Errorf("error")
	}
	fmt.Print(<-ch)

	request = []byte("*3$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")

	ch = ParseRequest(bytes.NewReader(request))
	if ch == nil {
		t.Errorf("error")
	}
	fmt.Print(<-ch)

	request = []byte("*10\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")

	ch = ParseRequest(bytes.NewReader(request))
	if ch == nil {
		t.Errorf("error")
	}
	fmt.Print(<-ch)
}
