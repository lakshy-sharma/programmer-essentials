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
	"matrix/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	portNumber    int
	replyMessage  string
	websocketMode bool
)

// launchServerCmd is the command that you use for starting a test server.
var launchServerCmd = &cobra.Command{
	Use:   "launchServer",
	Short: "Start a server for testing clients.",
	Long: `This command starts a testing server which replies back with Echo of what it receives.
	In case you want to send a specific reply, you can tell the server to send back that reply for each client message.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !websocketMode {
			utils.ServeTCP(portNumber, replyMessage)
		} else if websocketMode {
			utils.ServeWebsocket(portNumber, replyMessage)
		}
	},
}

func init() {
	rootCmd.AddCommand(launchServerCmd)
	launchServerCmd.Flags().IntVarP(&portNumber, "port", "p", 5000, "The port on which to host the server.")
	launchServerCmd.Flags().StringVarP(&replyMessage, "reply", "r", "ECHO", "The reply to send when the server accepts a client message.\nECHO server is default and sends back what client sent.")
	launchServerCmd.Flags().BoolVarP(&websocketMode, "wsmode", "w", false, "Start the server in web socket mode.")
}
