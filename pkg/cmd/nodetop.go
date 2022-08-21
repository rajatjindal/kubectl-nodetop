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

package cmd

import (
	"errors"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/discovery"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/metricsutil"
	"k8s.io/kubectl/pkg/util/completion"
	"k8s.io/kubectl/pkg/util/i18n"
	metricsapi "k8s.io/metrics/pkg/apis/metrics"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/rajatjindal/kubectl-group-top/pkg/k8s"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	sortByCPU    = "cpu"
	sortByMemory = "memory"
)

var supportedMetricsAPIVersions = []string{
	"v1beta1",
}

type TopPodOptions struct {
	NodeName           string
	Namespace          string
	LabelSelector      string
	FieldSelector      string
	SortBy             string
	AllNamespaces      bool
	PrintContainers    bool
	NoHeaders          bool
	UseProtocolBuffers bool
	Sum                bool
	ShowCapacity       bool

	Printer         *metricsutil.TopCmdPrinter
	DiscoveryClient discovery.DiscoveryInterface
	MetricsClient   metricsclientset.Interface

	factory cmdutil.Factory
	genericclioptions.IOStreams
	configFlags *genericclioptions.ConfigFlags
	k8sclient   corev1client.CoreV1Interface
}

var o *TopPodOptions

func NewCmdNodeTop(streams genericclioptions.IOStreams) *cobra.Command {
	o = &TopPodOptions{
		IOStreams:          streams,
		UseProtocolBuffers: false,
		configFlags:        genericclioptions.NewConfigFlags(true),
		Sum:                true,
		AllNamespaces:      true,
		SortBy:             "memory",
	}

	cmd := &cobra.Command{
		Use:                   "[NAME | -l label]",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Display resource (CPU/memory) usage of pods grouped by nodes"),
	}
	cmdutil.AddLabelSelectorFlagVar(cmd, &o.LabelSelector)
	cmd.Flags().StringVar(&o.FieldSelector, "field-selector", o.FieldSelector, "Selector (field query) to filter on, supports '=', '==', and '!='.(e.g. --field-selector key1=value1,key2=value2). The server only supports a limited number of field queries per type.")
	cmd.Flags().StringVar(&o.SortBy, "sort-by", o.SortBy, "If non-empty, sort pods list using specified field. The field can be either 'cpu' or 'memory'. defaults to 'memory'")
	cmd.Flags().BoolVar(&o.PrintContainers, "containers", o.PrintContainers, "If present, print usage of containers within a pod.")
	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", o.AllNamespaces, "If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
	cmd.Flags().BoolVar(&o.NoHeaders, "no-headers", o.NoHeaders, "If present, print output without headers.")
	cmd.Flags().BoolVar(&o.UseProtocolBuffers, "use-protocol-buffers", o.UseProtocolBuffers, "Enables using protocol-buffers to access Metrics API.")
	cmd.Flags().BoolVar(&o.Sum, "sum", o.Sum, "Print the sum of the resource usage")
	cmd.Flags().BoolVar(&o.ShowCapacity, "show-capacity", o.ShowCapacity, "Print node resources based on Capacity instead of Allocatable(default) of the nodes.")

	o.configFlags.AddFlags(cmd.Flags())
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(o.configFlags)
	o.factory = cmdutil.NewFactory(matchVersionKubeConfigFlags)

	cmd.ValidArgsFunction = completion.ResourceNameCompletionFunc(o.factory, "pod")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		cmdutil.CheckErr(o.Complete(o.factory, cmd, args))
		cmdutil.CheckErr(o.Validate())
		cmdutil.CheckErr(o.RunTopPod())
	}

	return cmd
}

func (o *TopPodOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	if len(args) == 1 {
		o.NodeName = args[0]
	} else if len(args) > 1 {
		return cmdutil.UsageErrorf(cmd, "%s", cmd.Use)
	}

	o.Namespace, _, err = f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}
	clientset, err := f.KubernetesClientSet()
	if err != nil {
		return err
	}

	o.DiscoveryClient = clientset.DiscoveryClient
	config, err := f.ToRESTConfig()
	if err != nil {
		return err
	}
	if o.UseProtocolBuffers {
		config.ContentType = "application/vnd.kubernetes.protobuf"
	}
	o.MetricsClient, err = metricsclientset.NewForConfig(config)
	if err != nil {
		return err
	}

	o.k8sclient = clientset.CoreV1()
	o.Printer = metricsutil.NewTopCmdPrinter(o.Out)
	return nil
}

func (o *TopPodOptions) Validate() error {
	if len(o.SortBy) > 0 {
		if o.SortBy != sortByCPU && o.SortBy != sortByMemory {
			return errors.New("--sort-by accepts only cpu or memory")
		}
	}
	if len(o.NodeName) > 0 && (len(o.LabelSelector) > 0 || len(o.FieldSelector) > 0) {
		return errors.New("only one of NAME or selector can be provided")
	}
	return nil
}

func (o TopPodOptions) RunTopPod() error {
	var err error
	labelSelector := labels.Everything()
	if len(o.LabelSelector) > 0 {
		labelSelector, err = labels.Parse(o.LabelSelector)
		if err != nil {
			return err
		}
	}
	fieldSelector := fields.Everything()
	if len(o.FieldSelector) > 0 {
		fieldSelector, err = fields.ParseSelector(o.FieldSelector)
		if err != nil {
			return err
		}
	}

	apiGroups, err := o.DiscoveryClient.ServerGroups()
	if err != nil {
		return err
	}

	metricsAPIAvailable := SupportedMetricsAPIVersionAvailable(apiGroups)

	if !metricsAPIAvailable {
		return errors.New("metrics API not available")
	}

	podMetrics, err := k8s.GetPodMetricsGroupedByNode(o.MetricsClient, o.k8sclient)
	if err != nil {
		return err
	}

	nodeMetrics, err := k8s.GetNodeMetricsGroupedByNode(o.MetricsClient, o.NodeName, labelSelector, fieldSelector)
	if err != nil {
		return err
	}

	availableResources, err := k8s.GetAvailableResources(o.k8sclient, o.NodeName, labelSelector, o.ShowCapacity)
	if err != nil {
		return err
	}

	index := 1
	for node, nodeMetric := range nodeMetrics {
		fmt.Printf("%d. Node (%s)\n", index, node)
		fmt.Println("==============================================")
		fmt.Println()
		o.Printer.PrintNodeMetrics([]metricsapi.NodeMetrics{nodeMetric}, map[string]v1.ResourceList{node: availableResources[node]}, o.NoHeaders, o.SortBy)
		fmt.Println()
		fmt.Println("Pods")
		fmt.Println("====")
		fmt.Println()
		o.Printer.PrintPodMetrics(podMetrics[node], o.PrintContainers, o.AllNamespaces, o.NoHeaders, o.SortBy, o.Sum)

		index++
	}

	return nil
}

func SupportedMetricsAPIVersionAvailable(discoveredAPIGroups *metav1.APIGroupList) bool {
	for _, discoveredAPIGroup := range discoveredAPIGroups.Groups {
		if discoveredAPIGroup.Name != metricsapi.GroupName {
			continue
		}
		for _, version := range discoveredAPIGroup.Versions {
			for _, supportedVersion := range supportedMetricsAPIVersions {
				if version.Version == supportedVersion {
					return true
				}
			}
		}
	}
	return false
}
