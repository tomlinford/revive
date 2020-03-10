package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUncheckedTypeAssertion(t *testing.T) {
	testRule(t, "unchecked-type-assertion", &rule.UncheckedTypeAssertionRule{})
}
