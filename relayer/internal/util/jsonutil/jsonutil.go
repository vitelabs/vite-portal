package jsonutil

import (
	"encoding/json"

	"github.com/vitelabs/vite-portal/relayer/internal/logger"
	"github.com/vitelabs/vite-portal/relayer/internal/util/runtimeutil"
)

func ToByte(v any) ([]byte, error) {
	output, err := json.Marshal(v)
	return output, err
}

func ToByteOrExit(v any) []byte {
	output, err := ToByte(v)
	if err != nil {
		exit(err, "ToByteOrExit failed")
	}
	return output
}

func ToByteIndent(v any) ([]byte, error) {
	output, err := json.MarshalIndent(v, "", "    ")
	return output, err
}

func ToByteIndentOrExit(v any) []byte {
	output, err := ToByteIndent(v)
	if err != nil {
		exit(err, "ToByteIndentOrExit failed")
	}
	return output
}

func ToString(v any) string {
	output, err := ToByte(v)
	if err != nil {
		return ""
	}
	return string(output)
}

func FromByte(data []byte, v any) error {
	err := json.Unmarshal(data, &v)
	return err
}

func FromByteOrExit(data []byte, v any)  {
	err := FromByte(data, v)
	if err != nil {
		exit(err, "FromByteOrExit failed")
	}
}

func exit(err error, msg string) {
	logger.Logger().Fatal().Err(err).Str("stacktrace", string(ToByteOrExit(runtimeutil.CaptureStacktrace(3, 32)))).Msg(msg)
}