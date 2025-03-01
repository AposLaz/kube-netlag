/*
 Copyright 2024 Apostolos Lazidis

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package config

import "os"

type EnvVars struct {
	NetperfPort   string
	CurrentNodeIp string
	MetricsPort   string
}

// Env returns a Config object with environment variable values. If a variable is
// unset, it will use the following default values:
// - NETPERF_PORT: 12865
// - METRICS_PORT: 9090
// - HOST_IP: "" (must be set)
func Env() EnvVars {
	netperfPort := os.Getenv("NETPERF_PORT")
	if netperfPort == "" {
		netperfPort = "12865"
	}

	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "9090"
	}

	return EnvVars{
		NetperfPort:   netperfPort,
		CurrentNodeIp: os.Getenv("HOST_IP"),
		MetricsPort:   metricsPort,
	}
}
