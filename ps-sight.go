package main

import (
	"strconv"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/keybase/go-ps"
)

// Take an event dictionary, return its full process tree
func eventProcessTree(event map[string]interface{}) []int {
	var processTree []int

	var eventPid, _ = strconv.Atoi(event["pid"].(string))
	var eventPPid, _ = strconv.Atoi(event["ppid"].(string))

	processTree = append(processTree, eventPid)

	if eventPPid != 1 {
		processTree = append(processTree, eventPPid)
		for eventPPid != 1 {
			parentProcess, err := ps.FindProcess(eventPPid)
			if err == nil && parentProcess != nil {
				if eventPPid != 1 {
					eventPPid = parentProcess.PPid()
					processTree = append(processTree, eventPPid)
				}
			} else {
				break
			}
		}
	}
	return processTree
}

// Take PID and return true if it's SSH
func isSSH(sockProcess netstat.Process) bool {
	if &sockProcess != nil {
		process, err := ps.FindProcess(sockProcess.Pid)
		if err == nil && process != nil {
			binaryPath, err := process.Path()
			if err == nil && binaryPath == "/usr/sbin/sshd" {
				return true
			}
		} else {
			if sockProcess.Name == "sshd" {
				return true
			}
		}
	}
	return false
}

// Take slice of sockets and return true if any one of them is a bind shell
func isBindShell(relevantSockets []netstat.SockTabEntry) bool {
	for _, socket := range relevantSockets {
		if socket.State == netstat.Listen && isSSH(*socket.Process) == false {
			return true
		}
	}
	return false
}

// Take slice of sockets and return true if any one of them is a reverse shell
func isReverseShell(relevantSockets []netstat.SockTabEntry) bool {
	for _, socket := range relevantSockets {
		if socket.State == netstat.Established && isSSH(*socket.Process) == false {
			return true
		}
	}
	return false
}
