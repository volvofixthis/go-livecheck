package validator

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"net"
	"strings"
	"time"

	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
)

type L4ValidatorExtra struct {
	Timeout string `mapstructure:"timeout,omitempty"`
	Proto   string `mapstructure:"proto"`
	Data    string `mapstructure:"data,omitempty"`
	Send    bool   `mapstructure:"send,omitempty"`
}

type L4Config struct {
	proto   string
	timeout time.Duration
	data    string
	send    bool
}

type UDPValidator struct {
	*Validator
	l4Config L4Config
}

func (v *UDPValidator) Exec(data map[string]any) (bool, error) {
	addresses := strings.Split(v.Validator.config.Rule, ",")
	for _, address := range addresses {
		s, err := net.ResolveUDPAddr(v.l4Config.proto, address)
		if err != nil {
			return false, err
		}
		c, err := net.DialUDP(v.l4Config.proto, nil, s)
		if err != nil {
			return false, err
		}
		_, err = c.Write([]byte(v.l4Config.data))
		if err != nil {
			return false, err
		}
		time.Sleep(v.l4Config.timeout)
		_, err = c.Write([]byte(v.l4Config.data))
		if err != nil {
			return false, err
		}
	}
	if len(addresses) == 0 {
		return false, errors.New("no addresses")
	}
	return true, nil
}

type TCPValidator struct {
	*Validator
	l4Config L4Config
}

func (v *TCPValidator) Exec(data map[string]any) (bool, error) {
	addresses := strings.Split(v.Validator.config.Rule, ",")
	for _, address := range addresses {
		c, err := net.DialTimeout(v.l4Config.proto, address, v.l4Config.timeout)
		if err != nil {
			return false, nil
		}
		if v.l4Config.send {
			_, err = c.Write([]byte(v.l4Config.data))
			if err != nil {
				return false, nil
			}
		}
		if err := c.Close(); err != nil {
			return false, err
		}
	}
	return true, nil
}

func NewL4Validator(c *config.ValidatorConfig) (ValidatorInterface, error) {
	extra := L4ValidatorExtra{}
	if err := mapstructure.Decode(c.Extra, &extra); err != nil {
		return nil, err
	}
	l4Config := L4Config{
		proto: extra.Proto,
		send:  extra.Send,
		data:  "Bite my shiny metal ass! (c) Bender",
	}
	if extra.Data != "" {
		l4Config.data = extra.Data
	}
	if extra.Timeout == "" {
		l4Config.timeout = time.Second
	} else {
		d, err := time.ParseDuration(extra.Timeout)
		if err != nil {
			return nil, err
		}
		l4Config.timeout = d
	}
	switch extra.Proto {
	case "tcp", "tcp6":
		return &TCPValidator{
			Validator: NewValidatorBase(c),
			l4Config:  l4Config,
		}, nil
	case "udp", "udp6":
		return &UDPValidator{
			Validator: NewValidatorBase(c),
			l4Config:  l4Config,
		}, nil
	}

	return nil, errors.New("No such L4 validator")
}
