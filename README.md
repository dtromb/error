# errors

type-extended errors for Go with stack traces

```go  
func TestError(t *testing.T) {
	erf := func () error {
		return New("Error.  Error.  Danger Will Robinson!  Danger!")
	}
	PrintStackTrace(erf())
}
```

[0] github.com/dtromb/errors.TestError.func1@0x0000000000476B26 at /home/dtrombley/go/path/src/github.com/dtromb/errors/errors_test.go:10
[1] github.com/dtromb/errors.TestError@0x0000000000476A2D at /home/dtrombley/go/path/src/github.com/dtromb/errors/errors_test.go:12
[2] testing.tRunner@0x0000000000471161 at /home/dtrombley/go/root/src/testing/testing.go:611
[3] runtime.goexit@0x000000000045D5A1 at /home/dtrombley/go/root/src/runtime/asm_amd64.s:2087

```go 
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
```


[0] github.com/dtromb/errors.PanicToError@0x0000000000475DDE at /home/dtrombley/go/path/src/github.com/dtromb/errors/errors.go:145
[1] github.com/dtromb/errors.TestPanic.func1@0x0000000000476BC2 at /home/dtrombley/go/path/src/github.com/dtromb/errors/errors_test.go:19
[2] runtime.call32@0x000000000045AA5C at /home/dtrombley/go/root/src/runtime/asm_amd64.s:479
[3] runtime.gopanic@0x000000000042B473 at /home/dtrombley/go/root/src/runtime/panic.go:459
[4] github.com/dtromb/errors.TestPanic.func2@0x0000000000476C5D at /home/dtrombley/go/path/src/github.com/dtromb/errors/errors_test.go:23
[5] github.com/dtromb/errors.TestPanic@0x0000000000476AA9 at /home/dtrombley/go/path/src/github.com/dtromb/errors/errors_test.go:26
[6] testing.tRunner@0x0000000000471161 at /home/dtrombley/go/root/src/testing/testing.go:611
[7] runtime.goexit@0x000000000045D5A1 at /home/dtrombley/go/root/src/runtime/asm_amd64.s:2087
