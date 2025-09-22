package main

import (
	"api-watchdog/internal/runner"
	"encoding/json"
	"os"

	"github.com/yourpov/logrite"
)

type GlobalCfg struct {
	DiscordWebhook string `json:"discord_webhook"`
	RequestTimeout int    `json:"request_timeout_ms"`
}

type CheckCfg struct {
	Name                 string `json:"name"`
	URL                  string `json:"url"`
	CheckIntervalSeconds int    `json:"check_interval_seconds"`
}

type Config struct {
	Global GlobalCfg  `json:"global"`
	Checks []CheckCfg `json:"checks"`
}

// loadConfig loads the config
func loadConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		os.Exit(1)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		os.Exit(1)
	}

	if cfg.Global.DiscordWebhook == "" {
		logrite.Warn("[Config/Global] Webhook not set")
		os.Exit(1)
	}

	for i := range cfg.Checks {
		if cfg.Checks[i].Name == "" {
			logrite.Warn("[Config/Checks] API Name not set")
			os.Exit(1)
		}

		if cfg.Checks[i].URL == "" {
			logrite.Warn("[Config/Checks] API Link not set")
			os.Exit(1)
		}
	}

	return cfg, nil
}

func main() {
	logrite.SetConfig(logrite.Config{
		ShowIcons:    true,
		UppercaseTag: true,
		UseColors:    true,
	})

	cfg, err := loadConfig("config/config.json")
	if err != nil {
		logrite.Error("[Config]: %v", err)
		os.Exit(1)
	}

	r := runner.New(runner.Options{
		Webhook:     cfg.Global.DiscordWebhook,
		RequestTOms: cfg.Global.RequestTimeout,
	})

	for _, c := range cfg.Checks {
		r.AddCheck(runner.Check{
			Name:      c.Name,
			URL:       c.URL,
			IntervalS: c.CheckIntervalSeconds,
		})
	}

	r.Run()

	select {}
}
