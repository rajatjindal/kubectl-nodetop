package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	metricsapi "k8s.io/metrics/pkg/apis/metrics"
	metricsv1beta1api "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

func getPodNodeMap(k8sclient corev1client.CoreV1Interface) (map[string]string, error) {
	pods, err := k8sclient.Pods(v1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	m := map[string]string{}
	for _, pod := range pods.Items {
		m[pod.Name] = pod.Spec.NodeName
	}

	return m, nil
}

func getPodMetricsFromMetricsAPI(metricsclient metricsclientset.Interface, namespace, resourceName string, allNamespaces bool, labelSelector labels.Selector, fieldSelector fields.Selector) (*metricsapi.PodMetricsList, error) {
	var err error
	ns := metav1.NamespaceAll
	if !allNamespaces {
		ns = namespace
	}
	versionedMetrics := &metricsv1beta1api.PodMetricsList{}
	if resourceName != "" {
		m, err := metricsclient.MetricsV1beta1().PodMetricses(ns).Get(context.TODO(), resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		versionedMetrics.Items = []metricsv1beta1api.PodMetrics{*m}
	} else {
		versionedMetrics, err = metricsclient.MetricsV1beta1().PodMetricses(ns).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector.String(), FieldSelector: fieldSelector.String()})
		if err != nil {
			return nil, err
		}
	}
	metrics := &metricsapi.PodMetricsList{}
	err = metricsv1beta1api.Convert_v1beta1_PodMetricsList_To_metrics_PodMetricsList(versionedMetrics, metrics, nil)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func GetPodMetricsGroupedByNode(metricsclient metricsclientset.Interface, k8sclient corev1client.CoreV1Interface, namespace, resourceName string, allNamespaces bool, labelSelector labels.Selector, fieldSelector fields.Selector) (map[string][]metricsapi.PodMetrics, error) {
	metrics, err := getPodMetricsFromMetricsAPI(metricsclient, namespace, resourceName, allNamespaces, labelSelector, fieldSelector)
	if err != nil {
		return nil, err
	}

	podNodeMap, err := getPodNodeMap(k8sclient)
	if err != nil {
		return nil, err
	}

	m := map[string][]metricsapi.PodMetrics{}
	for _, metric := range metrics.Items {
		nodeName := podNodeMap[metric.Name]
		m[nodeName] = append(m[nodeName], metric)
	}

	return m, nil
}
