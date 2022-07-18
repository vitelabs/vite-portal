package types

import (
	"errors"
	"fmt"

	roottypes "github.com/vitelabs/vite-portal/internal/types"
)

const DefaultCodeNamespace roottypes.CodeNamespace = ModuleName

const (
	CodeInvalidInputError = iota + 1
	codeTypeLimit         // This will be the last CodeType + 1
)

var CodeTypeErrors = map[roottypes.CodeType]error{
	CodeInvalidInputError: errors.New("invalid input: "),
}

func GetCodeError(code roottypes.CodeType) error {
	err := CodeTypeErrors[code]
	if err == nil {
		return errors.New(fmt.Sprintf("unknown code %d", code))
	}
	return err
}

func NewBasicError(ns roottypes.CodeNamespace, code roottypes.CodeType) roottypes.Error {
	return roottypes.NewError(ns, code, GetCodeError(code).Error())
}

func NewError(ns roottypes.CodeNamespace, code roottypes.CodeType, err error) roottypes.Error {
	return roottypes.NewError(ns, code, GetCodeError(code).Error()+err.Error())
}
