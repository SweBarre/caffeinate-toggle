package tray

import (
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/SweBarre/caffeinate-toggle/internal/caffeinate"
	"github.com/SweBarre/caffeinate-toggle/internal/config"
	"github.com/getlantern/systray"
)

var cfg *config.Config

func Run() {
	cfg = config.Load()
	if cfg.CaffeinateOnStart && !caffeinate.IsRunning() {
		args := strings.Join(cfg.Options, " ")
		err := caffeinate.Start(args)
		if err != nil {
			log.Fatalln("Failed to start caffeinate!")
		}
	}
	systray.Run(onReady, onExit)
}

func onReady() {
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
				args := strings.Join(cfg.Options, " ")
				err := caffeinate.StartTimed(args, sec)
				if err != nil {
					log.Fatalln("Failed to start caffeinate!")
				}
				systray.SetTitle(cfg.RunningTitle)
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
					log.Println("Stopping....")
					caffeinate.Stop()
					systray.SetTitle(cfg.NotRunningTitle)
				} else {
					log.Println("Starting")
					args := strings.Join(cfg.Options, " ")
					err := caffeinate.Start(args)
					if err != nil {
						log.Fatalln("Failed to start caffeinate!")
					}
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
