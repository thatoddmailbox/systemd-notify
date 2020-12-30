# systemd-notify
This is a small program that can complain to you when the status of a systemd unit changes. It listens to notifications from systemd over D-Bus, so you don't need to modify your service, and notifications are immediate.

Currently, this only supports posting to Microsoft Teams, but other notification methods may be added at a later point.