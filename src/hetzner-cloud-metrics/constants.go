package main

const name = "hetzner-cloud-metrics"
const version = "1.0.0"
const userAgent = name + "/" + version + " (https://git.ypbind.de/cgit/hetzner-cloud-status/)"

const versionText = `%s version %s
Copyright (C) 2020 by Andreas Maus <maus@ypbind.de>
This program comes with ABSOLUTELY NO WARRANTY.

%s is distributed under the Terms of the GNU General
Public License Version 3. (http://www.gnu.org/copyleft/gpl.html)

Build with go version: %s
`

const helpText = `
Usage: %s [--help] --token-file <file> --cpu --disk --network --step <step> [<server> <server> ...]
	--cpu				Show CPU data
	
	--disk				Show disk data

	--help				This text

	--network			Show network data

	--step <step>		Data range in seconds.
						Default: %d

	--token-file <file>	Token file

`

const defaultStep = 60

const hetznerAllServersURL = "https://api.hetzner.cloud/v1/servers?sort=name"
const hetznerMetricsURLFmtString = "https://api.hetzner.cloud/v1/servers/%d/metrics"

const hetznerMetricsTypeCPU = "type=cpu"
const hetznerMetricsTypeDisk = "type=disk"
const hetznerMetricsTypeNetwork = "type=network"

// ISO8601Format - time format for ISO8601
const ISO8601Format = "2006-01-02T15:04:05Z"

const (
	// TypeNil - interface{} is nil
	TypeNil int = iota
	// TypeBool - interface{} is bool
	TypeBool
	// TypeString - interface{} is string
	TypeString
	// TypeInt - interface{} is int
	TypeInt
	// TypeByte - interface{} is byte
	TypeByte
	// TypeFloat - interface{} is float
	TypeFloat
	// TypeOther - anything else
	TypeOther
)

// TypeNameMap - map type value to name
var TypeNameMap = []string{
	"nil",
	"bool",
	"string",
	"int",
	"byte",
	"float",
	"other",
}

var mapSI = []string{
	"",
	"k",
	"M",
	"G",
	"T",
	"P",
}
