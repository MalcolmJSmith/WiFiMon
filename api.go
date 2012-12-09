package main

import (
	"syscall"
	"unsafe"
)

const (
	wlan_interface_state_not_ready             = iota
	wlan_interface_state_connected             = iota
	wlan_interface_state_ad_hoc_network_formed = iota
	wlan_interface_state_disconnecting         = iota
	wlan_interface_state_disconnected          = iota
	wlan_interface_state_associating           = iota
	wlan_interface_state_discovering           = iota
	wlan_interface_state_authenticating        = iota
)
const (
	ERROR_SUCCESS                                             = 0x0
	WLAN_AVAILABLE_NETWORK_INCLUDE_ALL_ADHOC_PROFILES         = 1
	WLAN_AVAILABLE_NETWORK_INCLUDE_ALL_MANUAL_HIDDEN_PROFILES = 2
	WLAN_AVAILABLE_NETWORK_CONNECTED                          = 1
	MAX_INDEX                                                 = 1000
)

var (
	hWlanOpenHandle              *syscall.LazyProc
	hWlanCloseHandle             *syscall.LazyProc
	hWlanEnumInterfaces          *syscall.LazyProc
	hWlanGetAvailableNetworkList *syscall.LazyProc
	hWlanGetNetworkBssList       *syscall.LazyProc
	hWlanFreeMemory              *syscall.LazyProc
)

type GUID struct {
	Data1 uint
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

type DOT11_SSID struct {
	uSSIDLength uint32
	ucSSID      [32]byte
}

type WLAN_INTERFACE_INFO_LIST struct {
	dwNumberOfItems uint32
	dwIndex         uint32
	InterfaceInfo   [MAX_INDEX + 1]WLAN_INTERFACE_INFO
}

type WLAN_INTERFACE_INFO struct {
	InterfaceGuid           GUID
	strInterfaceDescription [256]uint16
	isState                 uint32
}

type WLAN_AVAILABLE_NETWORK_LIST struct {
	dwNumberOfItems uint32
	dwIndex         uint32
	Network         [MAX_INDEX + 1]WLAN_AVAILABLE_NETWORK
}

type WLAN_AVAILABLE_NETWORK struct {
	strProfileName              [256]uint16
	dot11Ssid                   DOT11_SSID
	dot11BssType                uint32
	uNumberOfBssids             uint32
	bNetworkConnectable         int32
	wlanNotConnectableReason    uint32
	uNumberOfPhyTypes           uint32
	dot11PhyTypes               [8]uint32
	bMorePhyTypes               int32
	wlanSignalQuality           uint32
	bSecurityEnabled            int32
	dot11DefaultAuthAlgorithm   uint32
	dot11DefaultCipherAlgorithm uint32
	dwFlags                     uint32
	dwReserved                  uint32
}

type WLAN_BSS_LIST struct {
	dwTotalSize     uint32
	dwNumberOfItems uint32
	wlanBssEntries  [MAX_INDEX + 1]WLAN_BSS_ENTRY
}

type WLAN_BSS_ENTRY struct {
	dot11Ssid               DOT11_SSID
	uPhyId                  uint32
	dot11Bssid              [6]byte
	dot11BssType            uint32
	dot11BssPhyType         uint32
	lRssi                   int32
	uLinkQuality            uint32
	bInRegDomain            int32
	usBeaconPeriod          uint16
	ullTimestamp            uint64
	ullHostTimestamp        uint64
	usCapabilityInformation uint16
	ulChCenterFrequency     uint32
	wlanRateSet             WLAN_RATE_SET
	ulIeOffset              uint32
	ulIeSize                uint32
}

type WLAN_RATE_SET struct {
	uRateSetLength uint32
	usRateSet      [126]uint16
}

func WlanOpenHandle(dwClientVersion uint32,
	pReserved uintptr,
	pdwNegotiatedVersion *uint32,
	phClientHandle *uintptr) syscall.Errno {
	e, _, _ := hWlanOpenHandle.Call(uintptr(dwClientVersion),
		pReserved,
		uintptr(unsafe.Pointer(pdwNegotiatedVersion)),
		uintptr(unsafe.Pointer(phClientHandle)))

	return syscall.Errno(e)
}

func WlanCloseHandle(hClientHandle uintptr,
	pReserved uintptr) syscall.Errno {
	e, _, _ := hWlanCloseHandle.Call(hClientHandle,
		pReserved)

	return syscall.Errno(e)
}

func WlanEnumInterfaces(hClientHandle uintptr,
	pReserved uintptr,
	ppInterfaceList **WLAN_INTERFACE_INFO_LIST) syscall.Errno {
	e, _, _ := hWlanEnumInterfaces.Call(hClientHandle,
		pReserved,
		uintptr(unsafe.Pointer(ppInterfaceList)))

	return syscall.Errno(e)
}

func WlanGetAvailableNetworkList(hClientHandle uintptr,
	pInterfaceGuid *GUID,
	dwFlags uint32,
	pReserved uintptr,
	ppAvailableNetworkList **WLAN_AVAILABLE_NETWORK_LIST) syscall.Errno {
	e, _, _ := hWlanGetAvailableNetworkList.Call(hClientHandle,
		uintptr(unsafe.Pointer(pInterfaceGuid)),
		uintptr(dwFlags),
		pReserved,
		uintptr(unsafe.Pointer(ppAvailableNetworkList)))

	return syscall.Errno(e)
}

func WlanGetNetworkBssList(hClientHandle uintptr,
	pInterfaceGuid *GUID,
	pDot11Ssid *DOT11_SSID,
	dot11BssType uint32,
	bSecurityEnabled int32,
	pReserved uintptr,
	ppWlanBssList **WLAN_BSS_LIST) syscall.Errno {
	e, _, _ := hWlanGetNetworkBssList.Call(hClientHandle,
		uintptr(unsafe.Pointer(pInterfaceGuid)),
		uintptr(unsafe.Pointer(pDot11Ssid)),
		uintptr(dot11BssType),
		uintptr(bSecurityEnabled),
		pReserved,
		uintptr(unsafe.Pointer(ppWlanBssList)))

	return syscall.Errno(e)
}

func WlanFreeMemory(pMemory uintptr) {
	_, _, _ = hWlanFreeMemory.Call(pMemory)
}

func init() {
	hapi := syscall.NewLazyDLL("Wlanapi.dll")
	hWlanOpenHandle = hapi.NewProc("WlanOpenHandle")
	hWlanCloseHandle = hapi.NewProc("WlanCloseHandle")
	hWlanEnumInterfaces = hapi.NewProc("WlanEnumInterfaces")
	hWlanGetAvailableNetworkList = hapi.NewProc("WlanGetAvailableNetworkList")
	hWlanGetNetworkBssList = hapi.NewProc("WlanGetNetworkBssList")
	hWlanFreeMemory = hapi.NewProc("WlanFreeMemory")
}
