package parser

import (
	"testing"
)

func TestParseField_Valid(t *testing.T) {
	cases := []struct {
		value string
		ft    FieldType
	}{
		{"*", FieldMinute},
		{"0", FieldMinute},
		{"59", FieldMinute},
		{"*/5", FieldMinute},
		{"0-30", FieldMinute},
		{"0-30/5", FieldMinute},
		{"1,15,30", FieldMinute},
		{"0", FieldHour},
		{"23", FieldHour},
		{"*/2", FieldHour},
		{"1", FieldDayOfMonth},
		{"31", FieldDayOfMonth},
		{"1", FieldMonth},
		{"12", FieldMonth},
		{"1-6", FieldMonth},
		{"0", FieldDayOfWeek},
		{"7", FieldDayOfWeek},
		{"1-5", FieldDayOfWeek},
		{"1,3,5", FieldDayOfWeek},
	}

	for _, tc := range cases {
		t.Run(tc.value, func(t *testing.T) {
			if err := ParseField(tc.value, tc.ft); err != nil {
				t.Errorf("expected valid, got error: %v", err)
			}
		})
	}
}

func TestParseField_Invalid(t *testing.T) {
	cases := []struct {
		value string
		ft    FieldType
	}{
		{"60", FieldMinute},
		{"-1", FieldMinute},
		{"abc", FieldMinute},
		{"24", FieldHour},
		{"0", FieldDayOfMonth},
		{"32", FieldDayOfMonth},
		{"0", FieldMonth},
		{"13", FieldMonth},
		{"8", FieldDayOfWeek},
		{"5-3", FieldDayOfWeek},
		{"*/0", FieldMinute},
		{"1-100", FieldMinute},
		{"foo/2", FieldHour},
	}

	for _, tc := range cases {
		t.Run(tc.value, func(t *testing.T) {
			if err := ParseField(tc.value, tc.ft); err == nil {
				t.Errorf("expected error for %q on field %d, got nil", tc.value, tc.ft)
			}
		})
	}
}
