package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/RWAs-labs/muse/cmd/musetool/cctx"
	"github.com/RWAs-labs/muse/cmd/musetool/config"
	musecontext "github.com/RWAs-labs/muse/cmd/musetool/context"
)

func NewGetInboundBallotCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "get-ballot [inboundHash] [chainID]",
		Short: "fetch ballot identifier from the inbound hash",
		RunE:  GetInboundBallot,
		Args:  cobra.ExactArgs(2),
	}
}

func GetInboundBallot(cmd *cobra.Command, args []string) error {
	inboundHash := args[0]
	inboundChainID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse chain id")
	}
	configFile, err := cmd.Flags().GetString(config.FlagConfig)
	if err != nil {
		return fmt.Errorf("failed to read value for flag %s , err %w", config.FlagConfig, err)
	}

	ctx, err := musecontext.NewContext(context.Background(), inboundChainID, inboundHash, configFile)
	if err != nil {
		return fmt.Errorf("failed to create context: %w", err)
	}

	cctxTrackingDetails := cctx.NewTrackingDetails()

	err = cctxTrackingDetails.CheckInbound(ctx)
	if err != nil {
		return fmt.Errorf("failed to get ballot identifier: %w", err)
	}
	if cctxTrackingDetails.Status == cctx.PendingInboundConfirmation {
		log.Printf(
			"Ballot Identifier: %s, warning the inbound hash might not be confirmed yet",
			cctxTrackingDetails.CCTXIdentifier,
		)
		return nil
	}
	log.Print("Ballot Identifier: ", cctxTrackingDetails.CCTXIdentifier)
	return nil
}
