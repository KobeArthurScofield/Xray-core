package xhttp_test

import (
	"testing"

	. "github.com/xtls/xray-core/transport/internet/xhttp"
)

func Test_GetNormalizedPath(t *testing.T) {
	c := Config{
		Path: "/?world",
	}

	path := c.GetNormalizedPath()
	if path != "/" {
		t.Error("Unexpected: ", path)
	}
}
