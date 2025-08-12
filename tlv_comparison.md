# NSDP TLV Type Comparison

## Wikipedia TLV Types vs Current Implementation

| TLV Type | Wikipedia Description | Currently Requested | Implementation Status |
|----------|----------------------|-------------------|---------------------|
| 0x0001 | Device model | ✅ Yes (`EmptyDeviceModel()`) | ✅ Implemented |
| 0x0003 | Device given name | ✅ Yes (`EmptyDeviceName()`) | ✅ Implemented |
| 0x0004 | Device MAC-address | ✅ Yes (`EmptyDeviceMAC()`) | ✅ Implemented |
| 0x0005 | Device system location | ❌ No | ❌ Missing |
| 0x0006 | Device current IP-address | ✅ Yes (`EmptyDeviceIP()`) | ✅ Implemented |
| 0x0007 | Device IP-network mask | ✅ Yes (`EmptyDeviceNetmask()`) | ✅ Implemented |
| 0x0008 | Router IP-address | ✅ Yes (`EmptyRouterIP()`) | ✅ Implemented |
| 0x000a | Administration password | ❌ No | ❌ Missing |
| 0x000b | DHCP Mode | ✅ Yes (`EmptyDHCPMode()`) | ✅ Implemented |
| 0x000d | Device Firmware version slot 1 | ✅ Yes (`EmptyFWVersionSlot1()`) | ✅ Implemented |
| 0x000e | Device Firmware version slot 2 | ✅ Yes (`EmptyFWVersionSlot2()`) | ✅ Implemented |
| 0x000f | Next active firmware slot after reboot | ❌ No | ❌ Missing |
| 0x0c00 | Speed/link status of ports | ❌ No | ❌ Missing |
| 0x1000 | Port Traffic Statistic | ⚠️ Partial (`EmptyPortStatistic()`) | ⚠️ Queried separately |
| 0x2800 | Get VLAN info | ❌ No | ❌ Missing |
| 0x2c00 | Delete VLAN (write only) | ❌ No | ❌ Missing (write-only) |

## Summary

**Currently Implemented (9/16):**
- Basic device identification (MAC, name, model)
- Network configuration (IP, netmask, gateway, DHCP)
- Firmware versions (both slots)
- Port statistics (queried separately)

**Missing TLV Types (7/16):**
- 0x0005: Device system location
- 0x000a: Administration password (sensitive)
- 0x000f: Next active firmware slot
- 0x0c00: Speed/link status of ports
- 0x2800: VLAN information
- 0x2c00: Delete VLAN (write-only, not for read queries)

**Recommendations:**
1. Add missing read-only TLV types to the initial query
2. Consider security implications for password-related TLVs
3. Implement port speed/link status for better port information
4. Add VLAN information querying for network topology insights
