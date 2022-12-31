package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ddritzenhoff/dinny"
	rest "github.com/ddritzenhoff/dinny/http"
)

var mondaySlackUID string
var tuesdaySlackUID string
var wednesdaySlackUID string
var thursdaySlackUID string
var fridaySlackUID string
var saturdaySlackUID string
var sundaySlackUID string

// AssignCooksCommand is a command to assign cooks.
type AssignCooksCommand struct {
	ConfigPath string
}

// Run executes the assign_cooks command.
func (c *AssignCooksCommand) Run(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&mondaySlackUID, "monday", "", "set cook for monday <slackUID>")
	fs.StringVar(&tuesdaySlackUID, "tuesday", "", "set cook for tuesday <slackUID>")
	fs.StringVar(&wednesdaySlackUID, "wednesday", "", "set cook for wednesday <slackUID>")
	fs.StringVar(&thursdaySlackUID, "thursday", "", "set cook for thursday <slackUID>")
	fs.StringVar(&fridaySlackUID, "friday", "", "set cook for friday <slackUID>")
	fs.StringVar(&saturdaySlackUID, "saturday", "", "set cook for saturday <slackUID>")
	fs.StringVar(&sundaySlackUID, "sunday", "", "set cook for sunday <slackUID>")
	attachConfigFlags(fs, &c.ConfigPath)
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("Run fs.Parse: %w", err)
	}

	// Load the configuration.
	config, err := ReadConfigFile(c.ConfigPath)
	if err != nil {
		return fmt.Errorf("Run ReadConfigFile: %w", err)
	}

	now := time.Now()
	var request rest.AssignCooksRequest

	if sundaySlackUID != "" {
		request.CookAssignments = append(request.CookAssignments, buildCookAssignment(now, time.Sunday, sundaySlackUID))
	}
	if mondaySlackUID != "" {
		request.CookAssignments = append(request.CookAssignments, buildCookAssignment(now, time.Monday, mondaySlackUID))
	}
	if tuesdaySlackUID != "" {
		request.CookAssignments = append(request.CookAssignments, buildCookAssignment(now, time.Tuesday, tuesdaySlackUID))
	}
	if wednesdaySlackUID != "" {
		request.CookAssignments = append(request.CookAssignments, buildCookAssignment(now, time.Wednesday, wednesdaySlackUID))
	}
	if thursdaySlackUID != "" {
		request.CookAssignments = append(request.CookAssignments, buildCookAssignment(now, time.Thursday, thursdaySlackUID))
	}
	if fridaySlackUID != "" {
		request.CookAssignments = append(request.CookAssignments, buildCookAssignment(now, time.Friday, fridaySlackUID))
	}
	if saturdaySlackUID != "" {
		request.CookAssignments = append(request.CookAssignments, buildCookAssignment(now, time.Saturday, saturdaySlackUID))
	}

	if len(request.CookAssignments) == 0 {
		return fmt.Errorf("Run: didn't specify any days, so nothing happened")
	}

	buf, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("Run json.Marshal: %w", err)
	}
	responseBody := bytes.NewBuffer(buf)

	url := fmt.Sprintf("%s/cmd/assign-cooks", config.URL)
	resp, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		return fmt.Errorf("Run http.Post: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Run io.ReadAll: %w", err)
	}
	b, err := prettyPrint(body)
	if err != nil {
		return fmt.Errorf("Run prettyPrint: %w", err)
	}
	fmt.Println(string(b))
	return nil
}

// getDayDifference returns the number of days between now (today) and then (some other weekday).
// This function will return some value between 0 and 6. 0 --> now=Monday, then=Monday 6 --> now=Wednesday, then=Tuesday.
func getDayDifference(now time.Weekday, then time.Weekday) int {
	return (int(then) - int(now) + 7) % 7
}

// buildCookAssignment creates a pb.CookAssignment.
func buildCookAssignment(now time.Time, then time.Weekday, cookSlackUID string) rest.CookAssignment {
	year, month, day := now.AddDate(0, 0, getDayDifference(now.Weekday(), then)).Date()
	return rest.CookAssignment{
		Date: dinny.Date{
			Year:  year,
			Month: month,
			Day:   day,
		},
		CookSlackUID: cookSlackUID,
	}
}

// usage prints usage information for assign_cooks to STDOUT.
func (c *AssignCooksCommand) usage() {
	fmt.Println(`
Assign cooks for the next week starting with the current day.
If today were Monday, and a cook were to be assigned with the -monday flag, that cook would be set to cook today.

Usage:

		dinny assign_cooks -monday <slackUID> -tuesday <slackUID> -wednesday <slackUID> -thursday <slackUID> -friday <slackUID> -saturday <slackUID> -sunday <slackUID>

Arguments:

		-monday <slackUID>
			Set the cook for the upcoming Monday
		-tuesday <slackUID>
			Set the cook for the upcoming Tuesday
		-wednesday <slackUID>
			Set the cook for the upcoming Wednesday
		-thursday <slackUID>
			Set the cook for the upcoming Thursday
		-friday <slackUID>
			Set the cook for the upcoming Friday
		-saturday <slackUID>
			Set the cook for the upcoming Saturday
		-sunday <slackUID>
			Set the cook for the upcoming Sunday
	`[1:])
}
