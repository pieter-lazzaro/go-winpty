package winpty

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func openWithOptions(options Options) (*WinPTY, error) {
	setupDefines(options.DLLPrefix)

	agentCfg, err := createAgentCfg(0)

	if err != nil {
		return nil, err
	}

	// Set the initial size to 40x40.
	winpty_config_set_initial_size.Call(agentCfg, uintptr(uint32(40)), uintptr(uint32(40)))

	var openErr uintptr
	defer winpty_error_free.Call(openErr)
	wp, _, _ := winpty_open.Call(agentCfg, uintptr(unsafe.Pointer(openErr)))

	if wp == uintptr(0) {
		return nil, fmt.Errorf("Error Launching WinPTY agent, %s", GetErrorMessage(openErr))
	}

	winpty_config_free.Call(agentCfg)

	stdin_name, _, _ := winpty_conin_name.Call(wp)
	stdout_name, _, _ := winpty_conout_name.Call(wp)

	obj := &WinPTY{}
	stdin_handle, err := syscall.CreateFile((*uint16)(unsafe.Pointer(stdin_name)), syscall.GENERIC_WRITE, 0, nil, syscall.OPEN_EXISTING, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("Error getting stdin handle. %s", err)
	}
	obj.StdIn = os.NewFile(uintptr(stdin_handle), "stdin")
	stdout_handle, err := syscall.CreateFile((*uint16)(unsafe.Pointer(stdout_name)), syscall.GENERIC_READ, 0, nil, syscall.OPEN_EXISTING, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("Error getting stdout handle. %s", err)
	}
	obj.StdOut = os.NewFile(uintptr(stdout_handle), "stdout")

	spawnCfg, err := createSpawnCfg(WINPTY_SPAWN_FLAG_AUTO_SHUTDOWN, options.AppName, options.Command, options.Dir, options.Env)

	if err != nil {
		return nil, err
	}

	var (
		spawnErr  uintptr
		lastError *uint32
	)
	spawnRet, _, _ := winpty_spawn.Call(wp, spawnCfg, uintptr(unsafe.Pointer(obj.childHandle)), uintptr(0), uintptr(unsafe.Pointer(lastError)), uintptr(unsafe.Pointer(spawnErr)))
	winpty_spawn_config_free.Call(spawnCfg)
	defer winpty_error_free.Call(spawnErr)

	if spawnRet == 0 {
		return nil, fmt.Errorf("Error spawning process... %s", GetErrorMessage(spawnErr))
	} else {
		obj.wp = wp
		return obj, nil
	}
}

func setSize(obj *WinPTY, ws_col, ws_row uint32) {
	if ws_col == 0 || ws_row == 0 {
		return
	}
	winpty_set_size.Call(obj.wp, uintptr(ws_col), uintptr(ws_row), uintptr(0))
}

func close(obj *WinPTY) {
	if obj.closed {
		return
	}

	winpty_free.Call(obj.wp)

	obj.StdIn.Close()
	obj.StdOut.Close()

	syscall.CloseHandle(syscall.Handle(obj.childHandle))

	obj.closed = true
}
