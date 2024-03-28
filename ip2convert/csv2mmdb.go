package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"io"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
)

func ConvertCSV2MMDB(input string, output string, mmdbType string) {
	var err error
	var inFile *os.File
	inFile, err = os.Open(input)
	if err != nil {
		fmt.Printf("Invalid input file %v.\n", input)
		return
	}
	defer inFile.Close()

	var outFile *os.File
	outFile, err = os.Create(output)
	if err != nil {
		fmt.Printf("Could not create output file %v.\n", output)
		return
	}
	defer outFile.Close()

	delim := ','
	ipVersion := 6 // default should be 6 which should cover both IPv4 and IPv6

	var rdr *csv.Reader

	var dbDesc string

	if mmdbType == "country" {
		dbDesc = "GeoLite2Country database" // need this to be able to use the Maxmind API for GeoLite2 Country
	} else if mmdbType == "city" {
		dbDesc = "GeoLite2City database" // need this to be able to use the Maxmind API for GeoLite2 City
	} else {
		fmt.Println("Invalid MMDB type.")
		return
	}
	var tree *mmdbwriter.Tree

	inFileBuffered := bufio.NewReaderSize(inFile, 65536)

	entryCnt := 0
	csvRdr := csv.NewReader(inFileBuffered)
	csvRdr.Comma = delim
	csvRdr.LazyQuotes = true

	rdr = csvRdr

	for {
		parts, err := rdr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Unable to read input file.")
			return
		} else if mmdbType == "country" && len(parts) != 4 {
			fmt.Println("DB1 CSV should have 4 columns.")
			return
		} else if mmdbType == "city" && len(parts) != 9 {
			fmt.Println("DB9 CSV should have 9 columns.")
			return
		}

		if tree == nil {
			tree, err = mmdbwriter.New(
				mmdbwriter.Options{
					DatabaseType: dbDesc,
					Description: map[string]string{
						"en": dbDesc,
					},
					DisableIPv4Aliasing:     false,
					IncludeReservedNetworks: true,
					Languages:               []string{"en"},
					IPVersion:               ipVersion,
				},
			)
			if err != nil {
				fmt.Println("Could not create tree.")
				return
			}
		}

		if mmdbType == "country" {
			err = AppendDB1CSVRecord(delim, parts, tree)
			if err != nil {
				fmt.Println("Invalid CSV data.")
				return
			}
		} else if mmdbType == "city" {
			err = AppendDB9CSVRecord(delim, parts, tree)
			if err != nil {
				fmt.Println("Invalid CSV data.")
				return
			}
		}

		entryCnt += 1
	}

	if entryCnt == 0 {
		fmt.Println("Nothing to import.")
		return
	}

	fmt.Fprintf(os.Stderr, "Writing to %s (%v entries)\n", output, entryCnt)
	if _, err := tree.WriteTo(outFile); err != nil {
		fmt.Println("Writing out to tree failed.")
	}
}

func DecimalToIPv4(IPNum *big.Int) (net.IP, error) {
	if IPNum == nil || IPNum.Cmp(big.NewInt(0)) < 0 || IPNum.Cmp(maxIPv4Range) > 0 {
		return nil, errors.New("Invalid IP number.")
	}

	buf := make([]byte, 4)
	bytes := IPNum.FillBytes(buf)

	ip := net.IP(bytes)
	return ip, nil
}

func DecimalToIPv6(IPNum *big.Int) (net.IP, error) {
	if IPNum == nil || IPNum.Cmp(big.NewInt(0)) < 0 || IPNum.Cmp(maxIPv6Range) > 0 {
		return nil, errors.New("Invalid IP number.")
	}

	buf := make([]byte, 16)
	bytes := IPNum.FillBytes(buf)

	ip := net.IP(bytes)
	return ip, nil
}

