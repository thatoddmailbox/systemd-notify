package main

import (
	"log"
	"strings"

	"github.com/godbus/dbus/v5"
)

const (
	systemdDestination = "org.freedesktop.systemd1"
)

func main() {
	log.Println("systemd-notify")

	err := loadConfig()
	if err != nil {
		panic(err)
	}

	conn, err := dbus.SystemBus()
	if err != nil {
		panic(err)
	}

	systemdObj := conn.Object(systemdDestination, "/org/freedesktop/systemd1")

	unitObjects := []dbus.BusObject{}

	pathToName := map[dbus.ObjectPath]string{}

	// get all the units we're supposed to watch
	for _, unit := range currentConfig.Watch.Units {
		if !strings.HasSuffix(unit, ".service") {
			log.Printf("WARNING: Unit '%s' specified in config file is not a service. This is not supported.", unit)
		}
		call := systemdObj.Call("org.freedesktop.systemd1.Manager.GetUnit", 0, unit)
		if call.Err != nil {
			// this is a very ugly way of doing this
			// however, this is really only to catch a misconfiguration, so it should be ok?
			if strings.HasSuffix(strings.TrimSpace(call.Err.Error()), "not loaded.") {
				log.Fatalf("Unit '%s' specified in config file does not exist or is not loaded.", unit)
			}

			panic(call.Err)
		}

		unitPath := call.Body[0].(dbus.ObjectPath)
		unitObject := conn.Object(systemdDestination, unitPath)

		pathToName[unitPath] = unit

		unitObjects = append(unitObjects, unitObject)
	}

	for _, unitObject := range unitObjects {
		err = conn.AddMatchSignal(
			dbus.WithMatchObjectPath(unitObject.Path()),
			dbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
			dbus.WithMatchMember("PropertiesChanged"),
		)
		if err != nil {
			panic(err)
		}
	}

	signal := make(chan *dbus.Signal, 10)
	conn.Signal(signal)
	for {
		select {
		case msg := <-signal:
			body := msg.Body

			if len(body) != 3 {
				// something we don't know, ignore it
				continue
			}

			if body[0] != "org.freedesktop.systemd1.Unit" {
				// something we don't know, ignore it
				continue
			}

			variants, ok := body[1].(map[string]dbus.Variant)
			if !ok {
				// something we don't know, ignore it
				continue
			}

			// do we care about this object?
			unitName, careAboutUnit := pathToName[msg.Path]
			if !careAboutUnit {
				// apparently not
				continue
			}

			// can we find the state variables?
			activeState, foundActiveState := variants["ActiveState"]
			subState, foundSubState := variants["SubState"]

			if !foundActiveState || !foundSubState {
				// something else must have changed
				continue
			}

			activeStateString := activeState.Value().(string)
			subStateString := subState.Value().(string)

			if len(currentConfig.Watch.FilterActiveStates) != 0 {
				// we have a filter enabled, apply it
				stateInFilter := false
				for _, filteredState := range currentConfig.Watch.FilterActiveStates {
					if filteredState == activeStateString {
						stateInFilter = true
						break
					}
				}

				if !stateInFilter {
					// we don't care about this event, ignore it
					continue
				}
			}

			notify(unitName, activeStateString, subStateString)
		}
	}
}
