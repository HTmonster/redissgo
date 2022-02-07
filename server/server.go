/*
 * @Description:
 * @Autor: HTmonster
 * @Date: 2022-02-06 14:33:21
 */

package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/HTmonster/redissgo/internal/logger"
)

/**
 * @description: setup a new server
 * @param {string} addr
 * @param {int} port
 * @return {*}
 */
func SetupAndListen(addr string, port int) error {

	/* ctrl-C ────────[signalChan]───────>ServerSetup
	                                       │
	ListenAndServe<────[closeChan]<────────┘*/

	// close signal
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan,
		syscall.SIGHUP,  //hong up
		syscall.SIGINT,  //interrupt ctrl-C
		syscall.SIGQUIT, //quit ctrl-
		syscall.SIGTERM, //terminate
	)

	// close channel
	closeChan := make(chan struct{})
	go func() {
		sig := <-signalChan //recive close signal
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
			closeChan <- struct{}{}
		}

	}()

	// server setup
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		logger.Log.Error("server setup error: ", err)
		return err
	}
	logger.Log.Info(fmt.Sprintf("server bind address: %s:%d ", addr, port))
	logger.Log.Info("start listening...")

	//server listen port and handle request
	ListenAndServe(listener, closeChan)

	return nil
}

/**
 * @description: Listen port and handle request
 * @param {net.Listener} listener
 * @param {<-chanstruct{}} closeChan
 * @return {*}
 */
func ListenAndServe(listener net.Listener, closeChan <-chan struct{}) {

	// creat a new request handler
	handler := NewHandler()

	// when a signal reached, close the server
	go func() {
		<-closeChan //recived close signal
		logger.Log.Info("server closing...")

		_ = listener.Close()
		_ = handler.Close()
	}()

	// close while exit
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()

	// get context
	ctx := context.Background()

	// wait group
	var waitDone sync.WaitGroup
	// wait the request coming
	for {
		// accept a new connection
		conn, err := listener.Accept()
		if err != nil {
			logger.Log.Error("accept a connection error: ", err)
			break
		}
		logger.Log.Info("accept a new connection: ", conn.RemoteAddr())

		waitDone.Add(1)

		// handle the connection
		go func() {
			defer func() {
				waitDone.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}
