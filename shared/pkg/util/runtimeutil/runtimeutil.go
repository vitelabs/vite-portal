package runtimeutil

import (
	"fmt"
	"runtime"
	"strings"
)

type FuncInfos struct {
	Line      int
	Filename  string
	Funcname  string
	Formatted string
}

func ToFuncInfos(pc uintptr) FuncInfos {
	f := runtime.FuncForPC(pc)
	file, line := f.FileLine(pc)
	filename := file[strings.LastIndex(file, "/")+1:]
	funcname := f.Name()
	fn := funcname[strings.LastIndex(funcname, ".")+1:]
	return FuncInfos{
		Line:      line,
		Filename:  filename,
		Funcname:  fn,
		Formatted: fmt.Sprintf("%s:%d#%s", filename, line, fn),
	}
}

type StacktraceItem struct {
	Pointer uintptr
	Infos   string
}

func CaptureStacktrace(offset int, depth int) []StacktraceItem {
	var pcs = make([]uintptr, depth)
	n := runtime.Callers(offset, pcs)
	var items = make([]StacktraceItem, n)
	for i, s := range pcs[0:n] {
		items[i] = StacktraceItem{
			Pointer: s,
			Infos:   ToFuncInfos(s).Formatted,
		}
	}
	return items
}
