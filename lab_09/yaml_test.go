package main

import (
	"encoding/json"
	"testing"
)

func TestToYAML(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
		idErr    bool
	}{
		{
			name: "Valid Server Struct",
			input: Server{
				Host:       "localhost",
				Port:       8080,
				Debug:      true,
				AllowedIPs: []string{"192.168.1.1", "10.0.0.1"},
			},
			expected: "host: \"localhost\"\nport: 8080\ndebug: true\nallowed_ips: \n  - \"192.168.1.1\"\n  - \"10.0.0.1\"\n",
			idErr:    false,
		},
		{
			name: "Simple Struct",
			input: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{Name: "Tony", Age: 45},
			expected: "name: \"Tony\"\nage: 45\n",
			idErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToYAML(tt.input)
			if (err != nil) != tt.idErr {
				t.Errorf("ToYAML() error = %v, idErr %v", err, tt.idErr)
				return
			}
			if result != tt.expected {
				t.Errorf("ToYAML() got = %q, expected %q", result, tt.expected)
			}
		})
	}
}

var benchServer = Server{
	Host:       "localhost",
	Port:       8080,
	Debug:      true,
	AllowedIPs: []string{"192.168.1.1", "10.0.0.1", "172.16.0.1", "8.8.8.8"},
}

func BenchmarkToYAML(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ToYAML(benchServer)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkToJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(benchServer)
		if err != nil {
			b.Fatal(err)
		}
	}
}
