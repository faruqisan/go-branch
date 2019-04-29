package httpclient

import (
	"testing"
)

func TestDefaultHTTPClientDo(t *testing.T) {
	testHTTPClientDo(t, Do)
}

func TestDefaultHTTPClientDoJSON(t *testing.T) {
	testHTTPClientDoJSON(t, DoJSON)
}
