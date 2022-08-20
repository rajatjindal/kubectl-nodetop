package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	metricsapi "k8s.io/metrics/pkg/apis/metrics"
	metricsv1beta1api "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

func GetNodeMetricsGroupedByNode(metricsClient metricsclientset.Interface, resourceName string, selector labels.Selector) (map[string]metricsapi.NodeMetrics, error) {
	metrics, err := getNodeMetrics(metricsClient, resourceName, selector)
	if err != nil {
		return nil, err
	}

	m := map[string]metricsapi.NodeMetrics{}
	for _, metric := range metrics.Items {
		m[metric.Name] = metric
	}

	return m, nil
}

func getNodeMetrics(metricsClient metricsclientset.Interface, resourceName string, selector labels.Selector) (*metricsapi.NodeMetricsList, error) {
	var err error
	versionedMetrics := &metricsv1beta1api.NodeMetricsList{}
	mc := metricsClient.MetricsV1beta1()
	nm := mc.NodeMetricses()
	if resourceName != "" {
		m, err := nm.Get(context.TODO(), resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		versionedMetrics.Items = []metricsv1beta1api.NodeMetrics{*m}
	} else {
		versionedMetrics, err = nm.List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
		if err != nil {
			return nil, err
		}
	}

	metrics := &metricsapi.NodeMetricsList{}
	err = metricsv1beta1api.Convert_v1beta1_NodeMetricsList_To_metrics_NodeMetricsList(versionedMetrics, metrics, nil)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func GetAvailableResources(k8sclient corev1client.CoreV1Interface, resourceName string, selector labels.Selector, showCapacity bool) (map[string]v1.ResourceList, error) {
	var nodes []v1.Node
	if len(resourceName) > 0 {
		node, err := k8sclient.Nodes().Get(context.TODO(), resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, *node)
	} else {
		nodeList, err := k8sclient.Nodes().List(context.TODO(), metav1.ListOptions{
			LabelSelector: selector.String(),
		})
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, nodeList.Items...)
	}

	availableResources := make(map[string]v1.ResourceList)

	for _, n := range nodes {
		if !showCapacity {
			availableResources[n.Name] = n.Status.Allocatable
		} else {
			availableResources[n.Name] = n.Status.Capacity
		}
	}

	return availableResources, nil
}
