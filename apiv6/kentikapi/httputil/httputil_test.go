package httputil_test

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/kentik/community_sdk_golang/apiv6/kentikapi/httputil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Sequence of calls for single request:
// 1. retryingClient.Do()
// 2. retryingClient.retryableRoundTripper.RoundTrip()
// 3. retryingClient.retryableRoundTripper.retryableClient.Do()
// 4. retryingClient.retryableRoundTripper.retryableClient.httpClient.Do()
// 5. retryingClient.retryableRoundTripper.retryableClient.httpClient.httpTransport.RoundTrip()

//nolint:errcheck // https://github.com/kisielk/errcheck/issues/55
// Closing a response you only read from cannot yield a meaningful error.
func TestRetryingClient_Do_ReturnsHTTPTransportError(t *testing.T) {
	t.Parallel()

	// arrange
	c := httputil.NewRetryingClient(httputil.ClientConfig{})

	req, err := retryablehttp.NewRequest(http.MethodGet, "https://invalid.url", nil)
	require.NoError(t, err)

	// act
	resp, err := c.Do(req.WithContext(context.Background()))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// assert
	t.Logf("Got response: %v, err: %v", resp, err)

	var dnsErr *net.DNSError
	require.True(t, errors.As(err, &dnsErr))
	assert.Equal(t, "no such host", dnsErr.Err)
}

//nolint:errcheck // https://github.com/kisielk/errcheck/issues/55
// Closing a response you only read from cannot yield a meaningful error.
func TestRetryingClientWithSpyHTTPTransport_Do(t *testing.T) {
	t.Parallel()

	const retryMax = 5

	tests := []struct {
		name                  string
		transportError        error
		expectedRequestsCount int
	}{
		{
			name: "retries when underlying client returns temporary URL error",
			transportError: &url.Error{
				Err: &net.OpError{
					Err: &net.DNSError{
						Err:         "fake error",
						IsTemporary: true,
					},
				},
			},
			expectedRequestsCount: retryMax + 1,
		}, {
			name: "does not retry when underlying client returns non-temporary URL error",
			transportError: &url.Error{
				Err: &net.OpError{
					Err: &net.DNSError{
						Err:         "fake error",
						IsTemporary: false,
					},
				},
			},
			expectedRequestsCount: 1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			st := spyTransport{transportError: tt.transportError}
			c := httputil.NewRetryingClient(httputil.ClientConfig{
				HTTPClient: &http.Client{
					Transport: &st,
				},
				RetryCfg: httputil.RetryConfig{
					MaxAttempts: intPtr(retryMax),
					MinDelay:    durationPtr(1 * time.Microsecond),
					MaxDelay:    durationPtr(10 * time.Microsecond),
				},
			})

			req, err := retryablehttp.NewRequest(http.MethodGet, "https://dummy.url", nil)
			require.NoError(t, err)

			// act
			resp, err := c.Do(req.WithContext(context.Background()))
			if err != nil {
				return
			}
			defer resp.Body.Close()

			// assert
			t.Logf("Got response: %v, err: %v", resp, err)
			assert.Equal(t, tt.expectedRequestsCount, st.requestsCount)

			var dnsErr *net.DNSError
			require.True(t, errors.As(err, &dnsErr))
			assert.Equal(t, "fake error", dnsErr.Err)
		})
	}
}

type spyTransport struct {
	transportError error
	requestsCount  int
}

func (t *spyTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	t.requestsCount++
	return nil, t.transportError
}

func intPtr(v int) *int {
	return &v
}

func durationPtr(v time.Duration) *time.Duration {
	return &v
}
