/*
Copyright © [2022] [Lakshy Sharma] <lakshy.sharma@protonmail.com>

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
package utils

import (
	"net"
	"sort"
	"strconv"
	"sync"
	"time"
)

// Cant let you run wild with this thing. Can I?
// Important. Don't be a douche and turn this setting up wildly,
// you might accidentally launch a mild DOS attack.
// Trust me Bigger is not always better.
const SANITY_LIMIT = 50

type ScanResult struct {
	Port    int
	State   string
	Service string
}

/*
Helping Functions
*/
// This function scans a port on a particular host and returns the result in a struct.
func scanPort(protocol string, hostname string, port int, portResultChannel chan ScanResult) {
	result := ScanResult{Port: port, Service: protocol}
	address := hostname + ":" + strconv.Itoa(port)
	connect, err := net.DialTimeout(protocol, address, 10*time.Second)
	if err != nil {
		result.State = "Closed"
		portResultChannel <- result
		return
	}
	defer connect.Close()
	result.State = "Open"
	portResultChannel <- result
}

func resultCollector(startPort int, endPort int, portResultChannel chan ScanResult, resultCaptureChannel chan []ScanResult) {
	var results []ScanResult

	for port := startPort; port <= endPort; port++ {
		scanOutput := <-portResultChannel
		results = append(results, scanOutput)
	}

	// Once all outputs have been collected send them back to our main thread.
	resultCaptureChannel <- results
	close(resultCaptureChannel)
	close(portResultChannel)
}

/*
Main scan controller function.
This function spawns multiple goroutines to scan the ports on a host and then waits for them to finish before moving ahead.
*/
func ScanHostPorts(hostname string, startPort int, endPort int) []ScanResult {
	log.Println("Started the port scanner.")
	speedlimitChannel := make(chan struct{}, SANITY_LIMIT)
	portResultChannel := make(chan ScanResult)
	resultCaptureChannel := make(chan []ScanResult)
	wg := sync.WaitGroup{}
	defer wg.Wait()

	// Start a receiver for capturing the outputs of our scan.
	wg.Add(1)
	go func() {
		defer wg.Done()
		resultCollector(startPort, endPort, portResultChannel, resultCaptureChannel)
	}()
	log.Println("Started the result collector.")

	// Scan Ports asynchronously.
	log.Println("Started sending SYN packets to the ports of network machines.")
	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		speedlimitChannel <- struct{}{}
		go func(hostname string, port int, returnChannel chan ScanResult) {
			defer wg.Done()
			scanPort("tcp", hostname, port, returnChannel)
			<-speedlimitChannel
		}(hostname, port, portResultChannel)
	}

	// Capture and clean the TCP scan results.
	log.Println("Completed the scan, now publishing results.")
	finalResult := <-resultCaptureChannel
	sort.SliceStable(finalResult, func(i, j int) bool {
		return finalResult[i].Port < finalResult[j].Port
	})
	return finalResult
}
