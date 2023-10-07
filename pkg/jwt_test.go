package pkg

import "testing"

func TestValidateJwtAndExtractSubject(t *testing.T) {
	subject, err := ValidateJwtAndExtractSubject(
		"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ2b2RvbGF6MDk1IiwiaWF0IjoxNTE2MjM5MDIzLCJleHAiOjE4MTYyMzkwMjN9.ga-y5GLMIJaoO6-UPY3eTztUQjGNj1QD_7yS0oggse7MKoyUaEZDMcO8ADRm7m6F1oWvSWBpu6hoPfLCQ64Emg",
		"super_secret_for_purser",
	)
	if err != nil {
		t.Errorf("error validating token: %s", err)
	} else {
		t.Logf("Subject: %s", subject)
	}
}
