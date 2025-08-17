package easyjson

import (
	"errors"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

// rsp
// @Description:
type rsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
}

// ResponseJson
//
//	@Description:
//	@param w
//	@param resp
//	@param err
func ResponseJson(w http.ResponseWriter, resp interface{}, err error) {
	var body rsp
	if err != nil {
		body.Code = -1
		body.Msg = err.Error()
	} else {
		body.Data = resp
	}

	//do write json
	err = doWriteJson(w, http.StatusOK, &body)
	if err != nil {
		logx.Errorf("doWriteJson err:%s", err)
		return
	}
}

// doWriteJson
//
//	@Description:
//	@param w
//	@param statusCode
//	@param v
//	@return error
func doWriteJson(w http.ResponseWriter, statusCode int, v any) error {

	bs, err := easyjson.Marshal(v.(easyjson.Marshaler))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("marshal json failed, error: %w", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	//write data
	if n, err := w.Write(bs); err != nil {
		// http.ErrHandlerTimeout has been handled by http.TimeoutHandler,
		// so it's ignored here.
		if !errors.Is(err, http.ErrHandlerTimeout) {
			return fmt.Errorf("write response failed, error: %w", err)
		}
	} else if n < len(bs) {
		return fmt.Errorf("actual bytes: %d, written bytes: %d", len(bs), n)
	}
	//all bytes data write to client DONE !
	return nil
}
