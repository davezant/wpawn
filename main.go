package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func isRunningSudo(verbose bool) bool{
	cmd := exec.Command("bash", "-c", "sudo -n echo true")
	_, err := cmd.Output()
	if err != nil {
		if verbose{fmt.Println("Error: Not in sudo")}
		return false
	}
	if verbose{fmt.Println("Started with sudo permissions")}
	return true

}
// Utility functions
func removeTag(text, cutTerm string) string {
	return strings.ReplaceAll(text, cutTerm, "")
}

func parseStringToLists(text, cutTerm string) []string {
	return strings.Split(text, cutTerm)
}

// Software checks
func isInstalled(verbose bool) bool {
	cmd1 := exec.Command("bash", "-c", "iw --version")
	output1, err1 := cmd1.Output()
	cmd2 := exec.Command("bash", "-c", "wpa_supplicant --version")
	output2, err2 := cmd2.Output()
	var nonExistentSoftware []string
	if err1 != nil {
		nonExistentSoftware = append(nonExistentSoftware, "iw")
	}
	if err2 != nil {
		nonExistentSoftware = append(nonExistentSoftware, "wpa_supplicant")
	}
	if verbose {
		fmt.Println("----------SOFTWARE-CHECK----------")
		fmt.Println(string(output1))
		fmt.Println(string(output2))
		fmt.Println("----------SOFTWARE-CHECK-END----------")
	}
	if len(nonExistentSoftware) > 0 {
		fmt.Println("The following software is not installed:", nonExistentSoftware)
		return false
	}
	return true
}

func getServicesInterface(verbose bool) bool {
	cmd1 := exec.Command("bash", "-c", "systemctl status wpa_supplicant.service")
	output1, err1 := cmd1.Output()
	cmd2 := exec.Command("bash", "-c", "rc-service wpa_supplicant status")
	output2, err2 := cmd2.Output()
	if err1 != nil && err2 != nil {
		fmt.Println("wpa_supplicant service is not running")
		return false
	}
	if verbose {
		fmt.Println("----------Daemon-working----------")
		fmt.Println(string(output1))
		fmt.Println(string(output2))
	}
	return true
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
		fmt.Println("---------Wlan0-interface-detected----------")
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
		fmt.Println("----------Connection-test-----------")
		fmt.Println("Is actually connected to:", result, "\n")
	}
	return result
}

func networkScanIntoArray(verbose bool) []string {
	fmt.Println("Scanning for available network connection...")
	cmd := exec.Command("bash", "-c", "iw wlan0 scan | grep 'SSID:'")
	output, err := cmd.Output()
	if err != nil {
		if verbose {
			fmt.Println(err)
		}
		fmt.Println("Error while scanning")
		return []string{}
	}

	result := removeTag(string(output), "\tSSID:")
	if verbose {
		fmt.Print(string(result))
	}

	return parseStringToLists(result, "\n")
}

// Set authentication
func getConfigurationSimple(SSID, PWD string, verbose bool) string {
	cmd := exec.Command("bash", "-c", "wpa_passphrase "+SSID+" "+PWD)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error while creating the configuration")
		return ""
	}
	if verbose {
		fmt.Println(string(output))
	}
	return string(output)
}

// Configuration file operations
func writeConfigFile(config string, verbose bool) {
	f, err := os.Create("/etc/wpa_supplicant/wpa_supplicant.conf")
	if err != nil {
		if verbose {
			fmt.Println(err)
		}
		fmt.Println("Error while creating the file")
		return
	}
	defer f.Close()

	l, err := f.WriteString(config)
	if err != nil {
		if verbose {
			fmt.Println(err)
		}
		fmt.Println("Error while writing the file")
		return
	}

	if verbose {
		fmt.Println(l, " Bytes written successfully!")
	}
}

func readConfigFile(verbose bool) {
	f, err := os.Open("/etc/wpa_supplicant/wpa_supplicant.conf")
	if err != nil {
		fmt.Println("Error while reading the file")
		if verbose {
			fmt.Println(err)
		}

		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

// Terminal interaction
func openNetworkOption(netArray []string, verbose bool) {
	var i int
	var pwd string
	for index, ssid := range netArray {
		if(ssid != ""){
			fmt.Println("[", index, "] - ", ssid)
		}
	}

	fmt.Println("Input the number of the network you want to connect:")
	fmt.Scan(&i)
	if(i>len(netArray)){
		fmt.Println("Error: Invalid number")
		return	
	}
	
	SSIDnew := netArray[i]
	fmt.Println("Input the password of this network:")
	fmt.Scan(&pwd)
	writeConfigFile(getConfigurationSimple(SSIDnew, pwd, verbose), verbose)
}

func main() {
	if isRunningSudo(true){
		if getServicesInterface(false){
			if isInstalled(false) == true {
				openNetworkOption(networkScanIntoArray(false), false)
			}
		}
	} else {
		fmt.Println("Execute this program with sudo permissions")
	}
}

