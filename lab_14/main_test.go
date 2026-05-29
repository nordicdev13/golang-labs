package main

import "testing"

func TestIsEven(t *testing.T) {
	if !IsEven(2) {
		t.Errorf("IsEven(2) повернуло false, очікувалось true")
	}

	if IsEven(3) {
		t.Errorf("IsEven(3) повернуло true, очікувалось false")
	}
}
