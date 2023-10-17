package jwt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateJwtAndExtractSubject(t *testing.T) {
	type testCase struct {
		desc       string
		raw        string
		hmacSecret string
		subject    string
		err        error
	}
	var testCases = []testCase{
		{
			"valid",
			"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ2b2RvbGF6MDk1IiwiaWF0IjoxNTE2MjM5MDIzLCJleHAiOjE4MTYyMzkwMjN9.ga-y5GLMIJaoO6-UPY3eTztUQjGNj1QD_7yS0oggse7MKoyUaEZDMcO8ADRm7m6F1oWvSWBpu6hoPfLCQ64Emg",
			"super_secret_for_purser",
			"vodolaz095",
			nil,
		},
		{
			"invalid",
			"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ2b2RvbGF6MDk1IiwiaWF0IjoxNTE2MjM5MDIzLCJleHAiOjE4MTYyMzkwMjN9.ga-y5GLMIJaoO6-UPY3eTztUQjGNj1QD_7yS0oggse7MKoyUaEZDMcO8ADRm7m6F1oWvSWBpu6hoPfLCQ64Emg",
			"random_wrong_secret",
			"vodolaz095",
			fmt.Errorf("token signature is invalid: signature is invalid"),
		},
		{
			"expired",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ2b2RvbGF6MDk1IiwiZXhwIjoxNTE2MjMyMDIyfQ.Sb3uh1m8K-8gmTsjWUEiDLqRhL9g12tRsGR0YgJlh1E",
			"super_secret_for_purser",
			"vodolaz095",
			fmt.Errorf("token has invalid claims: token is expired"),
		},
		{
			"no issue time",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhIjoiYiJ9.Xl4yNieJTduaXjj2w1FgX9PklvcCtMcG3wGT8qH5g3k",
			"super_secret_for_purser",
			"vodolaz095",
			fmt.Errorf("token issue time is unknown"),
		},
		{
			"no subject",
			"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjMsImV4cCI6MTgxNjIzOTAyM30.IdqvOC10qDolHDzM6vOKFUKpQuSgX6FnB9e5IGWbTqGtCYXVFYW5sg7qXwP_PnyhJHLJZdP6oqkkyFbufEr5NA",
			"super_secret_for_purser",
			"", // but this is wrong!
			nil,
		},
	}
	for i := range testCases {
		t.Logf("Executing test case %v %s", i, testCases[i].desc)
		subject, err := ValidateJwtAndExtractSubject(testCases[i].raw, testCases[i].hmacSecret)
		if testCases[i].err != nil {
			assert.Equal(t, testCases[i].err.Error(), err.Error())
		} else {
			assert.Nil(t, err)
		}
		if err == nil {
			assert.Equal(t, testCases[i].subject, subject)
		}
	}
}