func AppendDB1CSVRecord(delim rune, parts []string, tree *mmdbwriter.Tree) error {
	var err error

	// these 2 fields are used for the special case where we need to split a range due the Go handling of IPv4-mapped IPv6 being treated as plain IPv4
	oriStartNum := parts[0]
	oriEndNum := parts[1]

	startNum := new(big.Int)
	startNum, _ = startNum.SetString(parts[0], 10)

	endNum := new(big.Int)
	endNum, _ = endNum.SetString(parts[1], 10)

	var startIp net.IP
	var endIp net.IP

	startIp, err = DecimalToIPv4(startNum)
	if err != nil {
		startIp, err = DecimalToIPv6(startNum)
		if err != nil {
			return err
		}
	}
	parts[0] = startIp.String()

	endIp, err = DecimalToIPv4(endNum)
	if err != nil {
		endIp, err = DecimalToIPv6(endNum)
		if err != nil {
			return err
		}
	}
	parts[1] = endIp.String()

	record := mmdbtype.Map{}

	country := mmdbtype.Map{
		"iso_code": mmdbtype.String(parts[2]),
		"names": mmdbtype.Map{
			"en": mmdbtype.String(parts[3]),
		},
	}
	record["country"] = country

	if err = tree.InsertRange(startIp, endIp, record); err != nil {
		if strings.Contains(err.Error(), "start & end IPs did not give valid range") { // special case where start IP is IPv4-mapped IPv6 (converted by Go into plain IPv4)
			// need to split into 2 ranges
			splitIPv4 := make([]string, len(parts))
			splitIPv6 := make([]string, len(parts))
			copy(splitIPv4, parts)
			copy(splitIPv6, parts)

			splitIPv4[0] = oriStartNum
			splitIPv4[1] = "281474976710655"
			splitIPv6[0] = "281474976710656"
			splitIPv6[1] = oriEndNum

			err = AppendDB1CSVRecord(delim, splitIPv4, tree)
			if err != nil {
				return err
			}
			err = AppendDB1CSVRecord(delim, splitIPv6, tree)
			if err != nil {
				return err
			}
		} else if !strings.Contains(err.Error(), "which is in an aliased network") {
			return err
		}
	}
	return nil
}

func AppendDB9CSVRecord(delim rune, parts []string, tree *mmdbwriter.Tree) error {
	var err error

	// these 2 fields are used for the special case where we need to split a range due the Go handling of IPv4-mapped IPv6 being treated as plain IPv4
	oriStartNum := parts[0]
	oriEndNum := parts[1]

	startNum := new(big.Int)
	startNum, _ = startNum.SetString(parts[0], 10)

	endNum := new(big.Int)
	endNum, _ = endNum.SetString(parts[1], 10)

	var startIp net.IP
	var endIp net.IP

	startIp, err = DecimalToIPv4(startNum)
	if err != nil {
		startIp, err = DecimalToIPv6(startNum)
		if err != nil {
			return err
		}
	}
	parts[0] = startIp.String()

	endIp, err = DecimalToIPv4(endNum)
	if err != nil {
		endIp, err = DecimalToIPv6(endNum)
		if err != nil {
			return err
		}
	}
	parts[1] = endIp.String()

	record := mmdbtype.Map{}

	country := mmdbtype.Map{
		"iso_code": mmdbtype.String(parts[2]),
		"names": mmdbtype.Map{
			"en": mmdbtype.String(parts[3]),
		},
	}
	subdivision := mmdbtype.Map{
		"names": mmdbtype.Map{
			"en": mmdbtype.String(parts[4]),
		},
	}
	subdivisions := mmdbtype.Slice{subdivision}

	city := mmdbtype.Map{
		"names": mmdbtype.Map{
			"en": mmdbtype.String(parts[5]),
		},
	}
	var lat float64
	var long float64
	if lat, err = strconv.ParseFloat(parts[6], 64); err != nil {
		return err
	}
	if long, err = strconv.ParseFloat(parts[7], 64); err != nil {
		return err
	}
	location := mmdbtype.Map{
		"latitude":  mmdbtype.Float64(lat),
		"longitude": mmdbtype.Float64(long),
	}
	postal := mmdbtype.Map{
		"code": mmdbtype.String(parts[8]),
	}
	record["country"] = country
	record["city"] = city
	record["postal"] = postal
	record["location"] = location
	record["subdivisions"] = subdivisions

	if err = tree.InsertRange(startIp, endIp, record); err != nil {
		if strings.Contains(err.Error(), "start & end IPs did not give valid range") { // special case where start IP is IPv4-mapped IPv6 (converted by Go into plain IPv4)
			// need to split into 2 ranges
			splitIPv4 := make([]string, len(parts))
			splitIPv6 := make([]string, len(parts))
			copy(splitIPv4, parts)
			copy(splitIPv6, parts)

			splitIPv4[0] = oriStartNum
			splitIPv4[1] = "281474976710655"
			splitIPv6[0] = "281474976710656"
			splitIPv6[1] = oriEndNum

			err = AppendDB9CSVRecord(delim, splitIPv4, tree)
			if err != nil {
				return err
			}
			err = AppendDB9CSVRecord(delim, splitIPv6, tree)
			if err != nil {
				return err
			}
		} else if !strings.Contains(err.Error(), "which is in an aliased network") {
			return err
		}
	}
	return nil
}
