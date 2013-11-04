// +build darwin freebsd netbsd openbsd

package goline

import "syscall"

const ioctlReadTermios = syscall.TIOCGETA
const ioctlWriteTermios = syscall.TIOCSETA
