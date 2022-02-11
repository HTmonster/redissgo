/*
 * @Description: RESP request
 * @Autor: HTmonster
 * @Date: 2022-02-11 14:39:45
 */

package request

// parse status code
const (
	Begin int = iota
	ParamLen
	ParamBytes
	ParamData
	End
)

type Request struct {
	Params [][]byte
	Len    int64
	Err    error
}
