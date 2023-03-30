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
	serverPort          int
	serverHost          string
	websocketClientMode bool
	websocketPath       string
)

// launchClientCmd represents the launchTestClient command
var launchClientCmd = &cobra.Command{
	Use:   "launchClient",
	Short: "Launch a interactive client to test server responses.",
	Long: `This command launches a client for a websocket, TCP or a gRPC server.
	It opens a interactive prompt and allows users to send customized messages to the server and test its output.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if !websocketClientMode {
			utils.TcpClient(serverPort, serverHost)
		} else if websocketClientMode {
			utils.WebsocketClient(serverPort, serverHost, websocketPath)
		}
	},
}

func init() {
	rootCmd.AddCommand(launchClientCmd)
	launchClientCmd.Flags().IntVarP(&serverPort, "serverport", "p", 5000, "The port number on which your server is active.")
	launchClientCmd.Flags().StringVarP(&serverHost, "server", "s", "localhost", "The address where your server is active.")
	launchClientCmd.Flags().BoolVarP(&websocketClientMode, "wsmode", "w", false, "Start the client in web socket mode.")
	launchClientCmd.Flags().StringVarP(&websocketPath, "wspath", "f", "/", "The path on the server where the socket is located.")
}
