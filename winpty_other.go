// +build !windows

package winpty

import (
	"errors"
)

func openWithOptions(options Options) (*WinPTY, error) {
	return nil, errors.New("operating system not supported")
}

func setSize(obj *WinPTY, ws_col, ws_row uint32) {

}

func close(obj *WinPTY) {

}
