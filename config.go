package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	From     string `toml:"from"`

	Protocol string `toml:"protocol"`

	SMTPHost string `toml:"smtp_host"`
	SMTPPort int    `toml:"smtp_port"`
	UseTLS   bool   `toml:"use_tls"`

	POP3Host string `toml:"pop3_host"`
	POP3Port int    `toml:"pop3_port"`

	IMAPHost string `toml:"imap_host"`
	IMAPPort int    `toml:"imap_port"`

	LogFile string `toml:"log_file"`
	Debug   bool   `toml:"debug"`
	Timeout int    `toml:"timeout"`
	Retry   int    `toml:"retry"`
}

func defaultConfigPath() string {
	if env := os.Getenv("MAILC_CONFIG"); env != "" {
		return env
	}

	home := os.Getenv("HOME")
	if home == "" {
		return ""
	}

	return filepath.Join(home, ".config/mailc/config.toml")
}

func ensureConfigDir() error {
	path := filepath.Dir(defaultConfigPath())
	return os.MkdirAll(path, 0700)
}

func LoadConfig() (*Config, error) {
	path := defaultConfigPath()

	data, err := os.ReadFile(path)
	if err != nil {
		// 自动创建默认配置
		if os.IsNotExist(err) {
			if err := ensureConfigDir(); err != nil {
				return nil, err
			}

			defaultCfg := Config{
				Username: "your@email.com",
				Password: "your_app_password",
				From:     "",
				Protocol: "smtp",

				SMTPHost: "smtp.gmail.com",
				SMTPPort: 587,
				UseTLS:   true,

				LogFile: filepath.Join(os.Getenv("HOME"), ".local/share/mailc/mailc.log"),
				Debug:   false,
				Timeout: 10,
				Retry:   3,
			}

			f, _ := os.Create(path)
			defer f.Close()

			enc := toml.NewEncoder(f)
			_ = enc.Encode(defaultCfg)

			return &defaultCfg, nil
		}
		return nil, err
	}

	var cfg Config
	err = toml.Unmarshal(data, &cfg)
	return &cfg, err
}

func (c *Config) Validate() error {
	if c.Username == "" || c.Password == "" {
		return fmt.Errorf("missing username/password in config")
	}
	return nil
}
