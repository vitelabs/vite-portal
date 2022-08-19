package httputil

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/util/jsonutil"
)

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

func ExtractQuery(w http.ResponseWriter, r *http.Request, _ httprouter.Params, model interface{}) error {
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
