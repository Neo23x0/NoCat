//go:build !windows

package main

import "syscall"

var newConsoleAttr *syscall.SysProcAttr = nil
