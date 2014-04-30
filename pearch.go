package main

import (
	"database/sql"
	"encoding/xml"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				go doTheDirtyWork()
			case <-sigchan:
				ticker.Stop()
				return
			}
		}
	}()
	<-sigchan
}

func doTheDirtyWork() {
	response, err := http.Get("http://map.pilotedge.net/status_live.xml")
	defer response.Body.Close()
	if err != nil {
		panic(err)
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	for i, ch := range contents {
		switch {
		case ch > '~':
			contents[i] = ' '
		case ch == '\r':
		case ch == '\n':
		case ch == '\t':
		case ch < ' ':
			contents[i] = ' '
		}
	}

	status := Status{}
	err = xml.Unmarshal([]byte(string(contents)), &status)
	if err != nil {
		panic(err)
	}

	connectionString := "user=pearch password=flylikeaneagle dbname=pearch sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	t := time.Now().UTC()

	for _, c := range status.Controllers {
		sql := "insert into controller_events (event_timestamp, name, role, callsign, primary_frequency, latitude, longitude) values ($1, $2, $3, $4, $5, $6, $7)"
		_, err = db.Exec(sql, t, c.Name, c.Role, c.Callsign, c.PrimaryFrequency, c.Position.Latitude, c.Position.Longitude)
		if err != nil {
			panic(err)
		}
	}

	for _, p := range status.Pilots {
		sql := "insert into pilot_events (event_timestamp, cid, name, equipment, callsign, frequency, radio, desired_role, latitude, longitude, altitude, ground_speed, true_heading, flight_plan_origin, flight_plan_destination, flight_plan_route, flight_plan_remarks) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)"
		_, err = db.Exec(sql, t, p.CID, p.Name, p.Equipment, p.Callsign, p.Frequency.Freq, p.Frequency.Radio, p.Frequency.DesiredRole, p.Position.Latitude, p.Position.Longitude, p.Position.Altitude, p.Position.GroundSpeed, p.Position.TrueHeading, p.FlightPlan.Origin, p.FlightPlan.Destination, p.FlightPlan.Route, p.FlightPlan.Remarks)
		if err != nil {
			panic(err)
		}
	}
}

type Status struct {
	Controllers []Controller `xml:"controllers>controller"`
	Pilots      []Pilot      `xml:"pilots>pilot"`
}

type Controller struct {
	Name             string   `xml:"name"`
	Role             string   `xml:"role"`
	Position         Position `xml:"position"`
	PrimaryFrequency int32    `xml:"primFreq"`
	Callsign         string   `xml:"callsign"`
}

type Pilot struct {
	CID        int32      `xml:"cid,attr"`
	Name       string     `xml:"name"`
	Equipment  string     `xml:"equipment"`
	Callsign   string     `xml:"callsign"`
	Position   Position   `xml:"position"`
	Frequency  Frequency  `xml:"frequency"`
	FlightPlan FlightPlan `xml:"flightplan"`
}

type FlightPlan struct {
	Origin      string `xml:"origin,attr"`
	Destination string `xml:"destination,attr"`
	Route       string `xml:"route"`
	Remarks     string `xml:"remarks"`
}

type Frequency struct {
	Radio       string `xml:"radio,attr"`
	DesiredRole string `xml:"desiredRole,attr"`
	Freq        string `xml:",chardata"`
}

type Position struct {
	Latitude    float64 `xml:"lat,attr"`
	Longitude   float64 `xml:"lon,attr"`
	Altitude    float64 `xml:"alt,attr"`
	GroundSpeed float64 `xml:"groundSpeed,attr"`
	TrueHeading float64 `xml:"trueHeading,attr"`
}
