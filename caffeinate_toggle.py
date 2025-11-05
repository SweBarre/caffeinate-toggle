#!/usr/bin/env python3
import subprocess
import rumps

class CaffeinateApp(rumps.App):
    def __init__(self):
        super(CaffeinateApp, self).__init__("☕️", quit_button=None)
        self.menu = ["Toggle Caffeinate", None, "Quit"]
        self.update_indicator()

        # Auto-refresh every 5 seconds
        self.timer = rumps.Timer(self.update_indicator, 5)
        self.timer.start()

    def is_caffeinate_running(self):
        """Return True if caffeinate process is active."""
        try:
            subprocess.check_output(["pgrep", "-x", "caffeinate"])
            return True
        except subprocess.CalledProcessError:
            return False

    def update_indicator(self, _=None):
        """Update menu bar icon based on caffeinate state."""
        if self.is_caffeinate_running():
            self.title = "☕️"
        else:
            self.title = "😴"

    @rumps.clicked("Toggle Caffeinate")
    def toggle_caffeinate(self, _):
        if self.is_caffeinate_running():
            subprocess.run(["killall", "caffeinate"])
        else:
            subprocess.Popen(["caffeinate", "-dims"])
        self.update_indicator()  # update immediately

    @rumps.clicked("Quit")
    def quit_app(self, _):
        if self.is_caffeinate_running():
            subprocess.run(["killall", "caffeinate"])
        rumps.quit_application()

if __name__ == "__main__":
    CaffeinateApp().run()
