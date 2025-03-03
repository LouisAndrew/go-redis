package main

import (
	"strings"
	"testing"
)

type testCase[T any] struct {
	input     string
	want      T
	desc      string
	shouldErr bool
}

func TestRespReadLine(t *testing.T) {
	testCases := []testCase[string]{
		{
			input: "Louis\r\n",
			want:  "Louis",
			desc:  "Returns valid string correctly",
		},
		{
			input: "",
			want:  "",
			desc:  "Returns empty string without delimiter correctly",
		},
		{
			input: "\r\n",
			want:  "",
			desc:  "Returns empty string if input consists only of delimiters",
		},
		{
			input: "Louis\r\nRest",
			want:  "Louis",
			desc:  "Returns only the part of string until the delimiter",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, _, _ := NewResp(strings.NewReader(tc.input)).readLine()
			if string(got) != tc.want {
				t.Errorf("got %v want %v", got, tc.want)
			}
		})
	}
}

func TestRespReadInt(t *testing.T) {
	testCases := []testCase[int]{
		{
			input: "12\r\n",
			want:  12,
			desc:  "should return a valid integer line",
		},
		{
			input:     "Louis\r\n",
			want:      0,
			desc:      "should throw an error if the string doesn't contain a valid integer",
			shouldErr: true,
		},
		{
			input: "123456\r\n",
			want:  123456,
			desc:  "should return a valid big integer",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, _, err := NewResp(strings.NewReader(tc.input)).readInt()
			if err != nil && !tc.shouldErr {
				t.Errorf("Uncaught error although it shouldn't -> %v", err)
			}

			if got != tc.want {
				t.Errorf("got %v want %v", got, tc.want)
			}
		})
	}
}
