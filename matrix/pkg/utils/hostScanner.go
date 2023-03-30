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
	"encoding/binary"
	"log"
	"net"
	"os"
	"time"

	"github.com/tatsushid/go-fastping"
)

type ipData struct {
	Ipaddress    string
	State        string
	Hostname     []string
	ResponseTime time.Duration
}

type pingResult struct {
	ipAddress    *net.IPAddr
	ipState      string
	responseTime time.Duration
}

// We need to convert the CIDR notation to a list of IP addresses that weaim to scan.
// Surely there should be a better way to do this but this is the best I could think of.
func convertToIPs(networkCidr string) []net.IP {
	// convert string to IPNet struct
	_, ipv4Net, err := net.ParseCIDR(networkCidr)
	if err != nil {
		log.Fatal(err)
	}

	// convert IPNet struct mask and address to uint32
	mask := binary.BigEndian.Uint32(ipv4Net.Mask)
	start := binary.BigEndian.Uint32(ipv4Net.IP)

	// find the final address
	finish := (start & mask) | (mask ^ 0xffffffff)

	// loop through addresses as uint32 and store them in a slice.
	ipStore := []net.IP{}
	for i := start; i <= finish; i++ {
		// convert back to net.IP
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		ipStore = append(ipStore, ip)
	}
	return ipStore
}

// This function is responsible for sending ping packets to all the machines within the network.
// Once we start receiving replies we send them to the output formatter for further processing.
func pingSender(ipsToScan []net.IP, pingResultChannel chan pingResult, finishChannel chan string, pingTimer int) {
	log.Println("Started the ping sender.")
	
	// Setup the IP pinger.
	p := fastping.NewPinger()
	p.MaxRTT = time.Second*time.Duration(pingTimer) + 1
	var pingOutput pingResult

	// Resolve the ip addresses and add them to a IP address pinger.
	for _, ip := range ipsToScan {
		resolvedAddress, err := net.ResolveIPAddr("ip4:icmp", ip.String())
		if err != nil {
			os.Exit(1)
		}
		p.AddIPAddr(resolvedAddress)
	}

	// Define the action to perform on receiving the replies.
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		pingOutput.ipAddress = addr
		pingOutput.ipState = "Up"
		pingOutput.responseTime = rtt
		pingResultChannel <- pingOutput
	}

	// Once the final time passes stop the program.
	p.OnIdle = func() {
		log.Println("The ping sender has completed.")
		finishChannel <- "Completed"
	}
	
	// Run the pinger.
	err := p.Run()
	if err != nil {
		panic(err)
	}
}

// This function receives the results from the ping Sender and formats them with more information to make it presentable.
func replyReceiver(pingResultChannel chan pingResult, ipDataResults chan ipData, ipsToScan []net.IP, finishChannel chan string) {
	log.Println("Started the ping receiver.")
	// Wait for each ping result and perform a IP lookup.
	for {
		select {
		case finishMessage := <-finishChannel:
			if finishMessage == "Completed" {
				log.Println("The reply receiver has completed. Will publish the results soon.")
				finishChannel <- "Finish"
			}

		case pingOutput := <-pingResultChannel:
			var ipDetails ipData

			// Capture the IP data from the received answers.
			ipDetails.Ipaddress = pingOutput.ipAddress.String()
			ipDetails.State = pingOutput.ipState
			ipDetails.ResponseTime = pingOutput.responseTime

			// Perform a Name Lookup.
			lookup, err := net.LookupAddr(pingOutput.ipAddress.String())
			if err == nil {
				ipDetails.Hostname = lookup
			} else {
				ipDetails.Hostname = []string{"N/A"}
			}
			ipDataResults <- ipDetails
		}
	}
}

// This function is the control function which controls how the hosts are discovered within the network.
func DiscoverHosts(networkCidr string, pingTimer int) []ipData {
	// Setting up the variables and the channels for communication between the threads.
	log.Println("Matrix will be scanning the following CIDR: ", networkCidr)
	ipsToScan := convertToIPs(networkCidr)
	pingResultChannel := make(chan pingResult)
	ipDataResults := make(chan ipData)
	finishChannel := make(chan string)
	var scanResults []ipData

	// Display welcome message and progress bar.
	log.Println("The scanner has started. Please sit back and relax while the scan runs. Depending on your network space the timings might vary.")

	// Start a ping sender and the reply receiver.
	go pingSender(ipsToScan, pingResultChannel, finishChannel, pingTimer)
	go replyReceiver(pingResultChannel, ipDataResults, ipsToScan, finishChannel)

	// Capture all the replies as the replyReceiver sends them.
	// Close when the formatter says it is done.
	for {
		select {
		case ipResults := <-ipDataResults:
			scanResults = append(scanResults, ipResults)
		case finishMessage := <-finishChannel:
			if finishMessage == "Finish" {
				return scanResults
			}
		}
	}
}
