package main

import (
	"github.com/cakturk/go-netstat/netstat"
)

// Get all established connections and return their sockets
func establishedConnections() []netstat.SockTabEntry {
	var establishedSockets []netstat.SockTabEntry

	// Get all TCP sockets that have an established state
	tcpSockets, _ := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
		return s.State == netstat.Established
	})
	// Get all UDP sockets that have an established state
	udpSockets, _ := netstat.UDPSocks(func(s *netstat.SockTabEntry) bool {
		return s.State == netstat.Established
	})

	// Combine TCP and UDP
	establishedSockets = append(establishedSockets, tcpSockets...)
	establishedSockets = append(establishedSockets, udpSockets...)

	return establishedSockets
}

// Get all listening connections and return their sockets
func listeningConnections() []netstat.SockTabEntry {
	var listeningSockets []netstat.SockTabEntry

	tcpSockets, _ := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
		return s.State == netstat.Listen
	})
	udpSockets, _ := netstat.UDPSocks(func(s *netstat.SockTabEntry) bool {
		return s.State == netstat.Listen
	})

	// Combine TCP and UDP
	listeningSockets = append(listeningSockets, tcpSockets...)
	listeningSockets = append(listeningSockets, udpSockets...)

	return listeningSockets
}

// Take an event and return its network-communicating sockets
func eventSockets(event map[string]interface{}) []netstat.SockTabEntry {
	var relevantSockets []netstat.SockTabEntry
	var processTree = eventProcessTree(event)
	var allSockets = append(establishedConnections(), listeningConnections()...)

	for _, pid := range processTree {
		for _, socket := range allSockets {
			if pid == socket.Process.Pid {
				relevantSockets = append(relevantSockets, socket)
			}
		}
	}
	return relevantSockets
}
