package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Utility functions
func removeTag(text, cutTerm string) string {
	return strings.ReplaceAll(text, cutTerm, "")
}

func parseStringToLists(text, cutTerm string) []string {
	return strings.Split(text, cutTerm)
}

// Network checks
func isAliveWlan(verbose bool) bool {
	cmd := exec.Command("bash", "-c", "iw wlan0 info")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Does not exist:", err)
		return false
	}

	if verbose {
		fmt.Println(string(output), "\n")
	}
	return true
}

func isActuallyConnected(verbose bool) string {
	cmd := exec.Command("bash", "-c", "iw wlan0 info | grep 'ssid'")
	output, err := cmd.Output()
	if err != nil {
		if verbose {
			fmt.Println("It's not connected to any network right now", "\n")
		}
		return "nil"
	}
	result := removeTag(string(output), "\tssid")
	if verbose {
		fmt.Println("Is actually connected to:", result, "\n")
	}
	return result
}

func networkScanIntoArray(verbose bool) []string {
	cmd := exec.Command("bash", "-c", "iw wlan0 scan | grep 'SSID:'")
	output, err := cmd.Output()
	if err != nil {
		if verbose {
			fmt.Println("Error while scanning")
		}
		return []string{}
	}

	result := removeTag(string(output), "\tSSID:")
	if verbose {
		fmt.Print(string(result))
	}

	return parseStringToLists(result, "\n")
}

// Set authentication
func getConfigurationSimple(SSID, PWD string) string {
	cmd := exec.Command("bash", "-c", "wpa_passphrase "+SSID+" "+PWD)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error")
		return ""
	}
	fmt.Println(string(output))
	return string(output)
}

// Configuration file operations
func writeConfigFile(config string) {
	f, err := os.Create("/etc/wpa_supplicant/wpa_supplicant.conf")
	if err != nil {
		fmt.Println("ERROR")
		return
	}
	defer f.Close()

	l, err := f.WriteString(config)
	if err != nil {
		fmt.Println("ERROR")
		return
	}
	fmt.Println(l, "bytes written successfully")
}

func readConfigFile() {
	f, err := os.Open("/etc/wpa_supplicant/wpa_supplicant.conf")
	if err != nil {
		fmt.Println("ERROR")
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

// Terminal interaction
func openNetworkOption(netArray []string) {
	var i int
	var pwd string
	for index, ssid := range netArray {
		fmt.Println(index, ssid)
	}
	fmt.Scan(&i)
	SSIDnew := netArray[i]
	fmt.Println("PASSWORD:")
	fmt.Scan(&pwd)
	writeConfigFile(getConfigurationSimple(SSIDnew, pwd))
}

func main() {
	openNetworkOption(networkScanIntoArray(false))
}
