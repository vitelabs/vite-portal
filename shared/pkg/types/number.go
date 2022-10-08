package types

import (
	"encoding/json"
	"errors"
)

type Int64 int64

func (u *Int64) UnmarshalJSON(bs []byte) error {
	var i int64
	if err := json.Unmarshal(bs, &i); err == nil {
		*u = Int64(i)
		return nil
	}
	var s string
	if err := json.Unmarshal(bs, &s); err != nil {
		return errors.New("expected a string or an integer")
	}
	if err := json.Unmarshal([]byte(s), &i); err != nil {
		return err
	}
	*u = Int64(i)
	return nil
}
