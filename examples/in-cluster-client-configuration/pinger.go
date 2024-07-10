package main

import (
	"bytes"
	"fmt"
	probing "github.com/prometheus-community/pro-bing"
	"io/ioutil"
	"net/http"
	"time"
)

func initalizePortRange() {
	payload := bytes.NewBufferString("0 2147483647").Bytes()
	ioutil.WriteFile("/proc/sys/net/ipv4/ping_group_range", payload, 0666)
}

func tryPinging(ip string) (*probing.Statistics, error) {
	pinger, err := probing.NewPinger(ip)
	if err != nil {
		return nil, err
	}
	pinger.SetPrivileged(true)
	pinger.Count = 3
	pinger.Interval = time.Millisecond
	pinger.Timeout = time.Millisecond * 999
	pinger.ResolveTimeout = time.Second
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return nil, err
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	return stats, err
}

func curlTo(h string, port int) error {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/", h, port))
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return err
	}

	if resp.StatusCode != 200 {
		fmt.Println("Status code:", resp.StatusCode)
		fmt.Println(string(body))
	}
	return nil
}
