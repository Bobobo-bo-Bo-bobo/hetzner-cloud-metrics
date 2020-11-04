package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

func formatSIUnits(i float64) string {
	if i < 1024 {
		return fmt.Sprintf("%.3f", i)
	}
	idx := int(math.Floor(math.Log(i) / math.Log(1024.)))
	v := i / math.Pow(1024., float64(idx))
	return fmt.Sprintf("%.3f%s", v, mapSI[idx])
}

func formatMetrics(s HetznerAllServer, step int64, cpu bool, disk bool, network bool, useProxy string, token string) error {
	var tableHeader []string
	var metricData HetznerMetrics
	var metrics []string

	_utc, err := time.LoadLocation("UTC")
	if err != nil {
		return err
	}

	tableHeader = append(tableHeader, "Name")

	if cpu {
		metrics = append(metrics, hetznerMetricsTypeCPU)
		tableHeader = append(tableHeader, "CPU")
	}
	if disk {
		metrics = append(metrics, hetznerMetricsTypeDisk)
		tableHeader = append(tableHeader, "Disk Iops")
		tableHeader = append(tableHeader, "Disk bandwidth")
	}
	if network {
		metrics = append(metrics, hetznerMetricsTypeNetwork)
		tableHeader = append(tableHeader, "Net pps")
		tableHeader = append(tableHeader, "Net bandwidth")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	table.SetHeader(tableHeader)

	for _, srv := range s.Server {
		var tData []string

		_url := fmt.Sprintf(hetznerMetricsURLFmtString, srv.ID)

		now := time.Now().In(_utc)
		prev := time.Unix(now.Unix()-step, 0).In(_utc)

		_url = fmt.Sprintf("%s?start=%s&end=%s&%s", _url, prev.Format(ISO8601Format), now.Format(ISO8601Format), strings.Join(metrics, "&"))

		data, err := httpRequest(_url, "GET", nil, nil, useProxy, token)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if data.StatusCode != http.StatusOK {
			fmt.Fprintln(os.Stderr, data.Status)
			os.Exit(1)
		}

		err = json.Unmarshal(data.Content, &metricData)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		tData = append(tData, srv.Name)

		if cpu {
			var sum float64
			var count uint64

			if len(metricData.Metrics.TimeSeries.CPU.Values) == 0 {
				tData = append(tData, "-")
			} else {
				for _, c := range metricData.Metrics.TimeSeries.CPU.Values {
					_, val, err := extractTimeValue(c)
					if err != nil {
						return err
					}
					count++
					sum += val
				}
				avg := sum / float64(count)
				// reported value for CPU is the sum of all CPUs, normalize it
				avg /= float64(srv.ServerType.Cores)
				tData = append(tData, fmt.Sprintf("%3.2f%%", avg))
			}
		}

		if disk {
			var rsum float64
			var wsum float64
			var rcount uint64
			var wcount uint64
			var rstr = "-"
			var wstr = "-"

			for _, r := range metricData.Metrics.TimeSeries.Disk0IopsRead.Values {
				_, val, err := extractTimeValue(r)
				if err != nil {
					return err
				}
				rsum += val
				rcount++
			}

			for _, w := range metricData.Metrics.TimeSeries.Disk0IopsWrite.Values {
				_, val, err := extractTimeValue(w)
				if err != nil {
					return err
				}
				wsum += val
				wcount++
			}

			if rcount > 0 {
				rstr = formatSIUnits(rsum / float64(rcount))
			}
			if wcount > 0 {
				wstr = formatSIUnits(wsum / float64(wcount))
			}

			tData = append(tData, rstr+" / "+wstr)

			rsum = 0.0
			wsum = 0.0
			rcount = 0
			wcount = 0
			rstr = "-"
			wstr = "-"

			for _, r := range metricData.Metrics.TimeSeries.Disk0BandwidthRead.Values {
				_, val, err := extractTimeValue(r)
				if err != nil {
					return err
				}
				rsum += val
				rcount++
			}
			for _, w := range metricData.Metrics.TimeSeries.Disk0BandwidthWrite.Values {
				_, val, err := extractTimeValue(w)
				if err != nil {
					return err
				}
				wsum += val
				wcount++
			}

			if rcount > 0 {
				rstr = formatSIUnits(rsum / float64(rcount))
			}
			if wcount > 0 {
				wstr = formatSIUnits(wsum / float64(wcount))
			}

			tData = append(tData, rstr+" / "+wstr)
		}

		if network {
			var rsum float64
			var wsum float64
			var rcount uint64
			var wcount uint64
			var rstr = "-"
			var wstr = "-"

			for _, r := range metricData.Metrics.TimeSeries.Network0PpsIn.Values {
				_, val, err := extractTimeValue(r)
				if err != nil {
					return err
				}
				rsum += val
				rcount++
			}

			for _, w := range metricData.Metrics.TimeSeries.Network0PpsOut.Values {
				_, val, err := extractTimeValue(w)
				if err != nil {
					return err
				}
				wsum += val
				wcount++
			}

			if rcount > 0 {
				rstr = formatSIUnits(rsum / float64(rcount))
			}
			if wcount > 0 {
				wstr = formatSIUnits(wsum / float64(wcount))
			}

			tData = append(tData, rstr+" / "+wstr)

			rsum = 0.0
			wsum = 0.0
			rcount = 0
			wcount = 0
			rstr = "-"
			wstr = "-"

			for _, r := range metricData.Metrics.TimeSeries.Network0BandwithIn.Values {
				_, val, err := extractTimeValue(r)
				if err != nil {
					return err
				}
				rsum += val
				rcount++
			}
			for _, w := range metricData.Metrics.TimeSeries.Network0BandwithOut.Values {
				_, val, err := extractTimeValue(w)
				if err != nil {
					return err
				}
				wsum += val
				wcount++
			}

			if rcount > 0 {
				rstr = formatSIUnits(rsum / float64(rcount))
			}
			if wcount > 0 {
				wstr = formatSIUnits(wsum / float64(wcount))
			}

			tData = append(tData, rstr+" / "+wstr)

		}
		table.Append(tData)
	}
	table.Render()
	return nil
}
