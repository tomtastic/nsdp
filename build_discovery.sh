#!/bin/bash

# Build script for NSDP TLV Discovery Tool

echo "Building NSDP TLV Discovery Tool..."

# Ensure go.mod exists
if [ ! -f "go.mod" ]; then
    echo "Initializing Go module..."
    go mod init nsdp-tool
fi

# Install dependencies
echo "Installing dependencies..."
go get github.com/hdecarne-github/go-nsdp

# Build the discovery tool
echo "Building nsdp_discovery..."
go build -o nsdp_discovery nsdp_discovery.go

if [ $? -eq 0 ]; then
    echo "Build successful!"
    echo ""
    echo "Usage examples:"
    echo "  # Full scan (0x0000 to 0xFFFF) - WARNING: This will take a long time!"
    echo "  ./nsdp_discovery -i eth0"
    echo ""
    echo "  # Quick scan of known ranges"
    echo "  ./nsdp_discovery -i eth0 -start 0C00 -end 9000"
    echo ""
    echo "  # Scan specific range with output file"
    echo "  ./nsdp_discovery -i eth0 -start 1000 -end 2000 -o results.txt"
    echo ""
    echo "  # Verbose mode with custom batch size and delay"
    echo "  ./nsdp_discovery -i eth0 -start 0000 -end 1000 -v -batch 50 -delay 200ms"
    echo ""
    echo "  # Fast scan of your known TLVs"
    echo "  ./nsdp_discovery -i eth0 -start 0C00 -end 0C00"  # Port status
    echo "  ./nsdp_discovery -i eth0 -start 1000 -end 1000"  # Port statistics  
    echo "  ./nsdp_discovery -i eth0 -start 2000 -end 2000"  # VLAN engine
    echo ""
else
    echo "Build failed!"
    exit 1
fi
