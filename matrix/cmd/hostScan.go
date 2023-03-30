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
package cmd

import (
	"fmt"
	"matrix/pkg/utils"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var (
	networkCidr string
	pingTimer   int
)

// hostScanCmd represents the hostScan command
var hostScanCmd = &cobra.Command{
	Use:   "hostScan",
	Short: "Discover active hosts in your network.",
	Long: `The hostScan allows you to scan all hosts inside a network and check if they are online or not.
	It is capable of mapping IPs to their hostnames, making it easier to find a rogue raspberry pi ;)
	`,
	Run: func(cmd *cobra.Command, args []string) {
		scanResults := utils.DiscoverHosts(networkCidr, pingTimer)
		writer := tabwriter.NewWriter(os.Stdout, 1, 8, 0, ' ', tabwriter.AlignRight|tabwriter.Debug)
		fmt.Fprintln(writer, "\nScan Complete")
		fmt.Fprintln(writer, "--------------------------------------------")
		fmt.Fprintln(writer, "IP Address\tState\tHostname\tResponse Time")
		fmt.Fprintln(writer, "--------------------------------------------")
		for _, result := range scanResults {
			if result.State == "Up" {
				fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", result.Ipaddress, result.State, result.Hostname, result.ResponseTime)
			}
		}
		writer.Flush()
	},
}

func init() {
	rootCmd.AddCommand(hostScanCmd)
	hostScanCmd.Flags().IntVarP(&pingTimer, "pingtime", "t", 20, "Number of seconds to wait for a ping reply. Default is 10 seconds.")
	hostScanCmd.Flags().StringVarP(&networkCidr, "networkcidr", "n", "192.168.0.0/24", "The CIDR notation of the network you want to scan.")
}
