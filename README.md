# TCP Health Exporter
[![Go Report Card](https://goreportcard.com/badge/github.com/xeonel2/tcp_health_exporter)](https://goreportcard.com/report/github.com/xeonel2/tcp_health_exporter)
[![GoDoc](https://godoc.org/github.com/xeonel2/tcp_health_exporter?status.svg)](https://godoc.org/github.com/xeonel2/tcp_health_exporter)

Exposes health check metrics of TCP endpoint(s) for prometheus. 

Uses TCP Shaker by Tevino. So it's silent and performs TCP handshakes without ACK.

The exporter will run on port 9112

## Requirements:
- Linux 2.4 or newer

## Usage

Example tcpservicenames.yml:
```yaml
services: 
 - {servicename: "google", host: "google.com", "port":"80","metricname":"googlehealthstatus","help":"HealthStatus for Google"}

```
