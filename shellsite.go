package main

import (
	"fmt"
	"io"

	auparse "github.com/elastic/go-libaudit/auparse"
	tail "github.com/nxadm/tail"
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
				fmt.Println(parsedEvent["tags"])
			}
		}
	}
}
