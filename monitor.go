package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type WiFiStatus struct {
	Ssid          string
	State         string
	SignalQuality uint32
	Rssi          int32
	RssiPercent       int32
	Error         string
}

var handle uintptr

var (
	minRssi int32 = -100
	maxRssi int32 = -50
)

var interfaceStates = map[uint32]string{
	wlan_interface_state_not_ready:             "Not Ready",
	wlan_interface_state_connected:             "Connected",
	wlan_interface_state_ad_hoc_network_formed: "Adhoc",
	wlan_interface_state_disconnecting:         "Disconnecting",
	wlan_interface_state_disconnected:          "Disconnected",
	wlan_interface_state_associating:           "Associating",
	wlan_interface_state_discovering:           "Discovering",
	wlan_interface_state_authenticating:        "Authenticating",
}

func Monitor() {
	
	StartMonitor();
	
	for {
		_ = <- chMon
		chMon <- QueryMonitor()
		}
		return
}

func StartMonitor() {

	var apiVersion uint32

	// xp (5.1) uses version 1 of the api, Vista (6.0) uses version 2
	switch v, _ := syscall.GetVersion(); {
	case v&0xff == 5:
		if (v>>8)&0xFF > 0 {
			apiVersion = 1
		} else {
			// Windows 2000
			fmt.Fprintln(os.Stderr, "Requires Windows XP or later")
			os.Exit(10)
		}
	case v&0xff > 5:
		apiVersion = 2
	default:
		// Earlier than Windows 2000. Its probably crashed before reaching here!
		fmt.Fprintln(os.Stderr, "Requires Windows XP or later")
		os.Exit(10)
	}

	e := WlanOpenHandle(apiVersion, 0, &apiVersion, &handle)

	if e != ERROR_SUCCESS {
		fmt.Fprintln(os.Stderr, "WlanOpenHandle: ", e.Error())
		os.Exit(int(e))
	}
}

func StopMonitor() {
	e := WlanCloseHandle(handle, 0)
	if e != ERROR_SUCCESS {
		fmt.Fprintln(os.Stderr, "WlanCloseHandle: ", e.Error())
		os.Exit(int(e))
	}
}

func QueryMonitor() (status WiFiStatus) {

	// Pointers to memory allocated by Windows. WlanFreeMemory must be called to free it
	var (
		ilist *WLAN_INTERFACE_INFO_LIST
		nlist *WLAN_AVAILABLE_NETWORK_LIST
		blist *WLAN_BSS_LIST
	)

	var (
		guid GUID
		e    syscall.Errno
	)
	e = WlanEnumInterfaces(handle, 0, &ilist)
	if e != ERROR_SUCCESS {
		status.Error = "WlanEnumInterfaces: " + e.Error()
		return status
	}
	defer WlanFreeMemory(uintptr(unsafe.Pointer(ilist)))

	if ilist.dwNumberOfItems == 0 {
		status.Error = "No wireless interfaces found"
		return
	}

	// Search for a connected interface
	for i := 0; i < int(ilist.dwNumberOfItems); i++ {
		status.State = interfaceStates[ilist.InterfaceInfo[i].isState]
		if ilist.InterfaceInfo[i].isState == wlan_interface_state_connected {
			guid = ilist.InterfaceInfo[i].InterfaceGuid
			break
		}
		if i == MAX_INDEX {
			break
		}
	}

	if status.State != interfaceStates[wlan_interface_state_connected] {
		return
	}

	e = WlanGetAvailableNetworkList(handle, &guid,
		WLAN_AVAILABLE_NETWORK_INCLUDE_ALL_ADHOC_PROFILES&WLAN_AVAILABLE_NETWORK_INCLUDE_ALL_MANUAL_HIDDEN_PROFILES,
		0, &nlist)
	if e != ERROR_SUCCESS {
		status.Error = "WlanGetAvailableNetworkList: " + e.Error()		
		return
	}
	defer WlanFreeMemory(uintptr(unsafe.Pointer(nlist)))

	for i := 0; i < int(nlist.dwNumberOfItems); i++ {
		// Find the connected network
		if nlist.Network[i].dwFlags&WLAN_AVAILABLE_NETWORK_CONNECTED == WLAN_AVAILABLE_NETWORK_CONNECTED {
			// The character set of the SSID is undefined. For display make sure it only contains printable UTF-8 characters
			ssid := make([]byte, nlist.Network[i].dot11Ssid.uSSIDLength)
			for j, _ := range ssid {
				ssid[j] = nlist.Network[i].dot11Ssid.ucSSID[j] & 0x7F
				if ssid[j] < 32 {
					ssid[j] = 32
				}
			}
			status.Ssid = string(ssid)

			e = WlanGetNetworkBssList(handle,
				&guid,
				&(nlist.Network[i].dot11Ssid),
				nlist.Network[i].dot11BssType,
				nlist.Network[i].bSecurityEnabled,
				0,
				&blist)
			if e != ERROR_SUCCESS {
				status.Error = "WlanGetNetworkBssList: " + e.Error()				
				return
			}
			defer WlanFreeMemory(uintptr(unsafe.Pointer(blist)))

			// Find the best signal	
			for j := 0; j < int(blist.dwNumberOfItems); j++ {
				if blist.wlanBssEntries[j].uLinkQuality > status.SignalQuality {
					status.SignalQuality = blist.wlanBssEntries[j].uLinkQuality
					status.Rssi = blist.wlanBssEntries[j].lRssi
				}
			}
			if status.Rssi < minRssi {
				minRssi = status.Rssi
			}
			if status.Rssi > maxRssi {
				maxRssi = status.Rssi
			}
			status.RssiPercent = ((status.Rssi - minRssi) *100)/ (maxRssi - minRssi)
			return
		}
		
	}
	return
}
