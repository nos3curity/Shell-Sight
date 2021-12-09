package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/elastic/go-libaudit/auparse"
	"github.com/nxadm/tail"
)

func main() {
	verboseFlag := flag.Bool("v", false, "Enable verbose output")
	flag.Parse()

	auditLog, _ := tail.TailFile("/var/log/audit/audit.log", tail.Config{
		// Continuously read from the very end
		Follow: true,
		Location: &tail.SeekInfo{
			Whence: io.SeekEnd,
		},
	})
	for line := range auditLog.Lines {
		auditEvent, err := auparse.ParseLogLine(line.Text)
		if err == nil && (auditEvent.RecordType == 1300) {
			// Get a dictionary of values
			parsedEvent := auditEvent.ToMapStr()
			if err == nil && (parsedEvent["tags"].([]string)[0] != "x86_64") {
				if *verboseFlag {
					fmt.Println(parsedEvent)
				}
				eventConnections := eventSockets(parsedEvent)

				if isBindShell(eventConnections) {
					fmt.Println("Bind Shell")
					fmt.Println(eventConnections)
				} else if isReverseShell(eventConnections) {
					fmt.Println("Reverse Shell")
					fmt.Println(eventConnections)
				} else {
					fmt.Println("SSH")
				}
			}
		}
	}
}
