[![Go Report Card](https://goreportcard.com/badge/github.com/ip2location/ip2convert)](https://goreportcard.com/report/github.com/ip2location/ip2convert)

ip2convert Geolocation File Format Converter
============================================
This Go command line tool enables user to convert between the following formats:

IP2Location DB1 CSV => MMDB (compatible with GeoLite2-Country MMDB format)

IP2Location DB9 CSV => MMDB (compatible with GeoLite2-City MMDB format)

IP2Location IPv6 CSV (DB1 to DB26 supported) => IP2Location BIN (compatible with all official IP2Location SDK & libraries)


Installation
============

#### `go install` Installation

```bash
go install github.com/ip2location/ip2convert/ip2convert@latest
```


#### Git Installation

```bash
git clone https://github.com/ip2location/ip2convert ip2convert
cd ip2convert
go install ./ip2convert/
$GOPATH/bin/ip2convert
```


#### Debian/Ubuntu (amd64)

```bash
curl -LO https://github.com/ip2location/ip2convert/releases/download/v1.2.2/ip2convert-1.2.2.deb
sudo dpkg -i ip2convert-1.2.2.deb
```


#### MacOS

```
curl -Ls https://raw.githubusercontent.com/ip2location/ip2convert/main/scripts/macos.sh | sh
```


### Windows Powershell

Launch Powershell as administrator then run the below:

```bash
iwr -useb https://raw.githubusercontent.com/ip2location/ip2convert/main/scripts/windows.ps1 | iex
```


### Download pre-built binaries

Supported OS/architectures below:

```
darwin_amd64
darwin_arm64
dragonfly_amd64
freebsd_386
freebsd_amd64
freebsd_arm
freebsd_arm64
linux_386
linux_amd64
linux_arm
linux_arm64
netbsd_386
netbsd_amd64
netbsd_arm
netbsd_arm64
openbsd_386
openbsd_amd64
openbsd_arm
openbsd_arm64
solaris_amd64
windows_386
windows_amd64
windows_arm
```

After choosing a platform `PLAT` from above, run:

```bash
# for Windows, use ".zip" instead of ".tar.gz"
curl -LO https://github.com/ip2location/ip2convert/releases/download/v1.2.2/ip2convert_1.2.2_${PLAT}.tar.gz
# OR
wget https://github.com/ip2location/ip2convert/releases/download/v1.2.2/ip2convert_1.2.2_${PLAT}.tar.gz

tar -xvf ip2convert_1.2.2_${PLAT}.tar.gz
mv ip2convert_1.2.2_${PLAT} /usr/local/bin/ip2convert
```


Usage Examples
==============

### Display help
```bash
ip2convert -h
```

### Convert IP2Location DB1 IPv6 CSV into MMDB format (compatible with GeoLite2-Country MMDB format)

For the commercial DB1, please go to https://www.ip2location.com/database/db1-ip-country

For the free LITE DB1, please go to https://lite.ip2location.com/database/db1-ip-country

NOTE: Not all fields in GeoLite2-Country are supported for this conversion.

```bash
ip2convert csv2mmdb -t country -i \myfolder\IPV6-COUNTRY.CSV -o \myfolder\DB1.MMDB
```


### Convert IP2Location DB9 IPv6 CSV into MMDB format (compatible with GeoLite2-City MMDB format)

For the commercial DB9, please go to https://www.ip2location.com/database/db9-ip-country-region-city-latitude-longitude-zipcode

For the free LITE DB9, please go to https://lite.ip2location.com/database/db9-ip-country-region-city-latitude-longitude-zipcode

NOTE: Not all fields in GeoLite2-City are supported for this conversion.

```bash
ip2convert csv2mmdb -t city -i \myfolder\IPV6-COUNTRY-REGION-CITY-LATITUDE-LONGITUDE-ZIPCODE.CSV -o \myfolder\DB9.MMDB
```


### Convert IP2Location IPv6 CSV into IP2Location BIN format (compatible with official IP2Location SDK & libraries)

For the commercial CSVs, please go to https://www.ip2location.com/database/ip2location

For the free LITE CSVs, please go to https://lite.ip2location.com/ip2location-lite

DB1 to DB26 are supported.

```bash
ip2convert csv2bin -d 26 -i \myfolder\IPV6-COUNTRY-REGION-CITY-LATITUDE-LONGITUDE-ZIPCODE-TIMEZONE-ISP-DOMAIN-NETSPEED-AREACODE-WEATHER-MOBILE-ELEVATION-USAGETYPE-ADDRESSTYPE-CATEGORY-DISTRICT-ASN.CSV -o \myfolder\DB26IPV6.BIN
```


LICENCE
=====================
See the LICENSE file.
