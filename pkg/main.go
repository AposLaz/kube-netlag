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

package main

import (
	"github.com/AposLaz/kube-netlag/config"
	"github.com/AposLaz/kube-netlag/promMetrics"
)

func main() {
	envVars := config.Env()

	// intialize prometheus metrics
	promMetrics.Init()
	// Initialize prometheus server
	go promMetrics.StartServer(envVars.MetricsPort)

	if err := StartNetperfServer(envVars.NetperfPort); err != nil {
		panic(err)
	}

	InitializeMonitoring(envVars)
}
