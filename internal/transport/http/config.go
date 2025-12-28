package http //revive:disable-line:var-naming

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// errParsingAdress indicates an error while parsing an address URL for an instance of http server
var errParsingAdress = errors.New("error parsing address")

// Config contains fields and [flag.Value]s to set up the [Server]
type Config struct {
	address        address
	accrualAddress address
}

// Address returns a pointer to the [flag.Value] to set up the [Server]
func (c *Config) Address() *address { // revive:disable-line:unexported-return provides the interface to the caller
	return &c.address
}

// AccrualAddress returns a pointer to the [flag.Value] to set up the accrual [Client]
func (c *Config) AccrualAddress() *address { // revive:disable-line:unexported-return provides the interface to the caller
	return &c.accrualAddress
}

// String returns [codings] elements separated by ", " as a single string
type address struct {
	host   string
	port   string
	Source string
}

func (a address) String() string {
	return a.host + ":" + a.port
}

func (a *address) Set(s string) error {
	if s == "" {
		return fmt.Errorf("%w: %s", errParsingAdress, "empty string")
	}

	if strings.HasPrefix(s, "localhost:") {
		s = "http://" + s
	}
	url, err := url.Parse(s)
	if err != nil {
		return err
	}

	if a.host = url.Hostname(); url.Hostname() == "" {
		return fmt.Errorf("%w: %s", errParsingAdress, "empty scheme")
	}

	if a.port = url.Port(); url.Port() == "" {
		return fmt.Errorf("%w: %s", errParsingAdress, "empty port")
	}

	return nil
}

type secret string

func (sec secret) String() string {
	return ""
}

func (sec *secret) Set(s string) error {
	*sec = secret(s)
	return nil
}
