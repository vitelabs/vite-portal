package httputil

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/jsonutil"
)

const (
	ContentTypeTextPlain = "text/plain; charset=UTF-8"
	ContentTypeJson      = "application/json; charset=UTF-8"
)

func SetFallbackClientIp(h http.Header, value string) {
	host, _, err := net.SplitHostPort(value)
	if err != nil || host == "" {
		logger.Logger().Error().Err(err).Msg(fmt.Sprintf("couldn't split host and port of '%s'", value))
		return
	}
	h.Set(types.HeaderFallbackClientIp, host)
}

func GetFallbackClientIp(h http.Header) string {
	return h.Get(types.HeaderFallbackClientIp)
}

func WriteResponse(w http.ResponseWriter, data, contentType string) {
	WriteResponseWithCode(w, data, contentType, http.StatusOK)
}

func WriteResponseWithCode(w http.ResponseWriter, data, contentType string, code int) {
	writeHeader(w, contentType)
	w.WriteHeader(code)
	_, err2 := w.Write([]byte(data))
	if err2 != nil {
		logger.Logger().Error().Err(err2).Msg("WriteResponseWithCode failed")
	}
}

func WriteJsonResponse(w http.ResponseWriter, data any) {
	WriteJsonResponseWithCode(w, data, http.StatusOK)
}

func WriteJsonResponseWithCode(w http.ResponseWriter, data any, code int) {
	b, err1 := jsonutil.ToByte(data)
	if err1 != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err1.Error())
		logger.Logger().Error().Err(err1).Msg("WriteJsonResponseWithCode failed")
		return
	}
	writeHeader(w, ContentTypeJson)
	w.WriteHeader(code)
	_, err2 := w.Write(b)
	if err2 != nil {
		logger.Logger().Error().Err(err2).Msg("WriteJsonResponseWithCode failed")
	}
}

func WriteErrorResponse(w http.ResponseWriter, code int, msg string) {
	err := &types.RpcError{
		Code:    code,
		Message: msg,
	}
	WriteJsonResponseWithCode(w, err, code)
}

func writeHeader(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Type", contentType)
}

func ExtractBody(r *http.Request, maxRequestContentLength int64) ([]byte, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxRequestContentLength))
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, nil
	}
	if err := r.Body.Close(); err != nil {
		return nil, err
	}
	if logger.DebugEnabled() {
		logger.Logger().Debug().Str("body", string(body)).Msg("request body")
	}
	return body, nil
}

func ExtractModel(w http.ResponseWriter, r *http.Request, model interface{}, maxRequestContentLength int64) error {
	body, err := ExtractBody(r, maxRequestContentLength)
	if err != nil {
		return err
	}
	return ExtractModelFromBody(body, &model)
}

func ExtractModelFromBody(body []byte, model interface{}) error {
	if len(body) == 0 {
		return nil
	}
	if err := jsonutil.FromByte(body, model); err != nil {
		return err
	}
	return nil
}

func ExtractQuery(w http.ResponseWriter, r *http.Request, model interface{}) error {
	q := r.URL.Query()
	if len(q) == 0 {
		return nil
	}

	m := map[string]interface{}{}
	for k, v := range q {
		m[k] = v[0]
	}

	b, err := jsonutil.ToByte(m)
	if err != nil {
		return err
	}
	return ExtractModelFromBody(b, &model)
}
