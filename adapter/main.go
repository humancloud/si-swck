// Copyright (c) 2019 Stackinsights to present
// All rights reserved

package main

import (
	"flag"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	openapinamer "k8s.io/apiserver/pkg/endpoints/openapi"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/logs"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/apiserver"
	basecmd "sigs.k8s.io/custom-metrics-apiserver/pkg/cmd"
	generatedopenapi "sigs.k8s.io/custom-metrics-apiserver/test-adapter/generated/openapi"

	swckprov "github.com/humancloud/si-swck/adapter/pkg/provider"
)

type Adapter struct {
	basecmd.AdapterBase
	// BaseURL is the address of OAP cluster
	BaseURL string
	// MetricRegex is a regular expression to filter metrics retrieved from OAP cluster
	MetricRegex string
	// RefreshRegistryInterval is the refresh window for syncing metrics with OAP cluster
	RefreshRegistryInterval time.Duration
	// Message is printed on successful startup
	Message string
	// Namespace groups metrics into a single set in case of duplicated metric name
	Namespace string
}

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()
	klog.InitFlags(nil)

	cmd := &Adapter{}

	cmd.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(generatedopenapi.GetOpenAPIDefinitions, openapinamer.NewDefinitionNamer(apiserver.Scheme))
	cmd.OpenAPIConfig.Info.Title = "stackinsights-adapter"
	cmd.OpenAPIConfig.Info.Version = "1.0.0"

	cmd.Flags().StringVar(&cmd.Message, "msg", "starting adapter...", "startup message")
	cmd.Flags().StringVar(&cmd.BaseURL, "oap-addr", "http://oap:12800/graphql", "the address of OAP cluster")
	cmd.Flags().StringVar(&cmd.MetricRegex, "metric-filter-regex", "", "a regular expression to filter metrics retrieved from OAP cluster")
	cmd.Flags().StringVar(&cmd.Namespace, "namespace", "stackinsights.apache.org", "a prefix to which metrics are appended. The format is 'namespace|metric_name'")
	cmd.Flags().DurationVar(&cmd.RefreshRegistryInterval, "refresh-interval", 10*time.Second,
		"the interval at which to update the cache of available metrics from OAP cluster")
	cmd.Flags().AddGoFlagSet(flag.CommandLine) // make sure we get the klog flags
	if err := cmd.Flags().Parse(os.Args); err != nil {
		klog.Fatalf("failed to parse arguments: %v", err)
	}

	p, err := swckprov.NewProvider(cmd.BaseURL, cmd.MetricRegex, cmd.RefreshRegistryInterval, cmd.Namespace)
	if err != nil {
		klog.Fatalf("unable to build p: %v", err)
	}
	cmd.WithExternalMetrics(p)

	klog.Infof(cmd.Message)
	if err := cmd.Run(wait.NeverStop); err != nil {
		klog.Fatalf("unable to run custom metrics adapter: %v", err)
	}
}
