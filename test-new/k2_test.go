package test-new

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/reiver/go-telnet"
)

func TestKubernetes(t *testing.T) {
	//namespaceName := "testingocdp"
	//kubeResourcePath := "../kubernetes.yml"
	options := k8s.NewKubectlOptions("", "", "default")
	decision := k8s.AreAllNodesReady(t, options)
	fmt.Println("Is Kubernetes up and running ?", decision)
	k8s.RunKubectlAndGetOutputE(t, options, "cluster-info")
	k8s.RunKubectlAndGetOutputE(t, options, "get", "nodes")

	//filter := metav1.ListOptions{
	//	LabelSelector: "",
	//}

	//finterno, _ := k8s.GetNodesByFilterE(t, options, filter)
	//fmt.Println("Filtering here", finterno)

	//k8s.RunKubectlE(t, options, "get", "pods", "-o=wide", "--field-selector", "status.phase=Running")
	nodesare, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "pods", "-o=jsonpath={.items[*].metadata.name}")
	nodesstatus, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "pods", "-o=jsonpath={.items[*].status.phase}")
	statuscount := strings.Count(nodesstatus, "Running")
	arraystatus := strings.Split(nodesstatus, " ")
	fmt.Println("checking node status", nodesstatus)
	arraynode := strings.Split(nodesare, " ")
	if len(arraynode) == statuscount {
		fmt.Println("All Pods are Up and Running")
	}
	i := len(arraynode) - 1
	fmt.Println("number of pods running are ", i+1)
	for j := 0; j <= (len(arraynode) - 1); j++ {
		fmt.Println("Filtering done ", arraynode[j])
		if arraystatus[i] == "Running" {
			fmt.Println("The Pods is running fine") // find something for else if it is not running
		}
		k8s.RunKubectl(t, options, "get", "pods", arraynode[i])
		i--
	}
	k8s.RunKubectlE(t, options, "get", "pods", "--field-selector", "status.phase=Running")

	servicesare, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "svc", "-o=jsonpath={.items[*].metadata.name}")
	fmt.Println("services are", servicesare)
	//Kubectl get svc  -o=jsonpath={.items[*].spec.ports[*].targetPort}
	//kubectl get svc hello-kubernetes -n default -o=jsonpath="{.status.loadBalancer.ingress[*]}"
	arrayservice := strings.Split(servicesare, " ")
	a := len(arrayservice) - 1
	for b := 0; b <= a; b++ {
		service := k8s.GetService(t, options, arrayservice[b])
		externalip, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "svc", arrayservice[b], "-o=jsonpath={.status.loadbalancer.ingress[*].hostname}")

		fmt.Println("found external ip: ", externalip)
		endpoint := k8s.GetServiceEndpoint(t, options, service, 8080)
		fmt.Println("test connection with ", endpoint)
		var caller telnet.Caller = telnet.StandardCaller
		telnet.DialToAndCall(endpoint, caller)

		a--
	}
	//k8s.RunKubectlAndGetOutputE(t, options, "get", "svc", "-o=jsonpath={.items[*].status.typeid}")
	//k8s.GetService(t, options, "mysite")

	//k8s.CreateNamespace(t, options, namespaceName)
	//Pod := "zookeeper-server-6f86f757d8-cwwkw"
	//defer k8s.DeleteNamespace(t, options, namespaceName)
	//defer k8s.KubectlDelete(t, options, kubeResourcePath)

	//k8s.KubectlApply(t, options, kubeResourcePath)
	//k8s.WaitUntilServiceAvailable(t, options, "zookeeper", 10, 5*time.Second)

	//fmt.Println(*exec.Command("bash", "kubectl get pods"))
	// exec.Command("bash", "kubectl get pods")
	// uses the current context in kubeconfig
	// path-to-kubeconfig -- for example, /root/.kube/config

	//readynodes := k8s.GetReadyNodes(t, options)
	//fmt.Println("ready nodes", readynodes)

	//nodedes := k8s.GetNodes(t, options)
	//fmt.Println("Nodes description", nodedes) // Prints all available nodes
	//namepods := k8s.GetPod(t, options, Pod)
	//fmt.Println("Pod description", namepods) //Prints all the description of pod name defined

	//service := k8s.GetService(t, options, "mysite")
	//endpoint := k8s.GetServiceEndpoint(t, options, service, 8080)

	//fmt.Println("test connection with ", endpoint)

	//var caller telnet.Caller = telnet.StandardCaller

	//	telnet.DialToAndCall(endpoint, caller)
}
