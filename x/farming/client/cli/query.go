package cli

// DONTCOVER
// client is excluded from test coverage in MVP version

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/tendermint/farming/x/farming/types"
)

// GetQueryCmd returns a root CLI command handler for all x/farming query commands.
func GetQueryCmd() *cobra.Command {
	farmingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the farming module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	farmingQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryPlans(),
		GetCmdQueryPlan(),
		GetCmdQueryStakings(),
		GetCmdQueryTotalStakings(),
		GetCmdQueryRewards(),
	)

	return farmingQueryCmd
}

func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current farming parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as farming parameters.

Example:
$ %s query %s params
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			resp, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&resp.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryPlans() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plans [optional flags]",
		Args:  cobra.NoArgs,
		Short: "Query for all plans",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all farming plans on a network.

Example:
$ %s query %s plans
$ %s query %s plans --plan-type private
$ %s query %s plans --farming-pool-addr %s1zaavvzxez0elundtn32qnk9lkm8kmcszzsv80v
$ %s query %s plans --termination-addr %s1zaavvzxez0elundtn32qnk9lkm8kmcszzsv80v
$ %s query %s plans --staking-coin-denom poolD35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4
`,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName, sdk.Bech32MainPrefix,
				version.AppName, types.ModuleName, sdk.Bech32MainPrefix,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			planType, _ := cmd.Flags().GetString(FlagPlanType)
			farmingPoolAddr, _ := cmd.Flags().GetString(FlagFarmingPoolAddr)
			terminationAddr, _ := cmd.Flags().GetString(FlagTerminationAddr)
			stakingCoinDenom, _ := cmd.Flags().GetString(FlagStakingCoinDenom)

			var resp *types.QueryPlansResponse

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryPlansRequest{
				FarmingPoolAddress: farmingPoolAddr,
				TerminationAddress: terminationAddr,
				StakingCoinDenom:   stakingCoinDenom,
				Pagination:         pageReq,
			}
			if planType != "" {
				if planType == types.PlanTypePublic.String() || planType == types.PlanTypePrivate.String() {
					req.Type = planType
				} else {
					return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "plan type must be either public or private")
				}
			}

			resp, err = queryClient.Plans(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	cmd.Flags().AddFlagSet(flagSetPlans())
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "plans")

	return cmd
}

func GetCmdQueryPlan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plan [plan-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query a specific plan",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about a specific plan.

Example:
$ %s query %s plan
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			planId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "plan-id %s is not valid", args[0])
			}

			queryClient := types.NewQueryClient(clientCtx)

			resp, err := queryClient.Plan(cmd.Context(), &types.QueryPlanRequest{
				PlanId: planId,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryStakings() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "stakings [farmer]",
		Args:  cobra.ExactArgs(1),
		Short: "Query stakings by a farmer",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all stakings by a farmer.

Optionally restrict coins for a staking coin denom.

Example:
$ %s query %s stakings %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
$ %s query %s stakings %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj --staking-coin-denom poolD35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4
`,
				version.AppName, types.ModuleName, bech32PrefixAccAddr,
				version.AppName, types.ModuleName, bech32PrefixAccAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			farmerAcc, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			stakingCoinDenom, _ := cmd.Flags().GetString(FlagStakingCoinDenom)

			resp, err := queryClient.Stakings(cmd.Context(), &types.QueryStakingsRequest{
				Farmer:           farmerAcc.String(),
				StakingCoinDenom: stakingCoinDenom,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	cmd.Flags().AddFlagSet(flagSetStakings())
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryTotalStakings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-stakings [staking-coin-denom]",
		Args:  cobra.ExactArgs(1),
		Short: "Query total staking amount for a staking coin denom",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query total staking amount for a staking coin denom.

Example:
$ %s query %s total-stakings poolD35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			stakingCoinDenom := args[0]
			if err := sdk.ValidateDenom(stakingCoinDenom); err != nil {
				return err
			}

			resp, err := queryClient.TotalStakings(cmd.Context(), &types.QueryTotalStakingsRequest{
				StakingCoinDenom: stakingCoinDenom,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryRewards() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "rewards [farmer]",
		Args:  cobra.ExactArgs(1),
		Short: "Query rewards for a farmer",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all rewards for a farmer.

Optionally restrict rewards for a staking coin denom.

Example:
$ %s query %s rewards %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
$ %s query %s rewards %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj --staking-coin-denom poolD35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4
`,
				version.AppName, types.ModuleName, bech32PrefixAccAddr,
				version.AppName, types.ModuleName, bech32PrefixAccAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			farmerAcc, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			stakingCoinDenom, _ := cmd.Flags().GetString(FlagStakingCoinDenom)

			resp, err := queryClient.Rewards(cmd.Context(), &types.QueryRewardsRequest{
				Farmer:           farmerAcc.String(),
				StakingCoinDenom: stakingCoinDenom,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	cmd.Flags().AddFlagSet(flagSetRewards())
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
