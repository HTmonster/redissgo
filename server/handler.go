/*
 * @Description:
 * @Autor: HTmonster
 * @Date: 2022-02-07 09:26:14
 */

package server

import (
	"context"
	"net"

	"github.com/HTmonster/redissgo/internal/logger"
)

//------------ handler --------------
type Handler struct {
	content string // stub
}

/**
 * @description: creat a new Handler instance
 */
func NewHandler() *Handler {
	return &Handler{
		content: "",
	}
}

/**
 * @description: handle a new connection
 * @param {context.Context} ctx
 * @param {net.Conn} conn
 * @return {*}
 */
func (*Handler) Handle(ctx context.Context, conn net.Conn) error {
	defer conn.Close()
	for {
		var buf [128]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			logger.Log.Error("Read from tcp server failed,err:", err)
			break
		}
		data := string(buf[:n])
		logger.Log.Info("Read from tcp server:", data)
	}

	return nil
}

/**
 * @description: close a handler
 * @event:
 * @param {*}
 * @return {*}
 */
func (*Handler) Close() error {
	logger.Log.Info("handler closed.")
	return nil
}
