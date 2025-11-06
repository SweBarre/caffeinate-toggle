package tray

import (
	"os"
	"time"
	"sort"
	"log"

	"github.com/SweBarre/caffeinate-toggle/internal/caffeinate"
	"github.com/SweBarre/caffeinate-toggle/internal/config"
	"github.com/getlantern/systray"
)

var cfg *config.Config

func Run() {
	cfg = config.Load()
	caffeinate.CaffeinateOptions = cfg.Options
	if cfg.CaffeinateOnStart && !caffeinate.IsRunning() {
        caffeinate.Start(cfg.Options)
    }
	systray.Run(onReady, onExit)
}

func onReady() {
	//cfg = config.Load()
	caffeinate.CaffeinateOptions = cfg.Options
	// Prepare a slice of labels
	timers := make([]string, 0, len(cfg.Timers))
	for label := range cfg.Timers {
		timers = append(timers, label)
	}

	// Sort labels by duration
	sort.Slice(timers, func(i, j int) bool {
		return cfg.Timers[timers[i]] < cfg.Timers[timers[j]]
	})
	systray.SetTitle(cfg.NotRunningTitle)
	systray.SetTooltip("Caffeinate Toggle")

	mToggle := systray.AddMenuItem("Toggle Caffeinate", "Start or stop caffeinate")

	mSaveConfig := systray.AddMenuItem("Save Config", "Save current configuration")
	systray.AddMenuItem("Timers:", "Available timed caffeinate options (non-clickable)")
	systray.AddSeparator()
	for _, label := range timers {
		seconds := cfg.Timers[label]
		item := systray.AddMenuItem(label, "Run caffeinate for "+label)
		go func(sec int, m *systray.MenuItem) {
			for range m.ClickedCh {
				caffeinate.StartTimed(sec, cfg.Options)
			}
		}(seconds, item)
	}


	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit Caffeinate Toggle")

	// auto-refresh
	go func() {
		for {
			if caffeinate.IsRunning() {
				systray.SetTitle(cfg.RunningTitle)
			} else {
				systray.SetTitle(cfg.NotRunningTitle)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// handle clicks
	go func() {
		for {
			select {
			case <-mToggle.ClickedCh:
				if caffeinate.IsRunning() {
					caffeinate.Stop()
					systray.SetTitle(cfg.NotRunningTitle)
				} else {
					caffeinate.Start(cfg.Options)
					systray.SetTitle(cfg.RunningTitle)
				}

			case <-mSaveConfig.ClickedCh:
				cfg.Save()

			case <-mQuit.ClickedCh:
				if caffeinate.IsRunning() {
					caffeinate.Stop()
				}
				systray.Quit()
				os.Exit(0)
			}
		}
	}()
}



func onExit() {
	if caffeinate.IsRunning() {
		caffeinate.Stop()
	}
	log.Println("Stopping CaffeniateToggle")
}
