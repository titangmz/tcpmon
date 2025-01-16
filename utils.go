package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func readTCPConnections(filePath string) ([]TCPConnection, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var connections []TCPConnection
	scanner := bufio.NewScanner(file)

	// Skip the header line
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 10 {
			continue // Skip malformed lines
		}

		localAddr, localPort := parseAddress(fields[1])
		remoteAddr, remotePort := parseAddress(fields[2])
		state := parseState(fields[3])

		connections = append(connections, TCPConnection{
			LocalAddress:  localAddr,
			LocalPort:     localPort,
			RemoteAddress: remoteAddr,
			RemotePort:    remotePort,
			State:         state,
		})
	}

	return connections, scanner.Err()
}

func parseAddress(hexAddress string) (string, int) {
	parts := strings.Split(hexAddress, ":")
	if len(parts) != 2 {
		return "", 0
	}

	ipHex := parts[0]
	portHex := parts[1]

	ipBytes, _ := hex.DecodeString(ipHex)
	ip := net.IP{ipBytes[3], ipBytes[2], ipBytes[1], ipBytes[0]}.String()

	port, _ := strconv.ParseInt(portHex, 16, 32)
	return ip, int(port)
}

func parseState(hexState string) string {
	stateMap := map[string]string{
		"01": "ESTABLISHED",
		"02": "SYN_SENT",
		"03": "SYN_RECV",
		"04": "FIN_WAIT1",
		"05": "FIN_WAIT2",
		"06": "TIME_WAIT",
		"07": "CLOSE",
		"08": "CLOSE_WAIT",
		"09": "LAST_ACK",
		"0A": "LISTEN",
		"0B": "CLOSING",
	}
	return stateMap[hexState]
}

func mapSocketToProcess() map[string]ProcessInfo {
	processMap := make(map[string]ProcessInfo)

	procDirs, _ := os.ReadDir("/proc")
	for _, procDir := range procDirs {
		if pid, err := strconv.Atoi(procDir.Name()); err == nil {
			fdDir := fmt.Sprintf("/proc/%d/fd", pid)
			files, err := os.ReadDir(fdDir)
			if err != nil {
				continue
			}

			for _, file := range files {
				link, err := os.Readlink(filepath.Join(fdDir, file.Name()))
				if err != nil || !strings.HasPrefix(link, "socket:[") {
					continue
				}

				inode := strings.TrimSuffix(strings.TrimPrefix(link, "socket:["), "]")
				if socketAddr, err := getSocketAddressByInode(inode); err == nil {
					processName := getProcessName(pid)
					processMap[socketAddr] = ProcessInfo{PID: pid, Name: processName}
				}
			}
		}
	}

	return processMap
}

func getSocketAddressByInode(inode string) (string, error) {
	file, err := os.Open("/proc/net/tcp")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan() // Skip header

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		if fields[9] == inode {
			localAddr, localPort := parseAddress(fields[1])
			return localAddr + ":" + strconv.Itoa(localPort), nil
		}
	}

	return "", fmt.Errorf("inode not found")
}

func getProcessName(pid int) string {
	statusPath := fmt.Sprintf("/proc/%d/comm", pid)
	nameBytes, err := os.ReadFile(statusPath)
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(nameBytes))
}
