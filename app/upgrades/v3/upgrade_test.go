package v3_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	app "github.com/dymensionxyz/dymension/v3/app"
	"github.com/dymensionxyz/dymension/v3/app/apptesting"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

// UpgradeTestSuite defines the structure for the upgrade test suite
type UpgradeTestSuite struct {
	suite.Suite
	Ctx types.Context
	App *app.App
}

// SetupTest initializes the necessary items for each test
func (s *UpgradeTestSuite) SetupTest(t *testing.T) {
	s.App = apptesting.Setup(t, false)
	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "dymension_110-1", Time: time.Now().UTC()})
}

// TestUpgradeTestSuite runs the suite of tests for the upgrade handler
func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

const (
	dummyUpgradeHeight          = 5
	dummyEndBlockHeight         = 10
	expectCreateGaugeFee        = 10_000_000_000_000_000_000
	expectAddToGaugeFee         = 0
	expectRollappsEnabled       = false
	expectDisputePeriodInBlocks = 120960
	expectMinBond               = 1_000_000_000
)

// TestUpgrade is a method of UpgradeTestSuite to test the upgrade process.
func (s *UpgradeTestSuite) TestUpgrade() {
	testCases := []struct {
		msg        string
		update     func()
		postUpdate func() error
		expPass    bool
	}{
		{
			"Test that upgrade does not panic and sets correct parameters",

			func() {
				// Run upgrade
				s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight - 1)
				plan := upgradetypes.Plan{Name: "v3", Height: dummyUpgradeHeight}
				err := s.App.UpgradeKeeper.ScheduleUpgrade(s.Ctx, plan)
				s.Require().NoError(err)
				_, exists := s.App.UpgradeKeeper.GetUpgradePlan(s.Ctx)
				s.Require().True(exists)

				s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight)
				// simulate the upgrade process not panic.
				s.Require().NotPanics(func() {
					// simulate the upgrade process.
					s.App.BeginBlocker(s.Ctx, abci.RequestBeginBlock{})
				})
			},
			func() error {
				// Post-update validation to ensure parameters are correctly set

				// Check Rollapp parameters
				rollappParams := s.App.RollappKeeper.GetParams(s.Ctx)
				if !rollappParams.RollappsEnabled != expectRollappsEnabled || rollappParams.DisputePeriodInBlocks != expectDisputePeriodInBlocks {
					return fmt.Errorf("rollapp parameters not set correctly")
				}

				// Check Sequencer parameters
				seqParams := s.App.SequencerKeeper.GetParams(s.Ctx)
				if seqParams.MinBond.Amount.Int64() != expectMinBond {
					return fmt.Errorf("sequencer parameters not set correctly")
				}

				// If you need to check more parameters, add the checks here

				return nil
			},
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			s.SetupTest(s.T()) // Reset for each case

			tc.update()
			err := tc.postUpdate()
			if tc.expPass {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}
