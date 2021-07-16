package farming_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/farming/app"
	"github.com/tendermint/farming/x/farming"
	"github.com/tendermint/farming/x/farming/keeper"
	"github.com/tendermint/farming/x/farming/types"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

// createTestInput returns a simapp with custom FarmingKeeper
// to avoid messing with the hooks.
func createTestInput() (*app.FarmingApp, sdk.Context, []sdk.AccAddress) {
	app := app.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.FarmingKeeper = keeper.NewKeeper(
		app.AppCodec(),
		app.GetKey(types.StoreKey),
		app.GetSubspace(types.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.DistrKeeper,
		map[string]bool{},
	)

	// TODO: test_helpers.go
	addrs := []sdk.AccAddress{}
	// addrs := app.AddTestAddrs(app, ctx, 1, sdk.NewInt(100000))

	return app, ctx, addrs
}

func TestMsgCreateFixedAmountPlan(t *testing.T) {
	app, ctx, addrs := createTestInput()

	farmingPoolAddr := addrs[0]
	stakingCoinWeights := sdk.NewDecCoins(
		sdk.DecCoin{Denom: "testFarmStakingCoinDenom", Amount: sdk.MustNewDecFromStr("1.0")},
	)
	startTime := time.Now().UTC()
	endTime := startTime.AddDate(1, 0, 0)
	epochAmount := sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1)))

	msg := types.NewMsgCreateFixedAmountPlan(
		farmingPoolAddr,
		stakingCoinWeights,
		startTime,
		endTime,
		epochAmount,
	)

	handler := farming.NewHandler(app.FarmingKeeper)
	_, err := handler(ctx, msg)
	require.NoError(t, err)

	plans := app.FarmingKeeper.GetAllPlans(ctx)
	require.Equal(t, 1, len(plans))
	require.Equal(t, farmingPoolAddr.String(), plans[0].GetFarmingPoolAddress().String())
}

func TestMsgCreateRatioPlan(t *testing.T) {
	app, ctx, _ := createTestInput()

	farmingPoolAddr := sdk.AccAddress([]byte("farmingPoolAddr"))
	stakingCoinWeights := sdk.NewDecCoins(
		sdk.DecCoin{Denom: "testFarmStakingCoinDenom", Amount: sdk.MustNewDecFromStr("1.0")},
	)
	startTime := time.Now().UTC()
	endTime := startTime.AddDate(1, 0, 0)
	epochAmount := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1)))

	msg := types.NewMsgCreateFixedAmountPlan(
		farmingPoolAddr,
		stakingCoinWeights,
		startTime,
		endTime,
		epochAmount,
	)

	handler := farming.NewHandler(app.FarmingKeeper)
	_, err := handler(ctx, msg)
	require.NoError(t, err)

	plans := app.FarmingKeeper.GetAllPlans(ctx)
	require.Equal(t, 1, len(plans))
	require.Equal(t, farmingPoolAddr.String(), plans[0].GetFarmingPoolAddress().String())
}

func TestMsgStake(t *testing.T) {
	// TODO: not implemented yet
}

func TestMsgUnstake(t *testing.T) {
	// TODO: not implemented yet
}

func TestMsgHarvest(t *testing.T) {
	// TODO: not implemented yet
}
