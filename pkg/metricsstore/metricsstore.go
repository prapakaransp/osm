package metricsstore

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// metricsRootNamespace is the root namespace for all the metrics emitted.
// Ex: osm_<metric-name>
const metricsRootNamespace = "osm"

// MetricsStore is a type that provides functionality related to metrics
type MetricsStore struct {
	// Define metrics by their category below ----------------------

	/*
	 * K8s metrics
	 */
	// K8sAPIEventCounter is the metric counter for the number of K8s API events
	K8sAPIEventCounter prometheus.Counter

	/*
	 * Proxy metrics
	 */
	// ProxyConnectCount is the metric for the total number of proxies connected to the controller
	ProxyConnectCount prometheus.Gauge

	/*
	 * Injector metrics
	 */
	// InjectorSidecarCount counts the number of injector webhooks dealt with over time
	InjectorSidecarCount prometheus.Counter

	// InjectorRqTime the histogram to track times for the injector webhook calls
	InjectorRqTime *prometheus.HistogramVec

	// MetricsStore internals should be defined below --------
	registry *prometheus.Registry
}

var defaultMetricsStore MetricsStore

// DefaultMetricsStore is the default metrics store
var DefaultMetricsStore = &defaultMetricsStore

func init() {
	/*
	 * K8s metrics
	 */
	defaultMetricsStore.K8sAPIEventCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: metricsRootNamespace,
		Subsystem: "k8s",
		Name:      "api_event_count",
		Help:      "represents the number of events received from the Kubernetes API Server",
	})

	/*
	 * Proxy metrics
	 */
	defaultMetricsStore.ProxyConnectCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: metricsRootNamespace,
		Subsystem: "proxy",
		Name:      "connect_count",
		Help:      "represents the number of proxies connected to OSM controller",
	})

	/*
	 * Injector metrics
	 */
	defaultMetricsStore.InjectorSidecarCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: metricsRootNamespace,
		Subsystem: "injector",
		Name:      "injector_sidecar_count",
		Help:      "Counts the number of injector webhooks dealt with over time",
	})

	defaultMetricsStore.InjectorRqTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metricsRootNamespace,
			Subsystem: "injector",
			Name:      "injector_rq_time",
			Buckets:   []float64{.1, .25, .5, 1, 2.5, 5, 10, 20, 40},
			Help:      "Histogram for time taken to perform sidecar injection",
		},
		[]string{
			"success",
		})

	defaultMetricsStore.registry = prometheus.NewRegistry()
}

// Start store
func (ms *MetricsStore) Start() {
	ms.registry.MustRegister(ms.K8sAPIEventCounter)
	ms.registry.MustRegister(ms.ProxyConnectCount)
	ms.registry.MustRegister(ms.InjectorSidecarCount)
	ms.registry.MustRegister(ms.InjectorRqTime)
}

// Stop store
func (ms *MetricsStore) Stop() {
	ms.registry.Unregister(ms.K8sAPIEventCounter)
	ms.registry.Unregister(ms.ProxyConnectCount)
	ms.registry.Unregister(ms.InjectorSidecarCount)
	ms.registry.Unregister(ms.InjectorRqTime)
}

// Handler return the registry
func (ms *MetricsStore) Handler() http.Handler {
	return promhttp.InstrumentMetricHandler(
		ms.registry,
		promhttp.HandlerFor(ms.registry, promhttp.HandlerOpts{}),
	)
}
