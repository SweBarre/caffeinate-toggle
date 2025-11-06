package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"strconv"

	"gopkg.in/ini.v1"
)

type Config struct {
	RunningTitle    string
	NotRunningTitle string
	CaffeinateOnStart bool
	Options         []string
	Timers          map[string]int
	TimerOrder       []string   // preserves menu order
}

func Load() *Config {
	cfg := &Config{
		RunningTitle:    "☕️",
		NotRunningTitle: "😴",
		Options:         []string{"-dims"},
		Timers:          map[string]int{"5m": 300, "30m": 1800, "1h": 3600},
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Println("Failed to get HOME dir:", err)
		return cfg
	}

	path := filepath.Join(home, ".caffeinate-toggle.conf")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Println("config file not found:", path)
		return cfg // config file not found, use defaults
	}

	iniCfg, err := ini.Load(path)
	if err != nil {
		log.Println("Failed to parse config:", err)
		return cfg
	}

	sec := iniCfg.Section("main")
	cfg.RunningTitle = sec.Key("running_title").MustString(cfg.RunningTitle)
	cfg.NotRunningTitle = sec.Key("not_running_title").MustString(cfg.NotRunningTitle)
	cfg.CaffeinateOnStart = sec.Key("caffeinate_on_start").MustBool(false)
	opts := sec.Key("options").String()
	if opts != "" {
		cfg.Options = splitArgs(opts)
	}

	timerSec := iniCfg.Section("Timers")
	for _, key := range timerSec.KeyStrings() {
		v, err := timerSec.Key(key).Int()
		if err == nil {
			cfg.Timers[key] = v
		}
	}
	log.Println("Loaded configuration file:", path)

	return cfg
}

func (cfg *Config) Save() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println("Unable to get home directory:", err)
		return
	}

	configPath := filepath.Join(home, ".caffeinate-toggle.conf")
	iniCfg := ini.Empty()

	// main section
	sec, _ := iniCfg.NewSection("main")
	sec.NewKey("running_title", cfg.RunningTitle)
	sec.NewKey("not_running_title", cfg.NotRunningTitle)
	sec.NewKey("options", strings.Join(cfg.Options, " "))
	sec.NewKey("caffeinate_on_start", strconv.FormatBool(cfg.CaffeinateOnStart))

	timerSec, _ := iniCfg.NewSection("Timers")
	for label, seconds := range cfg.Timers {
		timerSec.NewKey(label, strconv.Itoa(seconds))
	}

	if err := iniCfg.SaveTo(configPath); err != nil {
		log.Println("Failed to save config:", err)
	} else {
		log.Println("Config saved to", configPath)
	}
}


func splitArgs(s string) []string {
	return strings.Fields(s)
}

