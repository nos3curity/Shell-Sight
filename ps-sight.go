package main

import (
	"strconv"

	"github.com/mitchellh/go-ps"
)

func eventProcessTree(event map[string]interface{}) []int {
	var processTree []int

	var eventPid, _ = strconv.Atoi(event["pid"].(string))
	var eventPPid, _ = strconv.Atoi(event["ppid"].(string))

	processTree = append(processTree, eventPid)

	if eventPPid != 1 {
		processTree = append(processTree, eventPPid)
		for eventPPid != 1 {
			parentProcess, err := ps.FindProcess(eventPPid)
			eventPPid = parentProcess.PPid()
			if eventPPid != 1 && err == nil {
				processTree = append(processTree, eventPPid)
			}
		}
	}
	return processTree
}
