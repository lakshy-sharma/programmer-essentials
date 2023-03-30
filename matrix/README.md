# Matrix

The matrix is a collection of networking tools written in Go using cobra CLI framework.\n
The focus is on increasing the performance of available utilities and learn more about networks.\n
It currently allows you to perform the following actions.

## Features
1. Scan a particular host for open ports.
2. Scan a network for hosts that are active. (This feature needs superuser access)
3. Launch a test TCP/Websocket server for testing your clients.
4. Launch a test TCP/Websocket client for testing your servers.

## Example
1. Find open ports on a host: <i>matrix portScan -H [IP address to scan] -s [Start port] -e [End port]</i>
2. Find active hosts on a network: <i>matrix hostScan -c [Network CIDR to scan] -t [Time to wait for Ping reply]</i>

## TODO
1. Extend the server and client to include gRPC.
2. Add feature for creating network packets for testing high speed networks.
