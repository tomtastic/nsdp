#!/bin/bash

# Test known TLV ranges efficiently
# Based on the TLVs mentioned: 0x0c00, 0x1000, 0x1c00, 0x2000, 0x2400, 0x2800, 0x3000, 0x3400, 0x3800, 0x4c00, 0x5400, 0x5800, 0x5c00, 0x6000, 0x6400, 0x6800, 0x6c00, 0x7000, 0x8000, 0x8c00, 0x9000

if [ $# -eq 0 ]; then
    echo "Usage: $0 <interface_name> [output_file]"
    echo "Example: $0 eth0"
    echo "Example: $0 en0 known_tlvs_results.txt"
    exit 1
fi

INTERFACE=$1
OUTPUT_FILE=${2:-"known_tlvs_$(date +%Y%m%d_%H%M%S).txt"}

echo "Testing known TLV ranges on interface: $INTERFACE"
echo "Output will be saved to: $OUTPUT_FILE"
echo ""

# Build if needed
if [ ! -f "./nsdp_discovery" ]; then
    echo "Building discovery tool..."
    ./build_discovery.sh
fi

# Known TLV ranges to test
KNOWN_TLVS=(
    "0c00"  # Port status
    "1000"  # Port statistics
    "1c00"  # Unknown
    "2000"  # VLAN engine
    "2400"  # VLAN membership
    "2800"  # 802.1Q VLAN
    "3000"  # PVID
    "3400"  # QoS engine
    "3800"  # QoS priority
    "4c00"  # Rate limiting
    "5400"  # Unknown
    "5800"  # Unknown
    "5c00"  # Port mirroring
    "6000"  # Available ports
    "6400"  # Unknown
    "6800"  # IGMP snooping
    "6c00"  # Multicast blocking
    "7000"  # IGMPv3 validation
    "8000"  # Unknown
    "8c00"  # Unknown
    "9000"  # Loop detection
)

echo "Testing ${#KNOWN_TLVS[@]} known TLV ranges..."
echo ""

# Test each known TLV individually for quick results
for tlv in "${KNOWN_TLVS[@]}"; do
    echo "Testing TLV 0x$tlv..."
    ./nsdp_discovery -i "$INTERFACE" -start "$tlv" -end "$tlv" -o "temp_${tlv}.txt" -v
    echo ""
done

# Combine all results
echo "Combining results..."
{
    echo "NSDP Known TLV Discovery Results"
    echo "================================"
    echo "Scan Date: $(date)"
    echo "Interface: $INTERFACE"
    echo "Known TLVs tested: ${#KNOWN_TLVS[@]}"
    echo ""
    
    for tlv in "${KNOWN_TLVS[@]}"; do
        if [ -f "temp_${tlv}.txt" ]; then
            echo "--- TLV 0x$tlv ---"
            cat "temp_${tlv}.txt"
            echo ""
            rm "temp_${tlv}.txt"
        fi
    done
} > "$OUTPUT_FILE"

echo "All results combined in: $OUTPUT_FILE"
echo ""
echo "Quick summary of found TLVs:"
grep -E "Valid TLVs Found: [1-9]" "$OUTPUT_FILE" | head -10
