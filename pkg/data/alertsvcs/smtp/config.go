package smtp

import (
	"errors"
	"fmt"
	"strings"
)

// Config configuration data
type Config struct {
	Enabled  bool   `toml:"enabled" override:"enabled"`
	Host     string `toml:"host" override:"host"`
	Port     int    `toml:"port" override:"port"`
	Username string `toml:"username" override:"username"`
	Password string `toml:"password" override:"password,redact"`
	// Whether to skip TLS verify.
	InsecureSkipVerify bool `toml:"insecure-skip-verify" override:"insecure-skip-verify"`
	// From address
	From string `toml:"from" override:"from"`
	// Default To addresses
	To []string `toml:"to" override:"to"`
}

// Validate basic validations
func (c Config) Validate() error {
	if c.Host == "" {
		return errors.New("host cannot be empty")
	}
	if c.Port <= 0 {
		return fmt.Errorf("invalid port %d", c.Port)
	}
	if c.Enabled && c.From == "" {
		return errors.New("must provide a 'from' address")
	}
	// Poor mans email validation, but since emails have a very large domain this is probably good enough
	// to catch user error.
	if c.From != "" && !strings.ContainsRune(c.From, '@') {
		return fmt.Errorf("invalid from email address: %q", c.From)
	}
	for _, t := range c.To {
		if !strings.ContainsRune(t, '@') {
			return fmt.Errorf("invalid to email address: %q", t)
		}
	}
	return nil
}
