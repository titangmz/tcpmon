package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	live := flag.Bool("live", true, "Enable live monitoring")
	interval := flag.Int("interval", 2, "Refresh interval for live monitoring (seconds)")
	flag.Parse()

	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	textView.SetBorder(true).SetTitle("TCP Monitor")

	// Channel to signal goroutine termination
	done := make(chan bool)

	// Fetch data in a goroutine
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				// Fetch and process data
				output := fetchAndProcessData()

				// Update the TextView with new data
				app.QueueUpdateDraw(func() {
					textView.SetText(output)
				})

				if !*live {
					return
				}

				time.Sleep(time.Duration(*interval) * time.Second)
			}
		}
	}()

	// Add quit keybinding
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			close(done)
			app.Stop()
			return nil
		}
		return event
	})

	if err := app.SetRoot(textView, true).Run(); err != nil {
		fmt.Printf("Error starting application: %v\n", err)
	}
}

func fetchAndProcessData() string {
	var builder strings.Builder

	connections, err := readTCPConnections("/proc/net/tcp")
	if err != nil {
		return fmt.Sprintf("[red]Error reading TCP connections: %v\n", err)
	}

	processMap := mapSocketToProcess()

	builder.WriteString(fmt.Sprintf("[yellow]%-20s %-8s %-20s %-8s %-12s %-8s %s\n",
		"Local Address", "L-Port", "Remote Address", "R-Port", "State", "PID", "Process"))
	builder.WriteString("[green]" + strings.Repeat("-", 90) + "\n")
	for _, conn := range connections {
		if proc, ok := processMap[conn.LocalAddress+":"+strconv.Itoa(conn.LocalPort)]; ok {
			conn.ProcessID = proc.PID
			conn.ProcessName = proc.Name
		}
		builder.WriteString(fmt.Sprintf("%-20s %-8d %-20s %-8d %-12s %-8d %s\n",
			conn.LocalAddress, conn.LocalPort, conn.RemoteAddress, conn.RemotePort, conn.State, conn.ProcessID, conn.ProcessName))
	}

	return builder.String()
}
