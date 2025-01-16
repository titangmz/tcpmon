# TCP Monitor

A simple terminal-based TCP connection monitoring tool for Linux systems. This tool provides real-time monitoring of TCP connections with process information.

## Features

- Live monitoring of TCP connections
- Display of local and remote addresses/ports
- TCP connection states visualization
- Process information (PID and process name) for each connection
- Configurable refresh interval
- Interactive terminal UI using tview

## Requirements

- Linux operating system
- Go 1.23.4 or later

## Installation

```bash
git clone https://github.com/titangmz/tcpmon
cd tcpmon
go build
```

## Usage

Run the application with default settings:
```bash
./tcpmon
```

Available flags:
- `-live`: Enable/disable live monitoring (default: true)
- `-interval`: Set refresh interval in seconds (default: 2)

## Controls

- Press `q` to quit the application

## Output Format

The tool displays the following information for each TCP connection:
- Local Address and Port
- Remote Address and Port
- Connection State
- Process ID (PID)
- Process Name