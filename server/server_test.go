/*
 * @Description:
 * @Autor: HTmonster
 * @Date: 2022-02-07 10:40:50
 */
package server

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestServer(t *testing.T) {

	msg := "hello world"

	// setup server
	go func() {
		if err := SetupAndListen("127.0.0.1", 6379); err != nil {
			t.Error("error while setup server: ", err)
			return
		}

		defer func() {
			p, _ := os.FindProcess(os.Getpid())
			_ = p.Signal(syscall.SIGINT)
		}()

	}()

	time.Sleep(1 * time.Second)
	// client
	conn, err := net.Dial("tcp", ":6379")
	if err != nil {
		t.Error("error connecting to server: ", err)
	}
	defer conn.Close()

	if _, err := fmt.Fprint(conn, msg); err != nil {
		t.Fatal(err)
	}
}
