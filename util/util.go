package util

import (
	"fmt"
	"strings"
	"syscall"
	"testing"
)

func GetFileName(path string) string {
	names := strings.Split(path, "/")
	return names[len(names)-1]
}

func GetTime() (int64, int64) {
	var r syscall.Rusage
	err := syscall.Getrusage(syscall.RUSAGE_SELF, &r)
	if err != nil {
		panic(err)
	}
	return r.Utime.Nano(), r.Stime.Nano()
}

// Assert asserts b is true, and fails otherwise
func Assert(t *testing.T, b bool, msg string) {
	if !b {
		t.Fatalf(msg)
	}
}

// AssertPanic asserts that the function totest panics
func AssertPanic(t *testing.T, totest func(), errmsg string) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(errmsg)
		}
	}()
	totest()
}

func ENTER() {
	fmt.Println("Press ENTER")
	fmt.Scanln()
}
