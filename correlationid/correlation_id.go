package correlationid

import (
	"context"
	"net/http"
)

const correlationIDHeader = "X-Correlation-Id"

type correlationIDKey struct {}

func WithCorrelationID(parent context.Context, correlationID string) context.Context {
	return context.WithValue(parent, correlationIDKey{}, correlationID)
}

func FromResponse(resp *http.Response) context.Context {
	if resp == nil {
		return context.Background()
	}

	id := resp.Header.Get(correlationIDHeader)
	if id == "" {
		return context.Background()
	}

	return WithCorrelationID(context.Background(), id)
}

func FromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(correlationIDKey{}).(string)

	return id, ok
}

func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(correlationIDHeader)

		if id != "" {
			w.Header().Set(correlationIDHeader, id)
			h.ServeHTTP(w, r.WithContext(WithCorrelationID(r.Context(), id)))
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

type Transport struct {
	Transport http.RoundTripper
}

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
