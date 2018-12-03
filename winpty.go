package winpty

import (
	"fmt"
	"os"
)

type Options struct {
	// DLLPrefix is the path to winpty.dll and winpty-agent.exe
	DLLPrefix string

	// AppName sets the title of the console
	AppName string

	// Command is the full command to launch
	Command string

	// Dir sets the current working directory for the command
	Dir string

	// Env sets the environment variables. Use the format VAR=VAL.
	Env []string
}

type WinPTY struct {
	StdIn  *os.File
	StdOut *os.File

	wp          uintptr
	childHandle uintptr
	closed      bool
}

// accepts path to command to execute, then arguments.
// returns WinPTY object pointer, error.
// remember to call Close on WinPTY object when done.
func Open(dllPrefix, cmd string) (*WinPTY, error) {
	return OpenWithOptions(Options{
		DLLPrefix: dllPrefix,
		Command:   cmd,
	})
}

// the same as open, but uses defaults for Env & Dir
func OpenDefault(dllPrefix, cmd string) (*WinPTY, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Failed to get dir on setup: %s", err)
	}

	return OpenWithOptions(Options{
		DLLPrefix: dllPrefix,
		Command:   cmd,
		Dir:       wd,
		Env:       os.Environ(),
	})
}

func OpenWithOptions(options Options) (*WinPTY, error) {
	return openWithOptions(options)
}

func (obj *WinPTY) SetSize(ws_col, ws_row uint32) {
	setSize(obj, ws_col, ws_row)
}

func (obj *WinPTY) Close() {
	close(obj)
}
