/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"os"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//
	// Uncomment to load all auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var (
	dnsEntries []string
)

func init() {
	dnsEntries = []string{"tsdb.service.consul", "datanode-1", "vault.service.consul"}
}

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	podName := os.Getenv("POD_NAME")
	hostName := os.Getenv("VIRTUAL_HOSTNAME")
	go StartWebServer(podName)
	initalizePortRange()
	fmt.Printf("Running in pod %s\n", podName)
	for {
		summaryBuffer := bytes.Buffer{}
		// get pods in all the namespaces by omitting namespace
		// Or specify namespace to get pods in particular namespace
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error listing pods: %s\n", err.Error())
		} else {
			fmt.Printf("%s There are %d pods in the cluster\n", podName, len(pods.Items))
			for _, p := range pods.Items {
				if !isReady(&p) {
					continue
				}
				if p.Status.PodIP == "" {
					fmt.Printf("SKIPPING because %s has no pod ip\n", p.Name)
					continue
				}
				fmt.Printf("Ping %-30s => %-30s ...", podName, p.Status.PodIP)
				printPingTest(p.Status.PodIP, p.Name)

				if strings.Contains(p.Name, "flannel-probe-") {
					// Let's try to curl to other flannel probe HTTP server ports
					if curlTo(p.Status.PodIP, 8080) != nil {
						curlProbeFailed.Inc()
						summaryBuffer.WriteString(fmt.Sprintf("CURL %s => %s failed\n", podName, p.Name))
					} else {
						curlProbeSuccess.Inc()
						summaryBuffer.WriteString(fmt.Sprintf("Node %s => %s inter-probe connectivity works\n", hostName, p.Spec.NodeName))
					}
				}
			}
		}
		fmt.Println("FLANNEL inter-pod curl summary")
		fmt.Println(summaryBuffer.String())
		fmt.Println("DNS probes")
		for _, hostname := range dnsEntries {
			fmt.Printf("Generic DNS testing %s ", hostname)
			if printPingTest(hostname, "") == nil {
				dnsSuccessful.Inc()
			} else {
				dnsFailed.Inc()
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func isReady(pod *v1.Pod) bool {
	if pod.Status.Phase == "Running" {
		return true
	}
	return false
}

func printPingTest(host string, podName string) error {
	results, err := tryPinging(host)
	if err != nil {
		fmt.Printf("FAILED - %s\n", err.Error())
		return err
	}
	success := "FAILED"
	if results.PacketsRecv == results.PacketsSent {
		success = "SUCCESS"
		pingsSuccessful.Inc()
	} else {
		pingsFailed.Inc()
		err = fmt.Errorf("pings failed")
	}
	fmt.Printf(
		"%8s %2d/%2d pings - avgRTT %3d ms maxRTT %3d ms (podName: %s)\n",
		success,
		results.PacketsRecv,
		results.PacketsSent,
		results.AvgRtt.Milliseconds(),
		results.MaxRtt.Milliseconds(),
		podName,
	)
	return err
}
