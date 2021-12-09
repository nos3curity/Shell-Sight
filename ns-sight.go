package main

import (
	"github.com/cakturk/go-netstat/netstat"
)

// Get all established connections and return their sockets
func getConnectionSockets(state netstat.SkState) []netstat.SockTabEntry {
	var sockets []netstat.SockTabEntry

	// Get all TCP sockets with specified state
	tcpSockets, _ := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
		return s.State == state
	})
	// Get all UDP sockets with specified state
	udpSockets, _ := netstat.UDPSocks(func(s *netstat.SockTabEntry) bool {
		return s.State == state
	})

	// Combine TCP and UDP
	sockets = append(sockets, tcpSockets...)
	sockets = append(sockets, udpSockets...)

	return sockets
}

// Take an event and return its network-communicating sockets
func eventSockets(event map[string]interface{}) []netstat.SockTabEntry {
	var relevantSockets []netstat.SockTabEntry
	var processTree = eventProcessTree(event)
	var allSockets = append(getConnectionSockets(netstat.Established), getConnectionSockets(netstat.Listen)...)

	for _, pid := range processTree {
		for _, socket := range allSockets {
			process := socket.Process
			if process != nil {
				if pid == socket.Process.Pid {
					relevantSockets = append(relevantSockets, socket)
				}
			}
		}
	}
	return relevantSockets
}
