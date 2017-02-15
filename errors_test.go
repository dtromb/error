package errors

import (
	"testing"
)


func TestError(t *testing.T) {
	erf := func () error {
		return New("Error.  Error.  Danger Will Robinson!  Danger!")
	}
	PrintStackTrace(erf())
}

func TestPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			PrintStackTrace(PanicToError(r))
		}
	}()
	erf := func () error {
		panic("NOT YOURS")
		return New("can I have a pony?")
	}
	PrintStackTrace(erf())
}

func TestCause(t *testing.T) {
	erf1  := func() error {
		return New("this function is never gonna work")
	} 
	erf2 := func() error {
		err := erf1()
		if err != nil {
			return New("is this function gonna work?").CausedBy(err)
		}
		return nil
	}
	PrintStackTrace(erf2())
}