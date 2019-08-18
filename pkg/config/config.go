package config

import (
	"github.com/go-ini/ini"
)

func Load(source string, section string, v interface{}) {
	cfg, err := ini.Load(source)
	if err != nil {
		panic(err)
	}

	err = cfg.Section(section).MapTo(v)
	if err != nil {
		panic(err)
	}
}
