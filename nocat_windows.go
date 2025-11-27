//go:build windows

package main

import (
	"golang.org/x/sys/windows"
)

var newConsoleAttr = &windows.SysProcAttr{
	CreationFlags: windows.CREATE_NEW_CONSOLE,
}
