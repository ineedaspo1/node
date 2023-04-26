package sso

import "testing"

func TestMystnodesSSOLink(t *testing.T) {
	// given
	sso := NewMystnodes()

	// expect
	sso.SSOLink()
}
