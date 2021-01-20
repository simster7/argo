package apiserver

import (
	"github.com/simster7/argo/v2/config"
	"github.com/simster7/argo/v2/server/auth/sso"
)

var emptyConfigFunc = func() interface{} { return &Config{} }

type Config struct {
	config.Config
	// SSO in settings for single-sign on
	SSO sso.Config `json:"sso,omitempty"`
}
