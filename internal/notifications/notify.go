package notifications

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Notifier defines the interface for sending notifications
type Notifier interface {
	Send(title, message string) error
}

// SystemNotifier sends notifications using OS-specific commands
type SystemNotifier struct {
	fallbackToConsole bool
}

// NewSystemNotifier creates a new system notifier
func NewSystemNotifier(fallbackToConsole bool) *SystemNotifier {
	return &SystemNotifier{
		fallbackToConsole: fallbackToConsole,
	}
}

// Send sends a desktop notification
func (n *SystemNotifier) Send(title, message string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		// Use notify-send on Linux
		cmd = exec.Command("notify-send", title, message, "-u", "normal")

	case "darwin":
		// Use osascript on macOS
		script := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
		cmd = exec.Command("osascript", "-e", script)

	case "windows":
		// Use PowerShell on Windows
		script := fmt.Sprintf(`
			[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
			[Windows.UI.Notifications.ToastNotification, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
			[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null

			$template = @"
			<toast>
				<visual>
					<binding template="ToastText02">
						<text id="1">%s</text>
						<text id="2">%s</text>
					</binding>
				</visual>
			</toast>
"@

			$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
			$xml.LoadXml($template)
			$toast = New-Object Windows.UI.Notifications.ToastNotification $xml
			[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("godo").Show($toast)
		`, title, message)
		cmd = exec.Command("powershell", "-Command", script)

	default:
		if n.fallbackToConsole {
			return n.fallbackConsole(title, message)
		}
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	err := cmd.Run()
	if err != nil && n.fallbackToConsole {
		return n.fallbackConsole(title, message)
	}

	return err
}

// fallbackConsole prints notification to console when system notifications fail
func (n *SystemNotifier) fallbackConsole(title, message string) error {
	fmt.Printf("\n🔔 %s\n%s\n\n", title, message)
	return nil
}

// ConsoleNotifier always prints to console (useful for testing or when no GUI)
type ConsoleNotifier struct{}

// NewConsoleNotifier creates a new console-only notifier
func NewConsoleNotifier() *ConsoleNotifier {
	return &ConsoleNotifier{}
}

// Send prints the notification to console
func (n *ConsoleNotifier) Send(title, message string) error {
	fmt.Printf("\n🔔 %s\n%s\n\n", title, message)
	return nil
}

// NoOpNotifier does nothing (useful for disabling notifications)
type NoOpNotifier struct{}

// NewNoOpNotifier creates a notifier that does nothing
func NewNoOpNotifier() *NoOpNotifier {
	return &NoOpNotifier{}
}

// Send does nothing
func (n *NoOpNotifier) Send(title, message string) error {
	return nil
}

