package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandleOK(t *testing.T) {
	h := NewHandler()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	h.ServeHTTP(recorder, req)

	res, err := io.ReadAll(recorder.Body)
	require.NoError(t, err)
	require.Equal(t, "OK\n", string(res))
}
