package metrics

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/prometheus/client_golang/prometheus"
)

func ConstructMetric(metric *v1alpha1.Metric, value float64) prometheus.Metric {
	labelKeys, labelValues := metric.GetMetricLabels()

	var valueType prometheus.ValueType
	switch metric.GetMetricType() {
	case v1alpha1.MetricTypeGauge:
		valueType = prometheus.GaugeValue
	}

	metricDesc := prometheus.NewDesc(metric.Name, metric.Help, labelKeys, nil)
	return prometheus.MustNewConstMetric(metricDesc, valueType, value, labelValues...)
}
