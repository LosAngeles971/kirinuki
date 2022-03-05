package business

import (
	"testing"
)

const (
	test_email    = "losangeles971@gmail.com"
	test_password = "losangeles971@gmail.com"
)

func TestBasic(t *testing.T) {
	s, err := NewSession(test_email, test_password, true)
	if err != nil {
		t.Fatalf("failed to create session from scratch, %v", err)
	}
	if s.password != test_password {
		t.Fatalf("wrong password %s expected %s", s.password, test_password)
	}
	err = s.login()
	if err != nil {
		t.Fatalf("failed to login to an already open session, %v", err)
	}
	if !s.isOpen() {
		t.Fatal("session from scratch must be already open")
	}
}
