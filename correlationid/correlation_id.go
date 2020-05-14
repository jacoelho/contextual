package correlationid

import (
	"context"
	"net/http"
)

const correlationIDHeader = "X-Correlation-Id"

type correlationIDKey struct{}

// WithCorrelationID returns a new context with correlation id
func WithCorrelationID(parent context.Context, correlationID string) context.Context {
	return context.WithValue(parent, correlationIDKey{}, correlationID)
}

// FromResponse returns a new context
// Correlation id is propagated if X-Correlation-Id header is present
func FromResponse(resp *http.Response) context.Context {
	if resp == nil {
		return context.Background()
	}

	correlationID := resp.Header.Get(correlationIDHeader)
	if correlationID == "" {
		return context.Background()
	}

	return WithCorrelationID(context.Background(), correlationID)
}

// FromContext returns correlation id if present
func FromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(correlationIDKey{}).(string)
	return id, ok
}

// Handler wraps a http.Handler propagating correlation id to inner handlers
func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if correlationID := r.Header.Get(correlationIDHeader); correlationID != "" {
			w.Header().Set(correlationIDHeader, correlationID)
			h.ServeHTTP(w, r.WithContext(WithCorrelationID(r.Context(), correlationID)))
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

// Transport is a correlation id aware transport
type Transport struct {
	Transport http.RoundTripper
}

// RoundTrip implements http.RoundTripper interface
// Propagates correlation id using http header X-Correlation-Id
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.Transport
	if base == nil {
		base = http.DefaultTransport
	}

	r := req.Clone(req.Context())

	id, ok := FromContext(r.Context())
	if ok {
		r.Header.Add(correlationIDHeader, id)
	}

	return base.RoundTrip(r)
}
