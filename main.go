package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/sys/unix"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: tinyrun run <cmd> [args...]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		parent()
	case "child":
		child()
	default:
		panic("usage: tinyrun run <cmd>")
	}
}

func parent() {
	self, err := os.Executable()
	must(err)

	cmd := exec.Command(
		self,
		append([]string{"child"},
			os.Args[2:]...,
		)...,
	)

	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	// syscall.SysProcAttr stays: os/exec types this field as *syscall.SysProcAttr,
	// so unix has no equivalent here. It's a struct type, not a deprecated call.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: unix.CLONE_NEWUTS | unix.CLONE_NEWPID | unix.CLONE_NEWNS,
	}

	must(cmd.Run())
}

func child() {

	must(unix.Sethostname([]byte("tinyrun")))
	must(unix.Chroot("./rootfs"))
	must(os.Chdir("/"))

	must(unix.Mount("", "/", "", unix.MS_REC|unix.MS_PRIVATE, ""))
	must(unix.Mount("proc", "/proc", "proc", 0, ""))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	must(cmd.Run())

	must(unix.Unmount("/proc", 0))
}

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "tinyrun: %v\n", err)
		os.Exit(1)
	}
}
