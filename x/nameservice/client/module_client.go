package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	nameservicecmd "github.com/yujianFresh/nameservice/x/nameservice/client/cli"
)

type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	namesvcQueryCmd := &cobra.Command{
		Use:   "nameservice",
		Short: "Querying commands for the nameservice module",
	}

	namesvcQueryCmd.AddCommand(
		client.GetCommands(
			nameservicecmd.GetCmdResolveName(mc.storeKey, mc.cdc),
			nameservicecmd.GetCmdWhois(mc.storeKey, mc.cdc),
		)...
	)

	return namesvcQueryCmd
}

// returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	namesvcTxCmd := &cobra.Command{
		Use:   "nameservice",
		Short: "Nameservice transactions subcommands",
	}

	namesvcTxCmd.AddCommand(
		client.PostCommands(
			nameservicecmd.GetCmdBuyName(mc.cdc),
			nameservicecmd.GetCmdSetName(mc.cdc),
		)...
	)

	return namesvcTxCmd
}
