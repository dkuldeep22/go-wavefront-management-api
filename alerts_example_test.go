package wavefront_test

import (
	"fmt"
	"github.com/WavefrontHQ/go-wavefront-management-api"
	"io/ioutil"
	"log"
)

func ExampleAlerts() {
	config := &wavefront.Config{
		Address: "test.wavefront.com",
		Token:   "xxxx-xxxx-xxxx-xxxx-xxxx",
	}
	client, err := wavefront.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	alerts := client.Alerts()
	targets := client.Targets()

	// ######################
	// EXAMPLE: Classic Alert
	// ######################

	a := &wavefront.Alert{
		Name:                "My First Alert",
		Target:              "test@example.com",
		Condition:           "ts(servers.cpu.usage, dc=dc2) > 10 * 10",
		DisplayExpression:   "ts(servers.cpu.usage, dc=dc2)",
		Minutes:             2,
		ResolveAfterMinutes: 2,
		Severity:            "WARN",
		Tags:                []string{"dc1", "synergy"},
	}

	// Create the alert on Wavefront
	err = alerts.Create(a)
	if err != nil {
		log.Fatal(err)
	}

	// The ID field is now set, so we can update/delete the Alert
	fmt.Println("alert ID is", *a.ID)

	// Alternatively we could search for the Alert
	err = alerts.Get(&wavefront.Alert{
		ID: a.ID,
	})

	if err != nil {
		log.Fatal(err)
	}

	// Update the Alert
	a.Target = "test@example.com,bob@example.com"
	err = alerts.Update(a)
	if err != nil {
		log.Fatal(err)
	}

	// Delete the Alert
	err = alerts.Delete(a, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("alert deleted")

	// ########################
	// EXAMPLE: Threshold Alert
	// ########################

	// Threshold Alerts only accept custom alert targets
	// Create Alert Targets that can be used for the example

	tmpl, _ := ioutil.ReadFile("./target-template.tmpl")
	targetA := wavefront.Target{
		Title:       "test target",
		Description: "testing something A",
		Method:      "WEBHOOK",
		Recipient:   "https://hooks.slack.com/services/test/me",
		ContentType: "application/json",
		CustomHeaders: map[string]string{
			"Testing": "true",
		},
		Triggers: []string{"ALERT_OPENED", "ALERT_RESOLVED"},
		Template: string(tmpl),
	}

	targetB := wavefront.Target{
		Title:       "test target",
		Description: "testing something B",
		Method:      "WEBHOOK",
		Recipient:   "https://hooks.slack.com/services/test/me",
		ContentType: "application/json",
		CustomHeaders: map[string]string{
			"Testing": "true",
		},
		Triggers: []string{"ALERT_OPENED", "ALERT_RESOLVED"},
		Template: string(tmpl),
	}

	targetC := wavefront.Target{
		Title:       "test target",
		Description: "testing something C",
		Method:      "WEBHOOK",
		Recipient:   "https://hooks.slack.com/services/test/me",
		ContentType: "application/json",
		CustomHeaders: map[string]string{
			"Testing": "true",
		},
		Triggers: []string{"ALERT_OPENED", "ALERT_RESOLVED"},
		Template: string(tmpl),
	}

	// Create the targets on Wavefront
	err = targets.Create(&targetA)
	if err != nil {
		log.Fatal(err)
	}
	err = targets.Create(&targetB)
	if err != nil {
		log.Fatal(err)
	}
	err = targets.Create(&targetC)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("alert targets created")

	strA := fmt.Sprintf("target:%s", *targetA.ID)
	strB := fmt.Sprintf("target:%s", *targetB.ID)
	strC := fmt.Sprintf("target:%s", *targetC.ID)

	mta := &wavefront.Alert{
		Name:      "My First Threshold Alert",
		AlertType: "THRESHOLD",
		Targets: map[string]string{
			"smoke": strA,
			"warn":  strB,
		},
		Conditions: map[string]string{
			"smoke": "ts(servers.cpu.usage) > 70",
			"warn":  "ts(servers.cpu.usage) > 100",
		},
		DisplayExpression:   "ts(servers.cpu.usage)",
		Minutes:             2,
		ResolveAfterMinutes: 2,
		SeverityList:        []string{"SMOKE", "WARN"},
		Tags:                []string{"dc1", "synergy"},
	}

	// Create the thresold alert on Wavefront
	err = alerts.Create(mta)
	if err != nil {
		log.Fatal(err)
	}

	// The ID field is now set, so we can update/delete the Threshold Alert
	fmt.Println("threshold alert ID is", *mta.ID)

	// Alternatively we could search for the Threshold Alert
	err = alerts.Get(&wavefront.Alert{
		ID: mta.ID,
	})

	if err != nil {
		log.Fatal(err)
	}

	// Update the Threshold Alert
	mta.Targets["smoke"] = strC
	err = alerts.Update(mta)
	if err != nil {
		log.Fatal(err)
	}

	// Delete the Threshold Alert
	err = alerts.Delete(mta, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("threshold alert deleted")

	err = targets.Delete(&targetA)
	if err != nil {
		log.Fatal(err)
	}
	err = targets.Delete(&targetB)
	if err != nil {
		log.Fatal(err)
	}
	err = targets.Delete(&targetC)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("alert targets deleted")
}
