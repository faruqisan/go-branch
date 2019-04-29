package branch

import (
	"github.com/faruqisan/go-branch/httpclient"
)

type (
	// Engine struct define requirement and act as function receiver
	// for the whole library
	Engine struct {
		client *httpclient.Client
	}
)

// New function return setuped engine
func New() *Engine {
	return &Engine{
		client: httpclient.NewClient(),
	}
}
