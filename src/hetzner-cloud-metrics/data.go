package main

import (
	"net/http"
)

// HTTPResult - result of the http_request calls
type HTTPResult struct {
	Content    []byte
	Header     http.Header
	Status     string
	StatusCode int
	URL        string
}

// HetznerAllServer - list of all hetzner servers
type HetznerAllServer struct {
	Server []HetznerServer `json:"servers"`
}

// HetznerServer - single server
type HetznerServer struct {
	Created         string            `json:"created"`
	Datacenter      HetznerDatacenter `json:"datacenter"`
	ID              uint64            `json:"id"`
	Image           HetznerImage      `json:"image"`
	IncludedTraffic uint64            `json:"included_traffic"`
	IncomingTraffic uint64            `json:"ingoing_traffic"`
	Name            string            `json:"name"`
	OutgoingTraffic uint64            `json:"outgoing_traffic"`
	ServerType      HetznerServerType `json:"server_type"`
	Status          string            `json:"status"`
	Volumes         []uint64          `json:"volumes"`
	// ...
}

// HetznerServerType - server type
type HetznerServerType struct {
	Cores       uint64  `json:"cores"`
	Description string  `json:"description"`
	Disk        uint64  `json:"disk"`
	ID          uint64  `json:"id"`
	Memory      float64 `json:"memory"`
	Name        string  `json:"name"`
}

// HetznerDatacenter - datacenter
type HetznerDatacenter struct {
	Description string                    `json:"description"`
	ID          uint64                    `json:"id"`
	Location    HetznerDatacenterLocation `json:"location"`
	Name        string                    `json:"name"`
}

// HetznerDatacenterLocation - datacenter location
type HetznerDatacenterLocation struct {
	City        string  `json:"city"`
	Country     string  `json:"country"`
	Description string  `json:"description"`
	ID          uint64  `json:"id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Name        string  `json:"name"`
	NetworkZone string  `json:"network_zone"`
}

// HetznerImage - image info
type HetznerImage struct {
	Description string `json:"description"`
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	OSFlavor    string `json:"os_flavor"`
	OSVersion   string `json:"os_version"`
	Status      string `json:"status"`
	Type        string `json:"type"`
}

// HetznerAllVolumes - list of all volumes
type HetznerAllVolumes struct {
	Volume []HetznerVolume `json:"volumes"`
}

// HetznerVolume - single volume
type HetznerVolume struct {
	Created     string                    `json:"created"`
	Format      string                    `json:"format"`
	ID          uint64                    `json:"id"`
	LinuxDevice string                    `json:"linux_device"`
	Location    HetznerDatacenterLocation `json:"location"`
	Name        string                    `json:"name"`
	Size        uint64                    `json:"size"`
	Status      string                    `json:"status"`
}

// HetznerMetrics - metric for a server
type HetznerMetrics struct {
	Metrics HetznerMetricsData `json:"metrics"`
}

// HetznerMetricsData - metrics data
type HetznerMetricsData struct {
	End        string                   `json:"end"`
	Start      string                   `json:"start"`
	Step       float64                  `json:"step"`
	TimeSeries HetznerMetricsTimeSeries `json:"time_series"`
}

// HetznerMetricsTimeSeries - time series
type HetznerMetricsTimeSeries struct {
	CPU                 HetznerMetricsTimeSeriesData `json:"cpu"`
	Disk0BandwidthRead  HetznerMetricsTimeSeriesData `json:"disk.0.bandwidth.read"`
	Disk0BandwidthWrite HetznerMetricsTimeSeriesData `json:"disk.0.bandwidth.write"`
	Disk0IopsRead       HetznerMetricsTimeSeriesData `json:"disk.0.iops.read"`
	Disk0IopsWrite      HetznerMetricsTimeSeriesData `json:"disk.0.iops.write"`
	Network0BandwithIn  HetznerMetricsTimeSeriesData `json:"network.0.bandwidth.in"`
	Network0BandwithOut HetznerMetricsTimeSeriesData `json:"network.0.bandwidth.out"`
	Network0PpsIn       HetznerMetricsTimeSeriesData `json:"network.0.pps.in"`
	Network0PpsOut      HetznerMetricsTimeSeriesData `json:"network.0.pps.out"`
}

// HetznerTimeValue - time, value
type HetznerTimeValue = []interface{}

// HetznerMetricsTimeSeriesData - data
type HetznerMetricsTimeSeriesData struct {
	// Note: Values returned are an array containing the timestamp as uint64 and the value as string
	Values []HetznerTimeValue `json:"values"`
}
