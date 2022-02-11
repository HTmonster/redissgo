/*
 * @Description: RESP parser
 * @Autor: HTmonster
 * @Date: 2022-02-11 14:23:39
 */

package request

import (
	"bufio"
	"errors"
	"io"
	"strconv"

	. "github.com/HTmonster/redissgo/internal/logger"
)

/**
 * @description: given a request, parse it
 * @param {io.Reader} reader
 * @return {*}
 */
func ParseRequest(reader io.Reader) <-chan *Request {
	ch := make(chan *Request)
	go parse(reader, ch)
	return ch
}

/**
 * @description: parse request, using channel to send result
 * @param {io.Reader} reader
 * @param {<-chan*Request} ch
 * @return {*}
 */
func parse(reader io.Reader, ch chan<- *Request) {
	defer func() {
		if err := recover(); err != nil {
			Log.Error(err)
		}
	}()

	bufreader := bufio.NewReader(reader)

	var params [][]byte
	var err error
	var paramLen, paramBytes int64
	var paramCnt int64

	status := Begin
	for {
		switch status {
		case Begin:
			if err = parseBegin(bufreader); err != nil {
				goto End
			}
			status = ParamLen
		case ParamLen:
			paramLen, err = parseParamLen(bufreader)
			if err != nil {
				goto End
			}
			paramCnt = 0
			status = ParamBytes
		case ParamBytes:
			paramBytes, err = parseParamBytes(bufreader)
			if err != nil {
				goto End
			}
			status = ParamData
		case ParamData:
			var msg []byte
			msg, err = parseParamData(bufreader, paramBytes)
			if err != nil {
				goto End
			}
			params = append(params, msg)
			paramCnt++
			if paramCnt >= paramLen {
				goto End
			} else {
				status = ParamBytes
			}
		default:
			goto End
		}

	}

End:
	ch <- &Request{
		Params: params,
		Len:    paramLen,
		Err:    err,
	}
	close(ch)
}

/**
 * @description: parse request begin
 * @param {*bufio.Reader} bufreader
 * @return {*}
 */
func parseBegin(bufreader *bufio.Reader) error {
	b, err := bufreader.ReadByte()

	if err != nil {
		Log.Error(err)
		return err
	}

	if b != '*' {
		err := errors.New("Unkonw request begin symbol :" + string(b))
		Log.Error(err)
		return err
	}

	return nil
}

/**
 * @description: parse request length
 * @param {*bufio.Reader} bufreader
 * @return {*}
 */
func parseParamLen(bufreader *bufio.Reader) (int64, error) {
	msg, err := bufreader.ReadBytes('\n')
	if err != nil {
		Log.Error(err)
		return 0, err
	}
	if len(msg) == 0 || msg[len(msg)-2] != '\r' {
		err = errors.New("Error protocol, parameter length error:" + string(msg))
		Log.Error(err)
		return 0, err
	}

	msg = msg[:len(msg)-2]

	paramLen, err := strconv.ParseInt(string(msg), 10, 64)
	if err != nil {
		err = errors.New("Error parameter length error:" + string(msg))
		Log.Error(err)
		return 0, err
	}

	return paramLen, nil
}

func parseParamBytes(bufreader *bufio.Reader) (int64, error) {
	b, err := bufreader.ReadByte()

	if err != nil {
		Log.Error(err)
		return 0, err
	}

	if b != '$' {
		err := errors.New("Unkonw request paramater size symbol :" + string(b))
		Log.Error(err)
		return 0, err
	}

	msg, err := bufreader.ReadBytes('\n')
	if err != nil {
		Log.Error(err)
		return 0, err
	}
	if len(msg) == 0 || msg[len(msg)-2] != '\r' {
		err = errors.New("Error protocol, parameter size error:" + string(msg))
		Log.Error(err)
		return 0, err
	}

	msg = msg[:len(msg)-2]

	paramSize, err := strconv.ParseInt(string(msg), 10, 64)
	if err != nil {
		err = errors.New("Error parameter size:" + string(msg))
		Log.Error(err)
		return 0, err
	}

	return paramSize, nil
}

/**
 * @description: read request parameter content (binary safe)
 * @param {bufio.Reader} bufreader
 * @return {*}
 */
func parseParamData(bufreader *bufio.Reader, bytes int64) ([]byte, error) {

	msg := make([]byte, bytes+2)
	_, err := io.ReadFull(bufreader, msg)
	if err != nil {
		Log.Error(err)
		return nil, err
	}
	if len(msg) == 0 || msg[len(msg)-1] != '\n' || msg[len(msg)-2] != '\r' {
		err = errors.New("Error parameter content:" + string(msg))
		Log.Error(err)
		return nil, err
	}
	msg = msg[0 : len(msg)-2]
	return msg, nil
}
