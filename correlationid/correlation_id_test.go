package correlationid_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jacoelho/contextual/correlationid"
)

func TestContextCorrelationID(t *testing.T) {
	value := "my id"

	got, ok := correlationid.FromContext(correlationid.WithCorrelationID(context.Background(), value))
	if !ok {
		t.Fatalf("expected value in context")
	}

	if got != value {
		t.Errorf("expect %v, got %v", value, got)
	}
}

func TestContextReplaceCorrelationID(t *testing.T) {
	value := "my id"

	parent := correlationid.WithCorrelationID(context.Background(), "initial")
	got, ok := correlationid.FromContext(correlationid.WithCorrelationID(parent, value))
	if !ok {
		t.Fatalf("expected value in context")
	}

	if got != value {
		t.Errorf("expect %v, got %v", value, got)
	}
}

func TestPropagateCorrelationID(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := httptest.NewServer(correlationid.Handler(handler))
	defer srv.Close()

	client := &http.Client{
		Transport: &correlationid.Transport{},
	}

	expected := "my value"
	ctx := correlationid.WithCorrelationID(context.Background(), expected)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %v", resp)
	}

	got, ok := correlationid.FromContext(correlationid.FromResponse(resp))
	if !ok {
		t.Fatal("correlation id context not found")
	}

	if expected != got {
		t.Fatalf("expect %v, got %v", expected, got)
	}
}
