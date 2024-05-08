package main

import (
	"errors"
	"math/big"
	"net"
	"sort"
)

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

func GetIPv4First2Octet(bigStr string) (uint32, error) {
	bigNum := new(big.Int)
	ok := false
	if bigNum, ok = bigNum.SetString(bigStr, 10); !ok {
		return 0, errors.New("Error parsing IP number.")
	}

	buf := make([]byte, 4)
	bigBytes := bigNum.FillBytes(buf) // need to fill into buffer to preserve all bytes instead of reading bigNum.Bytes()

	var res uint32 = uint32(bigBytes[0])*256 + uint32(bigBytes[1])
	return res, nil
}

func GetIPv6First2Octet(bigStr string) (uint32, error) {
	bigNum := new(big.Int)
	ok := false
	if bigNum, ok = bigNum.SetString(bigStr, 10); !ok {
		return 0, errors.New("Error parsing IP number.")
	}

	buf := make([]byte, 16)
	bigBytes := bigNum.FillBytes(buf) // need to fill into buffer to preserve all bytes instead of reading bigNum.Bytes()

	var res uint32 = uint32(bigBytes[0])*256 + uint32(bigBytes[1])
	return res, nil
}

func ForceAsIPv4(bigStr string) ([]byte, error) {
	bigNum := new(big.Int)
	ok := false
	if bigNum, ok = bigNum.SetString(bigStr, 10); !ok {
		return nil, errors.New("Error parsing IP number.")
	}

	buf := make([]byte, 4)
	bigBytes := bigNum.FillBytes(buf) // need to fill into buffer to preserve all bytes instead of reading bigNum.Bytes()

	return bigBytes, nil
}

func ForceAsIPv6(bigStr string) ([]byte, error) {
	bigNum := new(big.Int)
	ok := false
	if bigNum, ok = bigNum.SetString(bigStr, 10); !ok {
		return nil, errors.New("Error parsing IP number.")
	}

	buf := make([]byte, 16)
	bigBytes := bigNum.FillBytes(buf) // need to fill into buffer to preserve all bytes instead of reading bigNum.Bytes()

	return bigBytes, nil
}

func IsIPv4(IP string) bool {
	ipaddr := net.ParseIP(IP)

	if ipaddr == nil {
		return false
	}

	v4 := ipaddr.To4()

	if v4 == nil {
		return false
	}

	return true
}

func IsIPv6(IP string) bool {
	if IsIPv4(IP) {
		return false
	}

	ipaddr := net.ParseIP(IP)

	if ipaddr == nil {
		return false
	}

	v6 := ipaddr.To16()

	if v6 == nil {
		return false
	}

	return true
}

func IPv4ToDecimal(IP string) (*big.Int, error) {
	if !IsIPv4(IP) {
		return nil, errors.New("Not a valid IPv4 address.")
	}

	ipnum := big.NewInt(0)
	ipaddr := net.ParseIP(IP)

	if ipaddr != nil {
		v4 := ipaddr.To4()

		if v4 != nil {
			ipnum.SetBytes(v4)
		}
	}

	return ipnum, nil
}

func IPv6ToDecimal(IP string) (*big.Int, error) {
	if !IsIPv6(IP) {
		return nil, errors.New("Not a valid IPv6 address.")
	}

	ipnum := big.NewInt(0)
	ipaddr := net.ParseIP(IP)

	if ipaddr != nil {
		v6 := ipaddr.To16()

		if v6 != nil {
			ipnum.SetBytes(v6)
		}
	}

	return ipnum, nil
}

func GetSortedKeys(myMap map[string]uint32) []string {
	keys := make([]string, 0, len(myMap))

	for k := range myMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func GetSortedKeysCountry(myMap map[string]*countryType) []string {
	keys := make([]string, 0, len(myMap))

	for k := range myMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func GetSortedKeysUint(myMap map[uint32]uint32) []uint32 {
	keys := make([]uint32, 0, len(myMap))

	for k := range myMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

func ReverseBytes(s []byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func IPv4ToBytes(IP string) ([]byte, error) {
	ipaddr := net.ParseIP(IP)

	if ipaddr == nil {
		return nil, errors.New("Invalid IP")
	}

	v4 := ipaddr.To4()

	if v4 == nil {
		return nil, errors.New("Bad IPv4")
	}

	return v4, nil
}

func IPv6ToBytes(IP string) ([]byte, error) {
	ipaddr := net.ParseIP(IP)

	if ipaddr == nil {
		return nil, errors.New("Invalid IP")
	}

	v6 := ipaddr.To16()

	if v6 == nil {
		return nil, errors.New("Bad IPv6")
	}

	return v6, nil
}

func concatSlice[T any](first []T, second []T) []T {
	n := len(first)
	return append(first[:n:n], second...)
}
