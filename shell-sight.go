package main

import (
	"fmt"
	"io"

	"github.com/elastic/go-libaudit/auparse"
	"github.com/nxadm/tail"
)

func main() {
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
				// Get relevant process info from the event
				fmt.Println(eventProcessTree(parsedEvent))
			}
		}
	}
}
