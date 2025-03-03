package main

import (
	"fmt"
	"os/exec"
	"os"
	"strings"
	"bufio"
)
/*Utils*/
func removeTag(text, cutTerm string) string {
	return strings.ReplaceAll(text, cutTerm, "")
}
func parseStringToLists(text, cutTerm string) []string{
	return strings.Split(text, cutTerm)
}
/*Network Checks*/
func IsAliveWlan(verbose bool) bool{
	cmd := exec.Command("bash","-c","iw wlan0 info")
	output, err := cmd.Output()
	if err!=nil {
		fmt.Println("Does not exists?:", err)
		return false
			
	}
	
	if verbose == true {fmt.Println(string(output), "\n")}
	return true

}

func IsActuallyConnected(verbose bool) string{
	cmd := exec.Command("bash","-c","iw wlan0 info | grep 'ssid'")
	output, err := cmd.Output()
	if err!=nil {
			if verbose == true {fmt.Println("It's not connected at any network right now", "\n")}
			return "nil"
	}
	result := removeTag(string(output), "\tssid")
	if verbose == true{fmt.Println("Is Actually Connected to: " + result, "\n")}
	return result
	
}
func NetworkScanIntoArray(verbose bool) []string{
	cmd := exec.Command("bash","-c","iw wlan0 scan | grep 'SSID:'")
	output, err := cmd.Output()
	if err!=nil {
		if verbose == true{
			fmt.Println("Error while scanning")
		}
		return []string{}
	}
	
	result := removeTag(string(output), "\tSSID:")
	if verbose == true{
		print(string(result))
	}

	return parseStringToLists(result, "\n")
}
/*Set Auth*/
func GetConfigurationSimple(SSID, PWD string) string{
	cmd := exec.Command("bash","-c","wpa_passphrase " + SSID +" " + PWD)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error")
		return ""
	}
	fmt.Println(string(output))
	return string(output)
}
/*Configuration File*/
func WriteConfigFile(config string){
	f, err := os.Create("/etc/wpa_supplicant/wpa_supplicant.conf")
	if err!=nil {
		fmt.Println("ERROR")
		return	
	}
	l, err := f.WriteString(config)
	fmt.Println(l, "bytes written sucessfully")
	err = f.Close()
	}

func ReadConfigFile(){
	f, err := os.Open("/etc/wpa_supplicant/wpa_supplicant.conf")
	if err!=nil {
		fmt.Println("ERROR")
		return	
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan(){
		fmt.Println(scanner.Text())
		
	}
	
}
/*Terminal Interaction*/
func OpenNetworkOption(netArray []string){
	var i int
	var pwd string
	for index, ssid := range netArray {
		fmt.Println(index, ssid)
	}
	fmt.Scan(&i)
	SSIDnew := netArray[i]
	fmt.Println("PASSWORD:")
	fmt.Scan(&pwd)
	WriteConfigFile(GetConfigurationSimple(SSIDnew, pwd))
}
func main(){
/*
	fmt.Println(IsAliveWlan(false))
	fmt.Println(IsActuallyConnected(false))
	fmt.Println(NetworkScanIntoArray(false))
	s := GetConfigurationSimple("LIMA_2.4G","878915@@")
	WriteConfigFile(s)
*/
	OpenNetworkOption(NetworkScanIntoArray(false))
}
