# systemd-notify
This is a small program that can complain to you when the status of a systemd unit changes. It listens to notifications from systemd over D-Bus, so you don't need to modify your service, and notifications are immediate.

It currently supports posting notifications to Slack.

> There is also legacy support for Microsoft Teams, but it relies on the deprecated webhook functionality. You will probably need to update it to use the new Power Automate flow to get it working.