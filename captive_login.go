package main

import (
    "net/http"
	"net/url"
    "time"
	"strings"
	"os/exec"
	"regexp"
	"runtime"
	// "errors"
	"fmt"
)


func isCaptivePortalDetected() bool {
    url := "http://detectportal.firefox.com/canonical.html"

    client := &http.Client{
        Timeout: 5 * time.Second,
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            // Prevent redirects to detect captive portals properly
            return http.ErrUseLastResponse
        },
    }

    resp, err := client.Get(url)
    if err != nil {
        // Network error assumed as no portal detected (or could be handled differently)
        return false
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        // Non-OK status (like redirect) means captive portal is likely present
        return true
    }
	return false
}

func getSSID() (string, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Use PowerShell to get the connected Wi-Fi name as JSON
		psCmd := `Get-NetConnectionProfile | Select-Object -First 1 -ExpandProperty Name`
		cmd = exec.Command("powershell", "-NoProfile", "-Command", psCmd)

	case "darwin":
		// macOS
		cmd = exec.Command("networksetup", "-getairportnetwork", "en0")

	case "linux":
		// Linux
    	cmd = exec.Command("sh", "-c", "nmcli -t -f active,ssid dev wifi | grep '^yes' | cut -d: -f2")
	
	default:
		return "", fmt.Errorf("unsupported platform")
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "windows":
		return strings.TrimSpace(string(output)), nil

	case "darwin":
		parts := strings.Split(string(output), ": ")
		if len(parts) == 2 {
			return strings.TrimSpace(parts[1]), nil
		}
		return "", fmt.Errorf("SSID not found")

	case "linux":
		return strings.TrimSpace(string(output)), nil

	default:
		return "", fmt.Errorf("unsupported platform")
	}
}

func is_kpr(ssid string)(bool){
	re := regexp.MustCompile(`(?i)kpr`)

	if re.MatchString(ssid) {
		return true
	} else {
		return false
	}
}

func post(roll_no string) error{
	// http://172.168.64.1:2280/submit/user_login.php
	form := url.Values{}
	form.Add("usrname", roll_no)
	form.Add("newpasswd", "123456")
	form.Add("terms", "on")
	form.Add("page_sid", "internal")
	form.Add("org_url", "http://172.168.64.1:2280/")

	resp, err := http.Post(
		"http://172.168.64.1:2280/submit/user_login.php",
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	
	return err
}


func main(){
	
	for{

		a,e:=getSSID()
		if e!=nil{
			time.Sleep(time.Second*5)
			continue
		}
		if(is_kpr(a) && isCaptivePortalDetected()){
			e=post("23ad058")
			time.Sleep(time.Minute*30)
			if e!=nil{
				time.Sleep(time.Second*5)
				continue
			}
		}
		time.Sleep(time.Second*5)
	}
}