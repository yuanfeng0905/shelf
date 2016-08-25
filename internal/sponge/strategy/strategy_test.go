package strategy_test

import (
	"os"
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/sponge/strategy"
)

func init() {
	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	tests.Init("TEST")
}

// TestUpsertDelete tests if we can create a new strategy.
func TestNew(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	// set STRATEGY_CONF with test string
	oStrategy := os.Getenv("STRATEGY_CONF")

	strategyTestFile := os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/internal/sponge/strategy/sfix/strategy_coral_test.json"
	if err := os.Setenv("STRATEGY_CONF", strategyTestFile); err != nil {
		t.Fatalf("\t%s\tShould be able to set test strategy configuration string. : %v", tests.Failed, err)
	}

	defer func() {
		if err := os.Setenv("STRATEGY_CONF", oStrategy); err != nil {
			t.Fatalf("\t%s\tShould be able to set back strategy configuration string. : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to set back strategy configuration string.", tests.Success)
	}()

	t.Log("Given the need to get a new strategy instance.")
	{
		t.Log("\tWhen starting from an existing strategy file")
		{
			//----------------------------------------------------------------------
			// Get the strategy.

			s, err := strategy.New()
			if err != nil {
				t.Fatalf("\t%s\tShould be able to get the strategy : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get the strategy.", tests.Success)

			//----------------------------------------------------------------------
			// Check that we got the strategy we expected.

			if s.Name != "New York Times" || len(s.GetEntities()) != 3 {
				t.Logf("\t%+v", s.Name)
				t.Logf("\t%+v", len(s.GetEntities()))
				t.Fatalf("\t%s\tShould be able to get back the right strategy.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the right strategy.", tests.Success)

		}
	}
}
