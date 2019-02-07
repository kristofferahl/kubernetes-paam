package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
)

const statusHealthy = "healthy"
const statusAlert = "alert"

func main() {
	log.SetHandler(text.New(os.Stderr))

	configureApp()

	clientset := createKubeClient()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method":    r.Method,
			"path":      r.RequestURI,
			"direction": "incoming",
		}).Infof("HTTP %s %s", r.Method, r.RequestURI)

		if strings.HasPrefix(r.RequestURI, "/favicon.ico") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		statusCode := http.StatusOK

		result, err := runPAAM(clientset)
		if err != nil {
			statusCode = http.StatusInternalServerError
		} else if result.Status != statusHealthy {
			statusCode = http.StatusFailedDependency
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		js, err := json.Marshal(result)
		if err == nil {
			w.Write(js)
		}

		log.WithFields(log.Fields{
			"status":    statusCode,
			"direction": "outgoing",
		}).Infof("HTTP %s %s - %s", r.Method, r.RequestURI, http.StatusText(statusCode))
	})

	log.Infof("paam started")

	srv := &http.Server{
		Addr: config.HTTPBindAddress,
	}
	srv.ListenAndServe()
}

type paamResult struct {
	Status      string            `json:"status"`
	Description string            `json:"description"`
	Results     []deploymenResult `json:"results"`
}

type deploymenResult struct {
	Status      string
	Description string `json:"description"`
	Deployment  string `json:"deployment"`
	Namespace   string `json:"namespace"`
	Pods        []pod  `json:"pods"`
	NodeSpread  int    `json:"nodeSpread"`
}

type pod struct {
	Name     string `json:"name"`
	NodeName string `json:"nodeName"`
}

func runPAAM(clientset *kubernetes.Clientset) (*paamResult, error) {
	failed := make([]deploymenResult, 0)
	passed := make([]deploymenResult, 0)
	deploymentsAnalyzedCount := 0
	namespacesAnalyzed := make(map[string]int)

	deps, err := clientset.AppsV1().Deployments("").List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, d := range deps.Items {
		if contains(config.ExcludeNamespaces, d.GetNamespace()) || contains(config.ExcludeDeployments, d.GetName()) {
			continue
		}

		nodes := make(map[string]int)
		podResults := make([]pod, 0)

		log.Debugf("checking pod anti affinity for deployment %s (namespace=%s)", d.GetName(), d.GetNamespace())

		for _, p := range pods.Items {
			prefix := fmt.Sprintf("%s-", d.GetName())
			if strings.HasPrefix(p.GetName(), prefix) {
				podResults = append(podResults, pod{
					Name:     p.GetName(),
					NodeName: p.Spec.NodeName,
				})
				nodes[p.Spec.NodeName] = nodes[p.Spec.NodeName] + 1
			}
		}

		r := deploymenResult{
			Status:      "",
			Description: fmt.Sprintf("%d pod(s) are spread across %d node(s).", len(podResults), len(nodes)),
			Deployment:  d.GetName(),
			Namespace:   d.GetNamespace(),
			NodeSpread:  len(nodes),
			Pods:        podResults,
		}

		if len(nodes) == 1 && len(podResults) > 1 {
			r.Status = statusAlert
			failed = append(failed, r)
		} else {
			r.Status = statusHealthy
			passed = append(passed, r)
		}

		deploymentsAnalyzedCount++
		namespacesAnalyzed[d.GetNamespace()] = namespacesAnalyzed[d.GetNamespace()] + 1
	}

	status := statusHealthy
	results := make([]deploymenResult, 0)

	if len(failed) > 0 {
		status = statusAlert
	}

	if config.OnlyFailedResults {
		results = append(results, failed...)
	} else {
		results = append(failed, passed...)
	}

	return &paamResult{
		Status:      status,
		Description: fmt.Sprintf("analyzed %d deployment(s) in %d namespace(s)", deploymentsAnalyzedCount, len(namespacesAnalyzed)),
		Results:     results,
	}, nil
}

func createKubeClient() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
