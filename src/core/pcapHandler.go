package core

import (
	"fmt"
	"time"

	pcap "github.com/google/gopacket/pcap"
)

type PCAPType int

const (
	PCAPOnline PCAPType = iota
	PCAPOffline
)

type PCAPConfig struct {
	Type int `json:"type"` // PCAPType, Online or Offline
	// Online
	InfName string `json:"inf_name"` // Network interface name for online processing
	// Offline
	FilePath string `json:"file_name"` // PCAP file path for offline processing
	// Common
	Filter string `json:"filter"` // PCAP filter in Berkeley Packet Filter
}

type PCAPHandler struct {
	// Public
	Type    PCAPType
	Handler *pcap.Handler
	// Private
	config PCAPConfig
}

// loadConfig : Load configuration from file to struct
func (ph *PCAPHandler) loadConfig(configPath string) error {
	return nil
}

// Init : Instantiate a Handler for PCAP, given configuration as input
func (ph *PCAPHandler) Init(configPath string) error {
	// Load config, maybe a common package in furture?

	// Setup Type from config
	ph.Type = PCAPType(ph.config.Type)

	// Instantiate Handler
	var err error
	switch ph.Type {
	case PCAPOnline:
		ph.Handler, err = pcap.OpenLive(
			ph.config.InfName, // Newtork interface name
			int32(2147483647), // Snapshot length
			false,             // Promiscuous mode or not
			-1*time.Second,    // Timeout
		)
	case PCAPOffline:
		ph.Handler, err = pcap.OpenOffline(ph.config.FilePath)
	default:
		err = fmt.Errorf("Invalid PCAPType: %v", ph.Type)
	}

	// Setup filter
	if err != nil {
		return err
	}
	err = ph.Handler.SetBPFFilter(ph.config.Filter)
	return err
}

// Work : Start to capture network traffic
func (ph *PCAPHandler) Work() {
}

// Close : Close Handler before exit
func (ph *PCAPHandler) Close() error {
	ph.Handler.Close()
	return nil
}
