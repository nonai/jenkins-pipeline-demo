package nandan

import (
	"crypto/tls"
	"fmt"
	"strings"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
)

func TestKubernetes1(t *testing.T) {
	options := k8s.NewKubectlOptions("", "", "default")
	//Smoke test for Nodes Health
	k8s.AreAllNodesReady(t, options)

	//Smoke test for Pods health
	podslist, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "pods", "-o=jsonpath={.items[*].metadata.name}")
	for _, pod := range strings.Split(podslist, " ") {
		fmt.Println("Checking pod: ", pod)
		checkPod := k8s.IsPodAvailable(k8s.GetPod(t, options, pod))
		if !checkPod {
			fmt.Println("POD is DOWN: ", pod)
			t.Fail()
		}
	}
	k8s.RunKubectlE(t, options, "get", "pods", "--field-selector", "status.phase=Running")

	//Smoke test for Service health
	servicesare, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "svc", "-l", "component!=apiserver", "-o=jsonpath={.items[*].metadata.name}")
	fmt.Println("services are", servicesare)
	for _, svc := range strings.Split(servicesare, " ") {
		fmt.Println("Checking service: ", svc)
		s := k8s.GetService(t, options, svc)
		checkSvc := k8s.IsServiceAvailable(s)
		if !checkSvc {
			fmt.Println("This Service is DOWN: ", svc)
			t.Fail()
		}
	}

	// Smoke test for LB health(HTTP)
	servicesLb, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "svc", "-l", "component!=apiserver", "-o=jsonpath={.items[*].metadata.name}")
	for _, svcLb := range strings.Split(servicesLb, " ") {
		ep, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "svc", svcLb, "-o=jsonpath={.status.loadBalancer.ingress.*.ip}")
		if len(ep) != 0 {
			port, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "svc", svcLb, "-o=jsonpath={.spec.ports[*].port}")
			println("External LBs: ", svcLb)
			tlsConfig := tls.Config{}
			http_helper.HttpGetWithRetryWithCustomValidation(
				t,
				fmt.Sprintf("http://%s:%s", ep, port),
				&tlsConfig,
				5,
				10*time.Second,
				func(statusCode int, body string) bool {
					return statusCode == 200
				},
			)
		}
	}
}
