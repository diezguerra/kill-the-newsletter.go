package models

import "testing"

func TestConvertToRFC3339(t *testing.T) {
	timestamp := "2022-04-01 15:30:45.123456"
	expectedRFC3339 := "2022-04-01T15:30:45Z"

	actualRFC3339, err := ConvertToRFC3339(timestamp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if actualRFC3339 != expectedRFC3339 {
		t.Errorf("Expected RFC3339: %s, but got: %s", expectedRFC3339, actualRFC3339)
	}
}

func TestConvertToRFC3339_Failure(t *testing.T) {
	timestamp := "2022-04-01T15:30:45.123456"
	expectedRFC3339 := "2022-04-01T15:30:45.123456Z"

	actualRFC3339, err := ConvertToRFC3339(timestamp)
	if err == nil {
		t.Errorf("Expected error, but got no error")
	}

	if actualRFC3339 == expectedRFC3339 {
		t.Errorf("Expected failure, but conversion succeeded with RFC3339: %s", actualRFC3339)
	}
}

func TestConvertToRFC3339_NoTimeZone(t *testing.T) {
	timestamp := "2022-04-01 15:30:45.123456"
	expectedRFC3339 := "2022-04-01T15:30:45Z"

	actualRFC3339, err := ConvertToRFC3339(timestamp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if actualRFC3339 != expectedRFC3339 {
		t.Errorf("Expected RFC3339: %s, but got: %s", expectedRFC3339, actualRFC3339)
	}
}
