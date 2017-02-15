package errors

import (
	"fmt"
	"errors"
	"net/http" 
	"reflect"
	"runtime"
)

type StackTraceEntry struct {
	pc uintptr
	file string
	line int
	f *runtime.Func
}

type StackTracedError interface {
	error
	Cause() error
	CausedBy(cause error) StackTracedError
	Trace() []StackTraceEntry
	Complete() bool
}

type StdStackTracedError struct {
	error
	trace []StackTraceEntry
	cause error
	ellipsis bool
}

type RuntimeStackTracedError struct {
	StdStackTracedError
}

func (ste *StdStackTracedError) Error() string {
	return ste.error.Error()
}

func (ste *RuntimeStackTracedError) RuntimeError() {}


func (ste *StdStackTracedError) Trace() []StackTraceEntry {
	return ste.trace
}

func (ste *StdStackTracedError) Cause() error {
	return ste.cause
}

func (ste *StdStackTracedError) CausedBy(cause error) StackTracedError {
	if ste.cause != nil {
		panic("CausedBy() called twice")
	}
	ste.cause = cause
	return ste
}

func (ste *StdStackTracedError) Complete() bool {
	return !ste.ellipsis
}

var global_DEFAULT_STACK_TRACE_MAX_DEPTH int = 20

func NewStackTracedError(err error) StackTracedError {
	var ste StackTracedError
	var sste *StdStackTracedError
	if st, isSTE := err.(*StdStackTracedError); isSTE {
		return st
	}
	if _, isRTE := err.(runtime.Error); isRTE {
		rste := &RuntimeStackTracedError{
			StdStackTracedError: StdStackTracedError{error: err},
		}
		ste = rste
		sste = &rste.StdStackTracedError
	} else {
		sste = &StdStackTracedError{
			error: err,
		}
		ste = sste
	}
	sste.trace = make([]StackTraceEntry, 0, global_DEFAULT_STACK_TRACE_MAX_DEPTH)
	for i := 1; i < global_DEFAULT_STACK_TRACE_MAX_DEPTH+2; i++ {
		pc, file, line, ok := runtime.Caller(1+i)
		if !ok {
			break
		}
		if i == global_DEFAULT_STACK_TRACE_MAX_DEPTH+1 {
			sste.ellipsis = true
		} else {
			sste.trace = append(sste.trace, StackTraceEntry{
				pc: pc,
				file: file,
				line: line,
			})
		}
	}
	return ste
}

// Duplicated here to avoid circular package refs
type Stringable interface {
	String() string
}

type WebserverError interface {
	error
	StatusCode() int
}

type stdWebserverError struct {
	status int
	msg string
}

func (we *stdWebserverError) Error() string {
	return we.msg
}

func (we *stdWebserverError) StatusCode() int {
	return we.status
}

func ErrBadRequest(msg string) WebserverError {
	return &stdWebserverError{
		status: http.StatusBadRequest,
		msg: msg,
	}
}

func ErrNotFound(msg string) WebserverError {
	return &stdWebserverError{
		status: http.StatusNotFound,
		msg: msg,
	}
}

func PanicToError(panicRecover interface{}) error {
	if panicRecover == nil {
		return nil
	}
	if str, isStr := panicRecover.(string); isStr {
		return New(str)
	} else if str, isStr := panicRecover.(Stringable); isStr {
		return New(str.String())
	} else if e, isErr := panicRecover.(error); isErr {
		return NewStackTracedError(e)
	} 
	return New(reflect.TypeOf(panicRecover).String())
}

func New(msg string) StackTracedError {
	return NewStackTracedError(errors.New(msg))
}

func (ste *StackTraceEntry) fillIn() {
	ste.f = runtime.FuncForPC(ste.pc)
}

func StackTrace(err error) string {
	var buf []byte
	e := err
	if e != nil {
		buf = append(buf, (e.Error()+"\n")...)
	}
	for e != nil {
		if ste, ok := e.(StackTracedError); ok {
			for i, entry := range ste.Trace() {
				if entry.f == nil {
					entry.fillIn()
				}
				if entry.f != nil {
					file, line := entry.f.FileLine(entry.pc)
					buf = append(buf, fmt.Sprintf("[%d] %s@0x%16.16X at %s:%d\n", i, entry.f.Name(), uint64(entry.pc), file, line)...)
				} else {
					buf = append(buf, fmt.Sprintf("[%d] ???@%16.16X at %s:%d\n", i, uint64(entry.pc), entry.file, entry.line)...)
				}
			}
			if !ste.Complete() {
				buf = append(buf,"...\n"...)
			}
			if ste.Cause() != nil {
				buf = append(buf,fmt.Sprintf("... caused by: %s\n",ste.Cause().Error())...)
			}
			e = ste.Cause()
		} else {
			buf = append(buf,"<no stack trace available>\n"...)
			e = nil
		}
	}
	return string(buf)
}

func PrintStackTrace(err error) {
	e := err
	for e != nil {
		if ste, ok := e.(StackTracedError); ok {
			for i, entry := range ste.Trace() {
				if entry.f == nil {
					entry.fillIn()
				}
				if entry.f != nil {
					file, line := entry.f.FileLine(entry.pc)
					fmt.Printf("[%d] %s@0x%16.16X at %s:%d\n", i, entry.f.Name(), uint64(entry.pc), file, line)
				} else {
					fmt.Printf("[%d] ???@%16.16X at %s:%d\n", i, uint64(entry.pc), entry.file, entry.line)
				}
			}
			if !ste.Complete() {
				fmt.Println("...")
			}
			if ste.Cause() != nil {
				fmt.Println("... caused by: "+ste.Cause().Error())
			}
			e = ste.Cause()
		} else {
			fmt.Println("<no stack trace available>")
			e = nil
		}
	}
}