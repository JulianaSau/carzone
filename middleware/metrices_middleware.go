package middleware

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// MetricsMiddleware is a middleware that logs the request and response
	// and records the duration of the request
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests made.",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests in seconds",
		},
		[]string{"method", "path"},
	)

	statusCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_status_total",
			Help: "Total number of HTTP responses by status code.",
		},
		[]string{"method", "path", "status_code"},
	)

	fuelConsumedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "fuel_consumed_liters_total",
			Help: "Total fuel consumed across all trips.",
		},
	)

	distanceTraveledTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "distance_traveled_km_total",
			Help: "Total distance traveled across all trips in kilometers.",
		},
	)

	averageTripDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "trip_duration_seconds",
			Help:    "Distribution of trip durations in seconds.",
			Buckets: prometheus.DefBuckets,
		},
	)
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func init() {
	prometheus.MustRegister(requestCounter, requestDuration, statusCounter, fuelConsumedTotal, distanceTraveledTotal, averageTripDuration)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// record the start time of the request
		start := time.Now()

		// create a new response writer
		ww := &responseWriter{ResponseWriter: w}

		// call the next handler
		next.ServeHTTP(ww, r)

		// record the duration
		duration := time.Since(start).Seconds()

		// record the request
		requestCounter.WithLabelValues(r.URL.Path, r.Method, http.StatusText(ww.statusCode)).Inc()

		// record the duration
		requestDuration.WithLabelValues(r.URL.Path, r.Method).Observe(duration)

		// record the status code
		statusCounter.WithLabelValues(r.URL.Path, r.Method, http.StatusText(ww.statusCode)).Inc()

	})
}

func RecordTripMetrics(fuelConsumed float64, distanceTraveled float64, tripDuration time.Duration) {
	fuelConsumedTotal.Add(fuelConsumed)
	distanceTraveledTotal.Add(distanceTraveled)
	averageTripDuration.Observe(tripDuration.Seconds())
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
