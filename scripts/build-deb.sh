#!/bin/bash

VERSION="1.0.0"

rm -rf ../dist
mkdir -p ../dist/DEBIAN/
mkdir -p ../dist/usr/local/bin/
cp ../debian/control ../dist/DEBIAN/
cd ../ip2convert
rm -f go.*
go mod init ip2convert && go mod tidy && go build -o ../dist/usr/local/bin/ip2convert
cd ..
dpkg-deb -Zgzip --build dist ip2convert-$VERSION.deb