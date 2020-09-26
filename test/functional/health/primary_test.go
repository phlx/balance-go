package health

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

// TestHealth tests application health — postgres and redis are functional and latency is low
func TestHealth(t *testing.T) {
	carcass := functional.NewCarcass(t)

	json := carcass.Expectations.GET("/_health").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()

	json.ValueEqual("postgres", true)
	json.ValueEqual("redis", true)
	json.Value("latency_ms").Number().Le(100)
	json.Value("errors").Array().Empty()
}
