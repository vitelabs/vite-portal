package types

import (
	"errors"
	"fmt"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
)

type Relayer struct {
	Id      string `json:"id"`
	Version string `json:"version"`
}

func (r *Relayer) IsValid() bool {
	return r != nil && r.Id != ""
}

func (r *Relayer) Validate() error {
	if !r.IsValid() {
		msg := "relayer is invalid"
		err := errors.New(msg)
		logger.Logger().Error().Err(err).Str("relayer", fmt.Sprintf("%#v", r)).Msg(msg)
		return err
	}
	return nil
}
