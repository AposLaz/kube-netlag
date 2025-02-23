package netperf

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

  "github.com/AposLaz/kube-netlag/config"
)

// ComputeLatency measures the network latency for a given IP and port using the netperf tool.
// It returns a slice containing the minimum, maximum, and mean latency values in milliseconds.
// The function runs the netperf command with a TCP_RR test and processes the output to extract
// the latency metrics. The operation is subject to a timeout to prevent hanging. In case of
// errors during command execution or output parsing, an error is returned.
func ComputeLatency(ip string, port string) ([]float64,error) {
    // Set a timeout context
    ctx, cancel := context.WithTimeout(context.Background(), 30* time.Second)
    defer cancel()  // releases resources if slowOperation completes before timeout elapses

    netperfCmd := exec.CommandContext(ctx, "netperf", "-H", ip, "-p", port, "-t", "TCP_RR", "--", "-o", "min_latency,max_latency,mean_latency")
    awkCmd := exec.Command("awk", "-F,", "/^[0-9]/ {print $1, $2, $3}")
    
    netperfOut, err := netperfCmd.StdoutPipe()
    if(err != nil) {
        return nil, fmt.Errorf("failed to get netperf output: %v", err)
    }

    // This line sets the standard input (stdin) of the awk command to be the output pipe from the netperf command
    awkCmd.Stdin = netperfOut

    // capture the output of the netperf command 
    var latencyBuffer bytes.Buffer
    
    // Redirects the standard output of the awk command so that its output is written into the out buffer 
    // instead of being printed directly to the terminal
    awkCmd.Stdout = &latencyBuffer

    // Start allow the 2 commands to run simultaneously
    if err:= netperfCmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start netperf: %v", err)
    }
    
    if err:= awkCmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start awk: %v", err)
    }

    // Wait the 2 commands to finish
    if err := netperfCmd.Wait(); err != nil {
      // Check if the context must be canceled
      if ctx.Err() == context.DeadlineExceeded {
          return nil, fmt.Errorf("netperf execution for the Node [%s] timed out after %v", ip, 30*time.Second)
      }
		return nil,fmt.Errorf("netperf execution failed: %v", err)
    }

    if err := awkCmd.Wait(); err != nil {
		return nil, fmt.Errorf("awk execution failed: %v", err)
    }

    latencyArray := strings.Fields(latencyBuffer.String())
    if len(latencyArray) != 3 {
		return nil,fmt.Errorf("expected 3 latency values, got %d", len(latencyArray))
    }

    nodeLatencies := make([]float64, 0, len(latencyArray))

    for _, v := range latencyArray {
        num, err := strconv.ParseFloat(v, 64)
        if err != nil {
            
			return nil,fmt.Errorf("invalid latency value [%s]: %v", v, err)
        }
        nodeLatencies = append(nodeLatencies,num)
    }
    
	return nodeLatencies, nil
}

// StartServer launches the netperf server on the specified port. It attempts to start the server
// up to a maximum number of retries if initial attempts fail. The function logs the success or
// failure of starting the server and returns an error if all attempts are unsuccessful.
func StartServer(port string) error {
  cmd := exec.Command("netserver", "-p", port)
  maxRetries := 5
  var err error

  for attempt := 1; attempt <= maxRetries; attempt++ {
      err := cmd.Start(); 
  
      if err == nil {
        config.Logger("INFO", "Netperf server started on port %s", port)
        return nil
      }

      config.Logger("ERROR", "Failed to start netserver (Attempt %d/%d): %v", attempt, maxRetries, err)

      time.Sleep(2 * time.Second)
  }

  return err
}