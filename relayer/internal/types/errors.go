package types

import (
	"fmt"

	"github.com/vitelabs/vite-portal/relayer/internal/util/runtimeutil"
)

// Error code namespace identifier
type CodeNamespace string

// Error code identifier within a code namespace
type CodeType uint

// Root error code namespaces
const (
	// Unset namespace can be overriden with Error.WithDefaultNamespace()
	CodeNsUndefined CodeNamespace = ""
	// For error codes in this file only
	CodeNsRoot CodeNamespace = "root"
)

// Core error codes
const (
	CodeOK = iota
	CodeInternal
	codeTypeLimit // This will be the last CodeType + 1
)

// Core code messages
var CodeTypeMessages = map[CodeType]string{
	CodeOK:       "everything ok",
	CodeInternal: "internal error",
}

func GetCodeMessage(code CodeType) string {
	message := CodeTypeMessages[code]
	if message == "" {
		return fmt.Sprintf("unknown code %d", code)
	}
	return message
}

type Error interface {
	CodeNamespace() CodeNamespace
	Code() CodeType
	Error() string
	ErrorFormatted() string
	InnerError() string
	Data() interface{}
	WithDefaultCodeNamespace(CodeNamespace) Error
}

func NewError(ns CodeNamespace, code CodeType, format string, args ...interface{}) Error {
	return newError(ns, code, format, args...)
}

func newErrorWithRootCodeNamespace(code CodeType, format string, args ...interface{}) *rootError {
	return newError(CodeNsRoot, code, format, args...)
}

func newError(ns CodeNamespace, code CodeType, format string, args ...interface{}) *rootError {
	if format == "" {
		format = GetCodeMessage(code)
	}
	err := &rootError{
		ns:   ns,
		code: code,
		data: FmtError{
			format: format,
			args:   args,
		},
	}
	err.doStacktrace()
	return err
}

type rootError struct {
	ns         CodeNamespace
	code       CodeType
	data       FmtError
	stacktrace []runtimeutil.StacktraceItem
}

// ---
// Implement "Error" interface

func (err *rootError) CodeNamespace() CodeNamespace {
	return err.ns
}

func (err *rootError) Code() CodeType {
	return err.code
}

func (err *rootError) Error() string {
	return fmt.Sprintf("%#v", err)
}

func (err *rootError) ErrorFormatted() string {
	return fmt.Sprintf(`ERROR:
Namespace: %s
Code: %d
Message: %#v
`, err.ns, err.code, err)
}

func (err *rootError) InnerError() string {
	return err.data.Error()
}

func (err *rootError) Data() interface{} {
	return err.data
}

func (err *rootError) WithDefaultCodeNamespace(ns CodeNamespace) Error {
	namespace := err.ns
	if namespace == CodeNsUndefined {
		namespace = ns
	}
	return &rootError{
		ns:   namespace,
		code: err.code,
	}
}

// ---
// FmtError

type FmtError struct {
	format string
	args   []interface{}
}

func (fe FmtError) Error() string {
	return fmt.Sprintf(fe.format, fe.args...)
}

func (fe FmtError) String() string {
	return fmt.Sprintf("FmtError{format:%v,args:%v}", fe.format, fe.args)
}

func (fe FmtError) Format() string {
	return fe.format
}

// ---
// Stacktrace

// Captures stacktrace if not present already
func (err *rootError) doStacktrace() {
	if err.stacktrace == nil {
		var offset = 3
		var depth = 32
		err.stacktrace = runtimeutil.CaptureStacktrace(offset, depth)
	}
}
