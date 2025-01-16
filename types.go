package main

type TCPConnection struct {
	LocalAddress  string
	LocalPort     int
	RemoteAddress string
	RemotePort    int
	State         string
	ProcessID     int
	ProcessName   string
}

type ProcessInfo struct {
	PID  int
	Name string
}
