// Package mdmtest provides testing utilities for MDM functionality
package mdmtest

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

// TestAppleMDMClient represents a test client for Apple MDM
type TestAppleMDMClient struct {
	ServerAddress string
	SerialNumber  string
	UDID          string
	EnrollInfo    AppleEnrollInfo
}

// TestWindowsMDMClient represents a test client for Windows MDM
type TestWindowsMDMClient struct {
	ServerAddress string
	DeviceID      string
}

// AppleEnrollInfo contains enrollment information for Apple devices
type AppleEnrollInfo struct {
	UDID         string
	SerialNumber string
	DeviceType   string
	Model        string
	ServerURL    string
}

// NewTestMDMClientAppleDirect creates a new Apple MDM test client
func NewTestMDMClientAppleDirect(enrollInfo AppleEnrollInfo) *TestAppleMDMClient {
	return &TestAppleMDMClient{
		EnrollInfo:   enrollInfo,
		UDID:         enrollInfo.UDID,
		SerialNumber: enrollInfo.SerialNumber,
	}
}

// NewTestMDMClientWindowsProgramatic creates a new Windows MDM test client
func NewTestMDMClientWindowsProgramatic(serverAddress, nodeKey string) *TestWindowsMDMClient {
	return &TestWindowsMDMClient{
		ServerAddress: serverAddress,
		DeviceID:      nodeKey,
	}
}

// RandSerialNumber generates a random serial number
func RandSerialNumber() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	return strings.ToUpper(hex.EncodeToString(bytes))
}

// RandUDID generates a random UDID
func RandUDID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return fmt.Sprintf("%08X-%04X-%04X-%04X-%012X",
		bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16])
}
