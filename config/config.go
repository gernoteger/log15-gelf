// Package config helps to add the gelf writer to log15-config
package config

import (
	"reflect"

	"github.com/gernoteger/log15-config"
	"github.com/gernoteger/log15-gelf"
	"github.com/gernoteger/mapstructure-hooks"
	"github.com/inconshreveable/log15"
)

func init() {
	Register()
}

// use for registry functions
var HandlerConfigType = reflect.TypeOf((*config.HandlerConfig)(nil)).Elem()

// registers all handlers
func Register() {
	hooks.Register(HandlerConfigType, "gelf", NewGelfConfig)
}

type GelfConfig struct {
	config.LevelHandlerConfig `mapstructure:",squash"`
	Address                   string
}

// make sure its's the right interface
var _ config.HandlerConfig = (*GelfConfig)(nil)

func NewGelfConfig() interface{} {
	return &GelfConfig{}
}

func (c *GelfConfig) NewHandler() (log15.Handler, error) {
	h, err := gelf.GelfHandler(c.Address)
	return h, err
}
