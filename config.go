package main

import (
	"os"
	"strconv"
	"strings"
)

type paamConfig struct {
	HTTPBindAddress    string
	OnlyFailedResults  bool
	ExcludeNamespaces  []string
	ExcludeDeployments []string
}

var config *paamConfig

func configureApp() {
	httpBindAddress := envOrDefault("PAAM_HTTP_BIND_ADDRESS", ":8113")
	onlyFailedResults, err := strconv.ParseBool(envOrDefault("PAAM_ONLY_FAILED_RESULTS", "false"))
	if err != nil {
		onlyFailedResults = false
	}
	excludedNamespaces := strings.Split(envOrDefault("PAAM_EXCLUDE_NAMESPACES", ""), ",")
	excludedDeployments := strings.Split(envOrDefault("PAAM_EXCLUDE_DEPLOYMENTS", ""), ",")

	config = &paamConfig{
		HTTPBindAddress:    httpBindAddress,
		OnlyFailedResults:  onlyFailedResults,
		ExcludeNamespaces:  excludedNamespaces,
		ExcludeDeployments: excludedDeployments,
	}
}

func envOrDefault(key string, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	return value
}
