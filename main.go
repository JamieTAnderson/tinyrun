package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/sys/unix"
)

func main() {
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
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: unix.CLONE_NEWUTS | unix.CLONE_NEWPID | unix.CLONE_NEWNS,
	}

	must(cmd.Run())
}

func child() {

	must(syscall.Sethostname([]byte("tinyrun")))
	must(syscall.Chroot("./rootfs"))
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "/proc", "proc", 0, ""))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	must(cmd.Run())

	must(syscall.Unmount("/proc", 0))
}

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "tinyrun: %v\n", err)
		os.Exit(1)
	}
}
