package metrics

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	wfv1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
)

type RealTimeMetric struct {
	Func func() float64
}

func ConstructOrUpdateMetric(metric prometheus.Metric, metricSpec *wfv1.Prometheus) (prometheus.Metric, error) {
	switch metricSpec.GetMetricType() {
	case wfv1.MetricTypeGauge:
		return constructOrUpdateGaugeMetric(metric, metricSpec)
	case wfv1.MetricTypeHistogram:
		return constructOrUpdateHistogramMetric(metric, metricSpec)
	case wfv1.MetricTypeCounter:
		return constructOrUpdateCounterMetric(metric, metricSpec)
	default:
		return nil, fmt.Errorf("invalid metric spec")
	}
}

func ConstructRealTimeGaugeMetric(metricSpec *wfv1.Prometheus, valueFunc func() float64) prometheus.Metric {
	gaugeOpts := prometheus.GaugeOpts{
		Namespace:   argoNamespace,
		Subsystem:   workflowsSubsystem,
		Name:        metricSpec.Name,
		Help:        metricSpec.Help,
		ConstLabels: metricSpec.GetMetricLabels(),
	}

	return prometheus.NewGaugeFunc(gaugeOpts, valueFunc)
}

func constructOrUpdateCounterMetric(metric prometheus.Metric, metricSpec *wfv1.Prometheus) (prometheus.Metric, error) {
	counterOpts := prometheus.CounterOpts{
		Namespace:   argoNamespace,
		Subsystem:   workflowsSubsystem,
		Name:        metricSpec.Name,
		Help:        metricSpec.Help,
		ConstLabels: metricSpec.GetMetricLabels(),
	}

	val, err := strconv.ParseFloat(metricSpec.Counter.Value, 64)
	if err != nil {
		return nil, err
	}

	if metric == nil {
		metric = prometheus.NewCounter(counterOpts)
	}
	counter := metric.(prometheus.Counter)
	counter.Add(val)
	return counter, nil
}

func constructOrUpdateGaugeMetric(metric prometheus.Metric, metricSpec *wfv1.Prometheus) (prometheus.Metric, error) {
	gaugeOpts := prometheus.GaugeOpts{
		Namespace:   argoNamespace,
		Subsystem:   workflowsSubsystem,
		Name:        metricSpec.Name,
		Help:        metricSpec.Help,
		ConstLabels: metricSpec.GetMetricLabels(),
	}

	val, err := strconv.ParseFloat(metricSpec.Gauge.Value, 64)
	if err != nil {
		return nil, err
	}

	if metric == nil {
		metric = prometheus.NewGauge(gaugeOpts)
	}
	gauge := metric.(prometheus.Gauge)
	gauge.Set(val)
	return gauge, nil
}

func constructOrUpdateHistogramMetric(metric prometheus.Metric, metricSpec *wfv1.Prometheus) (prometheus.Metric, error) {
	histOpts := prometheus.HistogramOpts{
		Namespace:   argoNamespace,
		Subsystem:   workflowsSubsystem,
		Name:        metricSpec.Name,
		Help:        metricSpec.Help,
		ConstLabels: metricSpec.GetMetricLabels(),
		Buckets:     metricSpec.Histogram.Buckets,
	}

	val, err := strconv.ParseFloat(metricSpec.Histogram.Value, 64)
	if err != nil {
		return nil, err
	}
	if metric == nil {
		metric = prometheus.NewHistogram(histOpts)
	}
	hist := metric.(prometheus.Histogram)
	hist.Observe(val)
	return hist, nil
}

func getWorkflowPhaseGauges() map[wfv1.NodePhase]prometheus.Gauge {

	getOptsByPahse := func(phase wfv1.NodePhase) prometheus.GaugeOpts {
		return prometheus.GaugeOpts{
			Namespace:   argoNamespace,
			Subsystem:   workflowsSubsystem,
			Name:        "count",
			Help:        "Number of Workflows currently accessible by the controller by status",
			ConstLabels: map[string]string{"status": string(phase)},
		}
	}
	return map[wfv1.NodePhase]prometheus.Gauge{
		wfv1.NodePending:   prometheus.NewGauge(getOptsByPahse(wfv1.NodePending)),
		wfv1.NodeRunning:   prometheus.NewGauge(getOptsByPahse(wfv1.NodeRunning)),
		wfv1.NodeSucceeded: prometheus.NewGauge(getOptsByPahse(wfv1.NodeSucceeded)),
		wfv1.NodeSkipped:   prometheus.NewGauge(getOptsByPahse(wfv1.NodeSkipped)),
		wfv1.NodeFailed:    prometheus.NewGauge(getOptsByPahse(wfv1.NodeFailed)),
		wfv1.NodeError:     prometheus.NewGauge(getOptsByPahse(wfv1.NodeError)),
	}
}
