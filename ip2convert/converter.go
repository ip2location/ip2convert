package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strings"
)

var cmdCSV2MMDBInput string
var cmdCSV2MMDBOutput string
var cmdCSV2MMDBType string
var cmdCSV2BINDBPackage string

var cmdCSV2BINInput string
var cmdCSV2BINOutput string

const version string = "1.2.0"
const programName string = "ip2convert Geolocation File Format Converter"

var showVer bool = false
var maxIPv4Range *big.Int
var maxIPv4RangePlusOne *big.Int
var maxIPv6Range *big.Int

func init() {
	maxIPv4Range = big.NewInt(4294967295)
	maxIPv4RangePlusOne = big.NewInt(4294967296)
	maxIPv6Range = big.NewInt(0)
	maxIPv6Range.SetString("340282366920938463463374607431768211455", 10)
}

func main() {
	cmdCSV2MMDB := flag.NewFlagSet("csv2mmdb", flag.ExitOnError)
	cmdCSV2MMDB.StringVar(&cmdCSV2MMDBInput, "i", "", "Input CSV file")
	cmdCSV2MMDB.StringVar(&cmdCSV2MMDBOutput, "o", "", "Output MMDB file")
	cmdCSV2MMDB.StringVar(&cmdCSV2MMDBType, "t", "", "MMDB file type")

	cmdCSV2BIN := flag.NewFlagSet("csv2bin", flag.ExitOnError)
	cmdCSV2BIN.StringVar(&cmdCSV2BINDBPackage, "d", "", "DB package")
	cmdCSV2BIN.StringVar(&cmdCSV2BINInput, "i", "", "Input CSV file")
	cmdCSV2BIN.StringVar(&cmdCSV2BINOutput, "o", "", "Output BIN file")

	flag.BoolVar(&showVer, "v", false, "Show version")

	flag.Usage = func() {
		PrintUsage()
	}

	if len(os.Args) < 2 {
		flag.Parse()
		if showVer {
			PrintVersion()
		} else {
			PrintUsage()
		}
		return
	}

	switch os.Args[1] {
	case "csv2mmdb":
		cmdCSV2MMDB.Parse(os.Args[2:])
		cmdCSV2MMDBInput = strings.TrimSpace(cmdCSV2MMDBInput)
		cmdCSV2MMDBOutput = strings.TrimSpace(cmdCSV2MMDBOutput)
		cmdCSV2MMDBType = strings.TrimSpace(cmdCSV2MMDBType)
		if cmdCSV2MMDBInput == "" {
			fmt.Println("Input file not specified.")
			return
		}
		if cmdCSV2MMDBOutput == "" {
			fmt.Println("Output file not specified.")
			return
		}
		if cmdCSV2MMDBType == "" {
			fmt.Println("MMDB type not specified.")
			return
		}
		ConvertCSV2MMDB(cmdCSV2MMDBInput, cmdCSV2MMDBOutput, cmdCSV2MMDBType)
	case "csv2bin":
		cmdCSV2BIN.Parse(os.Args[2:])
		cmdCSV2BINDBPackage = strings.TrimSpace(cmdCSV2BINDBPackage)
		cmdCSV2BINInput = strings.TrimSpace(cmdCSV2BINInput)
		cmdCSV2BINOutput = strings.TrimSpace(cmdCSV2BINOutput)
		regexDBPackage := regexp.MustCompile(`^(([1-9])|(1[0-9])|(2[0-6]))$`) // 1 to 26 for the DB packages

		if !regexDBPackage.MatchString(cmdCSV2BINDBPackage) {
			fmt.Println("DB package not specified.")
			return
		}
		if cmdCSV2BINInput == "" {
			fmt.Println("Input file not specified.")
			return
		}
		if cmdCSV2BINOutput == "" {
			fmt.Println("Output file not specified.")
			return
		}
		WriteBIN(cmdCSV2BINInput, cmdCSV2BINOutput, cmdCSV2BINDBPackage)
	default:
		flag.Parse()
		if showVer {
			PrintVersion()
		} else {
			PrintUsage()
		}
	}
}

func PrintVersion() {
	fmt.Printf("%s Version: %s\n", programName, version)
}

func PrintUsage() {
	PrintVersion()
	var usage string = `

  Usage: EXE [OPTION]

    -v                   Display the version and exit

    -h                   Print this help


To convert IP2Location DB1 CSV to MMDB (compatible with GeoLite2-Country MMDB format)

  Usage: EXE csv2mmdb -t country [OPTION]

    -i                   Specify the input path to the DB1 CSV file

    -o                   Specify the output path to the MMDB file

NOTE:

  The conversion requires the IP2Location DB1 IPv6 CSV file.

  You can either subscribe to the commercial DB1 at https://www.ip2location.com
  OR download the free LITE DB1 from https://lite.ip2location.com


To convert IP2Location DB9 CSV to MMDB (compatible with GeoLite2-City MMDB format)

  Usage: EXE csv2mmdb -t city [OPTION]

    -i                   Specify the input path to the DB6 CSV file

    -o                   Specify the output path to the MMDB file

NOTE:

  The conversion requires the IP2Location DB9 IPv6 CSV file.

  You can either subscribe to the commercial DB9 at https://www.ip2location.com
  OR download the free LITE DB9 from https://lite.ip2location.com


To convert IP2Location DB CSV to IP2Location BIN

  Usage: EXE csv2bin [OPTION]

    -d                   Specify the IP2Location DB package
                         Valid values: 1 to 26

    -i                   Specify the input path to the DB CSV file

    -o                   Specify the output path to the BIN file

NOTE:

  The conversion requires the IP2Location DB CSV file.

  You can either subscribe to the commercial DB at https://www.ip2location.com
  OR download the free LITE DB from https://lite.ip2location.com
`

	usage = strings.ReplaceAll(usage, "EXE", os.Args[0])
	fmt.Println(usage)
}
