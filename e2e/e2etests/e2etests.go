package e2etests

import (
	"github.com/RWAs-labs/muse/e2e/e2etests/legacy"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/testutil/sample"
)

// List of all e2e test names to be used in musee2e
const (
	/*
	  EVM chain tests
	*/
	TestETHDepositName                      = "eth_deposit"
	TestETHDepositAndCallName               = "eth_deposit_and_call"
	TestETHDepositFastConfirmationName      = "eth_deposit_fast_confirmation"
	TestETHDepositAndCallNoMessageName      = "eth_deposit_and_call_no_message"
	TestETHDepositAndCallRevertName         = "eth_deposit_and_call_revert"
	TestETHDepositAndCallRevertWithCallName = "eth_deposit_and_call_revert_with_call"
	TestETHDepositRevertAndAbortName        = "eth_deposit_revert_and_abort"

	TestETHWithdrawName                          = "eth_withdraw"
	TestETHWithdrawAndArbitraryCallName          = "eth_withdraw_and_arbitrary_call"
	TestETHWithdrawAndCallName                   = "eth_withdraw_and_call"
	TestETHWithdrawAndCallNoMessageName          = "eth_withdraw_and_call_no_message"
	TestETHWithdrawAndCallThroughContractName    = "eth_withdraw_and_call_through_contract"
	TestETHWithdrawAndCallRevertName             = "eth_withdraw_and_call_revert"
	TestETHWithdrawAndCallRevertWithCallName     = "eth_withdraw_and_call_revert_with_call"
	TestETHWithdrawRevertAndAbortName            = "eth_withdraw_revert_and_abort"
	TestETHWithdrawAndCallRevertWithWithdrawName = "eth_withdraw_and_call_revert_with_withdraw"
	TestDepositAndCallOutOfGasName               = "deposit_and_call_out_of_gas"

	TestERC20DepositName                      = "erc20_deposit"
	TestERC20DepositAndCallName               = "erc20_deposit_and_call"
	TestERC20DepositAndCallNoMessageName      = "erc20_deposit_and_call_no_message"
	TestERC20DepositAndCallRevertName         = "erc20_deposit_and_call_revert"
	TestERC20DepositAndCallRevertWithCallName = "erc20_deposit_and_call_revert_with_call"
	TestERC20DepositRevertAndAbortName        = "erc20_deposit_revert_and_abort"

	TestERC20WithdrawName                      = "erc20_withdraw"
	TestERC20WithdrawAndArbitraryCallName      = "erc20_withdraw_and_arbitrary_call"
	TestERC20WithdrawAndCallName               = "erc20_withdraw_and_call"
	TestERC20WithdrawAndCallNoMessageName      = "erc20_withdraw_and_call_no_message"
	TestERC20WithdrawAndCallRevertName         = "erc20_withdraw_and_call_revert"
	TestERC20WithdrawAndCallRevertWithCallName = "erc20_withdraw_and_call_revert_with_call"
	TestERC20WithdrawRevertAndAbortName        = "erc20_withdraw_revert_and_abort"

	TestMEVMToEVMArbitraryCallName       = "mevm_to_evm_arbitrary_call"
	TestMEVMToEVMCallName                = "mevm_to_evm_call"
	TestMEVMToEVMCallRevertName          = "mevm_to_evm_call_revert"
	TestMEVMToEVMCallRevertAndAbortName  = "mevm_to_evm_call_revert_and_abort"
	TestMEVMToEVMCallThroughContractName = "mevm_to_evm_call_through_contract"
	TestEVMToMEVMCallName                = "evm_to_mevm_call"
	TestEVMToMEVMCallAbortName           = "evm_to_mevm_abort_call"

	TestDepositAndCallSwapName      = "deposit_and_call_swap"
	TestEtherWithdrawRestrictedName = "eth_withdraw_restricted"
	TestERC20DepositRestrictedName  = "erc20_deposit_restricted" // #nosec G101: Potential hardcoded credentials (gosec), not a credential

	/*
	 * Solana tests
	 */
	TestSolanaDepositName                                 = "solana_deposit"
	TestSolanaDepositThroughProgramName                   = "solana_deposit_through_program"
	TestSolanaWithdrawName                                = "solana_withdraw"
	TestSolanaWithdrawRevertExecutableReceiverName        = "solana_withdraw_revert_executable_receiver"
	TestSolanaWithdrawAndCallName                         = "solana_withdraw_and_call"
	TestSolanaWithdrawAndCallInvalidMsgEncodingName       = "solana_withdraw_and_call_invalid_msg_encoding"
	TestMEVMToSolanaCallName                              = "mevm_to_solana_call"
	TestSolanaWithdrawAndCallRevertWithCallName           = "solana_withdraw_and_call_revert_with_call"
	TestSolanaDepositAndCallName                          = "solana_deposit_and_call"
	TestSolanaDepositAndCallRevertName                    = "solana_deposit_and_call_revert"
	TestSolanaDepositAndCallRevertWithCallName            = "solana_deposit_and_call_revert_with_call"
	TestSolanaDepositAndCallRevertWithCallThatRevertsName = "solana_deposit_and_call_revert_with_call_that_reverts"
	TestSolanaDepositAndCallRevertWithDustName            = "solana_deposit_and_call_revert_with_dust"
	TestSolanaToMEVMCallName                              = "solana_to_mevm_call"
	TestSolanaToMEVMCallAbortName                         = "solana_to_mevm_call_abort"
	TestSolanaDepositRestrictedName                       = "solana_deposit_restricted"
	TestSolanaWithdrawRestrictedName                      = "solana_withdraw_restricted"
	TestSPLDepositName                                    = "spl_deposit"
	TestSPLDepositAndCallName                             = "spl_deposit_and_call"
	TestSPLDepositAndCallRevertName                       = "spl_deposit_and_call_revert"
	TestSPLDepositAndCallRevertWithCallName               = "spl_deposit_and_call_revert_with_call"
	TestSPLDepositAndCallRevertWithCallThatRevertsName    = "spl_deposit_and_call_revert_with_call_that_reverts"
	TestSPLWithdrawName                                   = "spl_withdraw"
	TestSPLWithdrawAndCallName                            = "spl_withdraw_and_call"
	TestSPLWithdrawAndCallRevertName                      = "spl_withdraw_and_call_revert"
	TestSPLWithdrawAndCreateReceiverAtaName               = "spl_withdraw_and_create_receiver_ata"

	/**
	 * TON tests
	 */
	TestTONDepositName              = "ton_deposit"
	TestTONDepositAndCallName       = "ton_deposit_and_call"
	TestTONDepositAndCallRefundName = "ton_deposit_refund"
	TestTONDepositRestrictedName    = "ton_deposit_restricted"
	TestTONWithdrawName             = "ton_withdraw"
	TestTONWithdrawConcurrentName   = "ton_withdraw_concurrent"

	/*
	 Sui tests
	*/
	TestSuiDepositName                            = "sui_deposit"
	TestSuiDepositAndCallName                     = "sui_deposit_and_call"
	TestSuiDepositAndCallRevertName               = "sui_deposit_and_call_revert"
	TestSuiTokenDepositName                       = "sui_token_deposit"                 // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestSuiTokenDepositAndCallName                = "sui_token_deposit_and_call"        // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestSuiTokenDepositAndCallRevertName          = "sui_token_deposit_and_call_revert" // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestSuiWithdrawName                           = "sui_withdraw"
	TestSuiTokenWithdrawName                      = "sui_token_withdraw"                           // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestSuiTokenWithdrawAndCallName               = "sui_token_withdraw_and_call"                  // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestSuiTokenWithdrawAndCallRevertWithCallName = "sui_token_withdraw_and_call_revert_with_call" // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestSuiWithdrawAndCallName                    = "sui_withdraw_and_call"
	TestSuiWithdrawRevertWithCallName             = "sui_withdraw_revert_with_call"          // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestSuiWithdrawAndCallRevertWithCallName      = "sui_withdraw_and_call_revert_with_call" // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestSuiDepositRestrictedName                  = "sui_deposit_restricted"
	TestSuiWithdrawRestrictedName                 = "sui_withdraw_restricted"
	TestSuiWithdrawInvalidReceiverName            = "sui_withdraw_invalid_receiver"

	/*
	 Bitcoin tests
	 Test transfer of Bitcoin asset across chains
	*/
	TestBitcoinDepositName                                 = "bitcoin_deposit"
	TestBitcoinDepositAndCallName                          = "bitcoin_deposit_and_call"
	TestBitcoinDepositFastConfirmationName                 = "bitcoin_deposit_fast_confirmation"
	TestBitcoinDepositAndCallRevertName                    = "bitcoin_deposit_and_call_revert"
	TestBitcoinDepositAndCallRevertWithDustName            = "bitcoin_deposit_and_call_revert_with_dust"
	TestBitcoinDepositAndWithdrawWithDustName              = "bitcoin_deposit_and_withdraw_with_dust"
	TestBitcoinDonationName                                = "bitcoin_donation"
	TestBitcoinStdMemoDepositName                          = "bitcoin_std_memo_deposit"
	TestBitcoinStdMemoDepositAndCallName                   = "bitcoin_std_memo_deposit_and_call"
	TestBitcoinStdMemoDepositAndCallRevertName             = "bitcoin_std_memo_deposit_and_call_revert"
	TestBitcoinStdMemoDepositAndCallRevertOtherAddressName = "bitcoin_std_memo_deposit_and_call_revert_other_address"
	TestBitcoinStdMemoDepositAndCallRevertAndAbortName     = "bitcoin_std_memo_deposit_and_call_revert_and_abort"
	TestBitcoinStdMemoInscribedDepositAndCallName          = "bitcoin_std_memo_inscribed_deposit_and_call"
	TestBitcoinDepositAndAbortWithLowDepositFeeName        = "bitcoin_deposit_and_abort_with_low_deposit_fee"
	TestBitcoinWithdrawSegWitName                          = "bitcoin_withdraw_segwit"
	TestBitcoinWithdrawTaprootName                         = "bitcoin_withdraw_taproot"
	TestBitcoinWithdrawMultipleName                        = "bitcoin_withdraw_multiple"
	TestBitcoinWithdrawLegacyName                          = "bitcoin_withdraw_legacy"
	TestBitcoinWithdrawP2WSHName                           = "bitcoin_withdraw_p2wsh"
	TestBitcoinWithdrawP2SHName                            = "bitcoin_withdraw_p2sh"
	TestBitcoinWithdrawInvalidAddressName                  = "bitcoin_withdraw_invalid"
	TestBitcoinWithdrawRestrictedName                      = "bitcoin_withdraw_restricted"
	TestBitcoinDepositInvalidMemoRevertName                = "bitcoin_deposit_invalid_memo_revert"
	TestBitcoinWithdrawRBFName                             = "bitcoin_withdraw_rbf"

	/*
	 Application tests
	 Test various smart contract applications across chains
	*/
	TestCrosschainSwapName = "crosschain_swap"

	/*
	 Miscellaneous tests
	 Test various functionalities not related to assets
	*/
	TestDonationEtherName   = "donation_ether"
	TestInboundTrackersName = "inbound_trackers"

	/*
	 Stress tests
	 Test stressing networks with many cross-chain transactions
	*/
	TestStressEtherWithdrawName  = "stress_eth_withdraw"
	TestStressBTCWithdrawName    = "stress_btc_withdraw"
	TestStressEtherDepositName   = "stress_eth_deposit"
	TestStressBTCDepositName     = "stress_btc_deposit"
	TestStressSolanaDepositName  = "stress_solana_deposit"
	TestStressSPLDepositName     = "stress_spl_deposit"
	TestStressSolanaWithdrawName = "stress_solana_withdraw"
	TestStressSPLWithdrawName    = "stress_spl_withdraw"
	TestStressSuiDepositName     = "stress_sui_deposit"
	TestStressSuiWithdrawName    = "stress_sui_withdraw"

	/*
		Staking tests
	*/

	TestUndelegateToBelowMinimumObserverDelegation = "undelegate_to_below_minimum_observer_delegation"

	/*
	 Admin tests
	 Test admin functionalities
	*/
	TestWhitelistERC20Name               = "whitelist_erc20"
	TestDepositEtherLiquidityCapName     = "deposit_eth_liquidity_cap"
	TestMigrateChainSupportName          = "migrate_chain_support"
	TestPauseMRC20Name                   = "pause_mrc20"
	TestUpdateBytecodeMRC20Name          = "update_bytecode_mrc20"
	TestUpdateBytecodeConnectorName      = "update_bytecode_connector"
	TestRateLimiterName                  = "rate_limiter"
	TestCriticalAdminTransactionsName    = "critical_admin_transactions"
	TestPauseERC20CustodyName            = "pause_erc20_custody"
	TestMigrateERC20CustodyFundsName     = "migrate_erc20_custody_funds"
	TestMigrateTSSName                   = "migrate_tss"
	TestSolanaWhitelistSPLName           = "solana_whitelist_spl"
	TestUpdateMRC20NameName              = "update_mrc20_name"
	TestMuseclientRestartHeightName      = "museclient_restart_height"
	TestMuseclientSignerOffsetName       = "museclient_signer_offset"
	TestUpdateOperationalChainParamsName = "update_operational_chain_params"

	/*
	 Operational tests
	 Not used to test functionalities but do various interactions with the netwoks
	*/
	TestDeploy                                    = "deploy"
	TestOperationAddLiquidityETHName              = "add_liquidity_eth"
	TestOperationAddLiquidityERC20Name            = "add_liquidity_erc20"
	TestOperationAddLiquidityBTCName              = "add_liquidity_btc"
	TestOperationAddLiquiditySOLName              = "add_liquidity_sol"
	TestOperationAddLiquiditySPLName              = "add_liquidity_spl"
	TestOperationAddLiquiditySUIName              = "add_liquidity_sui"
	TestOperationAddLiquiditySuiFungibleTokenName = "add_liquidity_sui_fungible_token" // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestOperationAddLiquidityTONName              = "add_liquidity_ton"

	/*
	 Stateful precompiled contracts tests
	*/
	TestPrecompilesPrototypeName                 = "precompile_contracts_prototype"
	TestPrecompilesPrototypeThroughContractName  = "precompile_contracts_prototype_through_contract"
	TestPrecompilesStakingName                   = "precompile_contracts_staking"
	TestPrecompilesStakingThroughContractName    = "precompile_contracts_staking_through_contract"
	TestPrecompilesBankName                      = "precompile_contracts_bank"
	TestPrecompilesBankFailName                  = "precompile_contracts_bank_fail"
	TestPrecompilesBankThroughContractName       = "precompile_contracts_bank_through_contract"
	TestPrecompilesDistributeName                = "precompile_contracts_distribute"
	TestPrecompilesDistributeNonMRC20Name        = "precompile_contracts_distribute_non_mrc20"
	TestPrecompilesDistributeThroughContractName = "precompile_contracts_distribute_through_contract"

	/*
	 Legacy tests (using v1 protocol contracts)
	*/
	TestLegacyMessagePassingExternalChainsName              = "legacy_message_passing_external_chains"
	TestLegacyMessagePassingRevertFailExternalChainsName    = "legacy_message_passing_revert_fail"
	TestLegacyMessagePassingRevertSuccessExternalChainsName = "legacy_message_passing_revert_success"
	TestLegacyMessagePassingEVMtoMEVMName                   = "legacy_message_passing_evm_to_mevm"
	TestLegacyMessagePassingMEVMToEVMName                   = "legacy_message_passing_mevm_to_evm"
	TestLegacyMessagePassingMEVMtoEVMRevertName             = "legacy_message_passing_mevm_to_evm_revert"
	TestLegacyMessagePassingEVMtoMEVMRevertName             = "legacy_message_passing_evm_to_mevm_revert"
	TestLegacyMessagePassingMEVMtoEVMRevertFailName         = "legacy_message_passing_mevm_to_evm_revert_fail"
	TestLegacyMessagePassingEVMtoMEVMRevertFailName         = "legacy_message_passing_evm_to_mevm_revert_fail"
	TestLegacyEtherDepositName                              = "legacy_eth_deposit"
	TestLegacyEtherWithdrawName                             = "legacy_eth_withdraw"
	TestLegacyEtherDepositAndCallRefundName                 = "legacy_eth_deposit_and_call_refund"
	TestLegacyEtherDepositAndCallName                       = "legacy_eth_deposit_and_call"
	TestLegacyERC20WithdrawName                             = "legacy_erc20_withdraw"
	TestLegacyERC20DepositName                              = "legacy_erc20_deposit"
	TestLegacyMultipleERC20DepositName                      = "legacy_erc20_multiple_deposit"
	TestLegacyMultipleERC20WithdrawsName                    = "legacy_erc20_multiple_withdraw"
	TestLegacyERC20DepositAndCallRefundName                 = "legacy_erc20_deposit_and_call_refund"

	/*
	 MUSE tests
	 Test transfer of MUSE asset across chains
	 Note: It is still the only way to transfer MUSE across chains. Work to integrate MUSE transfers as part of the gateway is in progress
	 These tests are marked as legacy because there is no longer active development on MUSE transfers, and we stopped integrating MUSE support on new mainnet chains
	*/
	TestLegacyMuseDepositName           = "legacy_muse_deposit"
	TestLegacyMuseDepositNewAddressName = "legacy_muse_deposit_new_address"
	TestLegacyMuseDepositRestrictedName = "legacy_muse_deposit_restricted"
	TestLegacyMuseWithdrawName          = "legacy_muse_withdraw"
	TestLegacyMuseWithdrawBTCRevertName = "legacy_muse_withdraw_btc_revert" // #nosec G101 - not a hardcoded password

)

const (
	CountArgDescription = "count"
)

// Here are all the dependencies for the e2e tests, add more dependencies here if needed
var (
	// DepdencyAllBitcoinDeposits is a dependency to wait for all bitcoin deposit tests to complete
	DepdencyAllBitcoinDeposits = runner.NewE2EDependency("all_bitcoin_deposits")
)

// AllE2ETests is an ordered list of all e2e tests
var AllE2ETests = []runner.E2ETest{
	/*
	 EVM chain tests
	*/
	runner.NewE2ETest(
		TestETHDepositName,
		"deposit Ether into MEVM",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000000000000000000"},
		},
		TestETHDeposit,
	),
	runner.NewE2ETest(
		TestETHDepositAndCallName,
		"deposit Ether into MEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestETHDepositAndCall,
	),
	runner.NewE2ETest(
		TestETHDepositFastConfirmationName,
		"deposit Ether into MEVM using fast confirmation",
		[]runner.ArgDefinition{},
		TestETHDepositFastConfirmation,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestETHDepositAndCallNoMessageName,
		"deposit Ether into MEVM and call a contract using no message content",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestETHDepositAndCallNoMessage,
	),
	runner.NewE2ETest(
		TestETHDepositAndCallRevertName,
		"deposit Ether into MEVM and call a contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestETHDepositAndCallRevert,
	),
	runner.NewE2ETest(
		TestETHDepositAndCallRevertWithCallName,
		"deposit Ether into MEVM and call a contract that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestETHDepositAndCallRevertWithCall,
	),
	runner.NewE2ETest(
		TestETHDepositRevertAndAbortName,
		"deposit Ether into MEVM, revert, then abort with onAbort",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestETHDepositRevertAndAbort,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestETHWithdrawName,
		"withdraw Ether from MEVM",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestETHWithdraw,
	),
	runner.NewE2ETest(
		TestETHWithdrawAndArbitraryCallName,
		"withdraw Ether from MEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestETHWithdrawAndArbitraryCall,
	),
	runner.NewE2ETest(
		TestETHWithdrawAndCallName,
		"withdraw Ether from MEVM call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
			{Description: "gas limit for withdraw", DefaultValue: "250000"},
		},
		TestETHWithdrawAndCall,
	),
	runner.NewE2ETest(
		TestETHWithdrawAndCallNoMessageName,
		"withdraw Ether from MEVM call a contract with no message content",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
			{Description: "gas limit for withdraw", DefaultValue: "250000"},
		},
		TestETHWithdrawAndCallNoMessage,
	),
	runner.NewE2ETest(
		TestETHWithdrawAndCallThroughContractName,
		"withdraw Ether from MEVM call a contract through intermediary contract",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestETHWithdrawAndCallThroughContract,
	),
	runner.NewE2ETest(
		TestETHWithdrawAndCallRevertName,
		"withdraw Ether from MEVM and call a contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestETHWithdrawAndCallRevert,
	),
	runner.NewE2ETest(
		TestETHWithdrawAndCallRevertWithCallName,
		"withdraw Ether from MEVM and call a contract that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestETHWithdrawAndCallRevertWithCall,
	),
	runner.NewE2ETest(
		TestETHWithdrawRevertAndAbortName,
		"withdraw Ether from MEVM, revert, then abort with onAbort, check onAbort can created cctx",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "1000000000000000000"},
			{Description: "gas limit for withdraw", DefaultValue: "250000"},
		},
		TestETHWithdrawRevertAndAbort,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestETHWithdrawAndCallRevertWithWithdrawName,
		"withdraw Ether from MEVM and call a contract that reverts with a onRevert call that triggers a withdraw",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestETHWithdrawAndCallRevertWithWithdraw,
		runner.WithMinimumVersion("v26.0.0"),
	),
	runner.NewE2ETest(
		TestDepositAndCallOutOfGasName,
		"deposit Ether into MEVM and call a contract that runs out of gas",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestDepositAndCallOutOfGas,
	),
	runner.NewE2ETest(
		TestERC20DepositName,
		"deposit ERC20 into MEVM",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100000000000000000000"},
		},
		TestERC20Deposit,
	),
	runner.NewE2ETest(
		TestERC20DepositAndCallName,
		"deposit ERC20 into MEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100000"},
		},
		TestERC20DepositAndCall,
	),
	runner.NewE2ETest(
		TestERC20DepositAndCallNoMessageName,
		"deposit ERC20 into MEVM and call a contract with no message content",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100000"},
		},
		TestERC20DepositAndCallNoMessage,
	),
	runner.NewE2ETest(
		TestERC20DepositAndCallRevertName,
		"deposit ERC20 into MEVM and call a contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "10000000000000000000"},
		},
		TestERC20DepositAndCallRevert,
	),
	runner.NewE2ETest(
		TestERC20DepositAndCallRevertWithCallName,
		"deposit ERC20 into MEVM and call a contract that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "10000000000000000000"},
		},
		TestERC20DepositAndCallRevertWithCall,
	),
	runner.NewE2ETest(
		TestERC20DepositRevertAndAbortName,
		"deposit ERC20 into MEVM, revert, then abort with onAbort because revert fee cannot be paid",
		[]runner.ArgDefinition{},
		TestERC20DepositRevertAndAbort,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestERC20WithdrawName,
		"withdraw ERC20 from MEVM",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestERC20Withdraw,
	),
	runner.NewE2ETest(
		TestERC20WithdrawAndArbitraryCallName,
		"withdraw ERC20 from MEVM and arbitrary call a contract",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestERC20WithdrawAndArbitraryCall,
	),
	runner.NewE2ETest(
		TestERC20WithdrawAndCallName,
		"withdraw ERC20 from MEVM and authenticated call a contract",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestERC20WithdrawAndCall,
	),
	runner.NewE2ETest(
		TestERC20WithdrawAndCallNoMessageName,
		"withdraw ERC20 from MEVM and authenticated call a contract with no message",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestERC20WithdrawAndCallNoMessage,
	),
	runner.NewE2ETest(
		TestERC20WithdrawAndCallRevertName,
		"withdraw ERC20 from MEVM and call a contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestERC20WithdrawAndCallRevert,
	),
	runner.NewE2ETest(
		TestERC20WithdrawAndCallRevertWithCallName,
		"withdraw ERC20 from MEVM and call a contract that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestERC20WithdrawAndCallRevertWithCall,
	),
	runner.NewE2ETest(
		TestERC20WithdrawRevertAndAbortName,
		"withdraw ERC20 from MEVM, revert, then abort with onAbort",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestERC20WithdrawRevertAndAbort,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestMEVMToEVMArbitraryCallName,
		"mevm -> evm call",
		[]runner.ArgDefinition{},
		TestMEVMToEVMArbitraryCall,
	),
	runner.NewE2ETest(
		TestMEVMToEVMCallName,
		"mevm -> evm call",
		[]runner.ArgDefinition{},
		TestMEVMToEVMCall,
	),
	runner.NewE2ETest(
		TestMEVMToEVMCallRevertName,
		"mevm -> evm call that reverts and call onRevert",
		[]runner.ArgDefinition{},
		TestMEVMToEVMCallRevert,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestMEVMToEVMCallRevertAndAbortName,
		"mevm -> evm call that reverts and abort with onAbort",
		[]runner.ArgDefinition{},
		TestMEVMToEVMCallRevertAndAbort,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestMEVMToEVMCallThroughContractName,
		"mevm -> evm call through intermediary contract",
		[]runner.ArgDefinition{},
		TestMEVMToEVMCallThroughContract,
	),
	runner.NewE2ETest(
		TestEVMToMEVMCallName,
		"evm -> mevm call",
		[]runner.ArgDefinition{},
		TestEVMToMEVMCall,
	),
	runner.NewE2ETest(
		TestEVMToMEVMCallAbortName,
		"evm -> mevm call fails and abort with onAbort",
		[]runner.ArgDefinition{},
		TestEVMToMEVMCallAbort,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestDepositAndCallSwapName,
		"evm -> mevm deposit and call with swap and withdraw back to evm",
		[]runner.ArgDefinition{},
		TestDepositAndCallSwap,
	),

	/*
	 Solana tests
	*/
	runner.NewE2ETest(
		TestSolanaDepositName,
		"deposit SOL into MEVM",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "24000000"},
		},
		TestSolanaDeposit,
	),
	runner.NewE2ETest(
		TestSolanaDepositThroughProgramName,
		"deposit SOL into MEVM through example connected program",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "24000000"},
		},
		TestSolanaDepositThroughProgram,
	),
	runner.NewE2ETest(
		TestSolanaWithdrawName,
		"withdraw SOL from MEVM",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1000000"},
		},
		TestSolanaWithdraw,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaWithdrawAndCallName,
		"withdraw SOL from MEVM and call solana program",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1000000"},
		},
		TestSolanaWithdrawAndCall,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaWithdrawRevertExecutableReceiverName,
		"withdraw SOL from MEVM reverts if executable receiver",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1000000"},
		},
		TestSolanaWithdrawRevertExecutableReceiver,
	),
	runner.NewE2ETest(
		TestMEVMToSolanaCallName,
		"call solana program from MEVM",
		[]runner.ArgDefinition{},
		TestMEVMToSolanaCall,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaWithdrawAndCallInvalidMsgEncodingName,
		"withdraw SOL from MEVM and call solana program with invalid msg encoding",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1000000"},
		},
		TestSolanaWithdrawAndCallInvalidMsgEncoding,
	),
	runner.NewE2ETest(
		TestSolanaWithdrawAndCallRevertWithCallName,
		"withdraw SOL from MEVM and call solana program that reverts",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1000000"},
		},
		TestSolanaWithdrawAndCallRevertWithCall,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSPLWithdrawAndCallName,
		"withdraw SPL from MEVM and call solana program",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1000000"},
		},
		TestSPLWithdrawAndCall,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSPLWithdrawAndCallRevertName,
		"withdraw SPL from MEVM and call solana program that reverts",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1000000"},
		},
		TestSPLWithdrawAndCallRevert,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaDepositAndCallName,
		"deposit SOL into MEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1200000"},
		},
		TestSolanaDepositAndCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaToMEVMCallName,
		"call a mevm contract",
		[]runner.ArgDefinition{},
		TestSolanaToMEVMCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaToMEVMCallAbortName,
		"call a mevm contract and abort",
		[]runner.ArgDefinition{},
		TestSolanaToMEVMCallAbort,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSPLWithdrawName,
		"withdraw SPL from MEVM",
		[]runner.ArgDefinition{
			{Description: "amount in spl tokens", DefaultValue: "100000"},
		},
		TestSPLWithdraw,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSPLWithdrawAndCreateReceiverAtaName,
		"withdraw SPL from MEVM and create receiver ata",
		[]runner.ArgDefinition{
			{Description: "amount in spl tokens", DefaultValue: "1000000"},
		},
		TestSPLWithdrawAndCreateReceiverAta,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaDepositAndCallRevertName,
		"deposit SOL into MEVM and call a contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1200000"},
		},
		TestSolanaDepositAndCallRevert,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaDepositAndCallRevertWithCallName,
		"deposit SOL into MEVM and call a contract that reverts with call on revert",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1200000"},
		},
		TestSolanaDepositAndCallRevertWithCall,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaDepositAndCallRevertWithCallThatRevertsName,
		"deposit SOL into MEVM and call a contract that reverts with call on revert and connected program reverts",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1200000"},
		},
		TestSolanaDepositAndCallRevertWithCallThatReverts,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaDepositAndCallRevertWithDustName,
		"deposit SOL into MEVM; revert with dust amount that aborts the CCTX",
		[]runner.ArgDefinition{},
		TestSolanaDepositAndCallRevertWithDust,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaDepositRestrictedName,
		"deposit SOL into MEVM restricted address",
		[]runner.ArgDefinition{
			{Description: "receiver", DefaultValue: sample.RestrictedEVMAddressTest},
			{Description: "amount in lamport", DefaultValue: "1200000"},
		},
		TestSolanaDepositRestricted,
	),
	runner.NewE2ETest(
		TestSolanaWithdrawRestrictedName,
		"withdraw SOL from MEVM to restricted address",
		[]runner.ArgDefinition{
			{Description: "receiver", DefaultValue: sample.RestrictedSolAddressTest},
			{Description: "amount in lamport", DefaultValue: "1000000"},
		},
		TestSolanaWithdrawRestricted,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaWhitelistSPLName,
		"whitelist SPL",
		[]runner.ArgDefinition{},
		TestSolanaWhitelistSPL,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSPLDepositName,
		"deposit SPL into MEVM",
		[]runner.ArgDefinition{
			{Description: "amount of spl tokens", DefaultValue: "24000000"},
		},
		TestSPLDeposit,
	),
	runner.NewE2ETest(
		TestSPLDepositAndCallName,
		"deposit SPL into MEVM and call",
		[]runner.ArgDefinition{
			{Description: "amount of spl tokens", DefaultValue: "12000000"},
		},
		TestSPLDepositAndCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSPLDepositAndCallRevertName,
		"deposit SPL into MEVM and call which reverts",
		[]runner.ArgDefinition{
			{Description: "amount of spl tokens", DefaultValue: "12000000"},
		},
		TestSPLDepositAndCallRevert,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSPLDepositAndCallRevertWithCallName,
		"deposit SPL into MEVM and call which reverts with call on revert",
		[]runner.ArgDefinition{
			{Description: "amount of spl tokens", DefaultValue: "12000000"},
		},
		TestSPLDepositAndCallRevertWithCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSPLDepositAndCallRevertWithCallThatRevertsName,
		"deposit SPL into MEVM and call which reverts with call on revert that reverts",
		[]runner.ArgDefinition{
			{Description: "amount of spl tokens", DefaultValue: "12000000"},
		},
		TestSPLDepositAndCallRevertWithCallThatReverts,
		runner.WithMinimumVersion("v30.0.0"),
	),
	/*
	 TON tests
	*/
	runner.NewE2ETest(
		TestTONDepositName,
		"deposit TON into MEVM",
		[]runner.ArgDefinition{
			{Description: "amount in nano tons", DefaultValue: "1000000000"}, // 1.0 TON
		},
		TestTONDeposit,
	),
	runner.NewE2ETest(
		TestTONDepositAndCallName,
		"deposit TON into MEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in nano tons", DefaultValue: "1000000000"}, // 1.0 TON
		},
		TestTONDepositAndCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestTONDepositAndCallRefundName,
		"deposit TON into MEVM and call a smart contract that reverts; expect refund",
		[]runner.ArgDefinition{
			{Description: "amount in nano tons", DefaultValue: "1000000000"}, // 1.0 TON
		},
		TestTONDepositAndCallRefund,
	),
	runner.NewE2ETest(
		TestTONDepositRestrictedName,
		"deposit TON into MEVM restricted address",
		[]runner.ArgDefinition{
			{Description: "amount in nano tons", DefaultValue: "100000000"}, // 0.1 TON
		},
		TestTONDepositRestricted,
	),
	runner.NewE2ETest(
		TestTONWithdrawName,
		"withdraw TON from MEVM",
		[]runner.ArgDefinition{
			{Description: "amount in nano tons", DefaultValue: "2000000000"}, // 2.0 TON
		},
		TestTONWithdraw,
	),
	runner.NewE2ETest(
		TestTONWithdrawConcurrentName,
		"withdraw TON from MEVM for several recipients simultaneously",
		[]runner.ArgDefinition{},
		TestTONWithdrawConcurrent,
	),
	/*
	 Sui tests
	*/
	runner.NewE2ETest(
		TestSuiDepositName,
		"deposit SUI into MEVM",
		[]runner.ArgDefinition{
			{Description: "amount in mist", DefaultValue: "10000000000"},
		},
		TestSuiDeposit,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSuiDepositAndCallName,
		"deposit SUI into MEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in mist", DefaultValue: "1000000"},
		},
		TestSuiDepositAndCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSuiDepositAndCallRevertName,
		"deposit SUI into MEVM and call a contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount in mist", DefaultValue: "10000000000"},
		},
		TestSuiDepositAndCallRevert,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSuiTokenDepositName,
		"deposit fungible token SUI into MEVM",
		[]runner.ArgDefinition{
			{Description: "amount in base unit", DefaultValue: "10000000000"},
		},
		TestSuiTokenDeposit,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSuiTokenDepositAndCallName,
		"deposit fungible token into MEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in base unit", DefaultValue: "1000000"},
		},
		TestSuiTokenDepositAndCall,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSuiTokenDepositAndCallRevertName,
		"deposit fungible token into MEVM and call a contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount in base unit", DefaultValue: "10000000000"},
		},
		TestSuiTokenDepositAndCallRevert,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSuiWithdrawName,
		"withdraw SUI from MEVM",
		[]runner.ArgDefinition{
			{Description: "amount in mist", DefaultValue: "1000000"},
		},
		TestSuiWithdraw,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSuiWithdrawAndCallName,
		"withdraw SUI from MEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in mist", DefaultValue: "1000000"},
		},
		TestSuiWithdrawAndCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSuiWithdrawRevertWithCallName,
		"withdraw SUI from MEVM that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount in mist", DefaultValue: "1000000"},
		},
		TestSuiWithdrawRevertWithCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSuiWithdrawAndCallRevertWithCallName,
		"withdraw SUI from MEVM and call a contract that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount in mist", DefaultValue: "1000000"},
		},
		TestSuiWithdrawAndCallRevertWithCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSuiTokenWithdrawName,
		"withdraw fungible token from MEVM",
		[]runner.ArgDefinition{
			{Description: "amount in base unit", DefaultValue: "100000"},
		},
		TestSuiTokenWithdraw,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSuiTokenWithdrawAndCallName,
		"withdraw fungible token from MEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in base unit", DefaultValue: "100000"},
		},
		TestSuiTokenWithdrawAndCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSuiTokenWithdrawAndCallRevertWithCallName,
		"withdraw fungible token from MEVM and call a contract that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount in base unit", DefaultValue: "1000000"},
		},
		TestSuiTokenWithdrawAndCallRevertWithCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSuiDepositRestrictedName,
		"deposit SUI into MEVM restricted address",
		[]runner.ArgDefinition{
			{Description: "amount in mist", DefaultValue: "1000000"},
		},
		TestSuiDepositRestrictedAddress,
	),
	runner.NewE2ETest(
		TestSuiWithdrawRestrictedName,
		"withdraw SUI from MEVM to restricted address",
		[]runner.ArgDefinition{
			{Description: "receiver", DefaultValue: sample.RestrictedSuiAddressTest},
			{Description: "amount in mist", DefaultValue: "1000000"},
		},
		TestSuiWithdrawRestrictedAddress,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSuiWithdrawInvalidReceiverName,
		"withdraw SUI from MEVM to invalid receiver address",
		[]runner.ArgDefinition{
			{Description: "receiver", DefaultValue: "0x547a07f0564e0c8d48c4ae53305eabdef87e9610"},
			{Description: "amount in mist", DefaultValue: "1000000"},
		},
		TestSuiWithdrawInvalidReceiver,
	),
	/*
	 Bitcoin tests
	*/
	runner.NewE2ETest(
		TestBitcoinDonationName,
		"donate Bitcoin to TSS address", []runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.1"},
		},
		TestBitcoinDonation,
	),
	runner.NewE2ETest(
		TestBitcoinDepositName,
		"deposit Bitcoin into MEVM",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "1.0"},
		},
		TestBitcoinDeposit,
	),
	runner.NewE2ETest(
		TestBitcoinDepositFastConfirmationName,
		"deposit Bitcoin into MEVM using fast confirmation",
		[]runner.ArgDefinition{},
		TestBitcoinDepositFastConfirmation,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestBitcoinDepositAndCallName,
		"deposit Bitcoin into MEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinDepositAndCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestBitcoinDepositAndCallRevertName,
		"deposit Bitcoin into MEVM; expect refund",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.1"},
		},
		TestBitcoinDepositAndCallRevert,
	),
	runner.NewE2ETest(
		TestBitcoinDepositAndCallRevertWithDustName,
		"deposit Bitcoin into MEVM; revert with dust amount that aborts the CCTX",
		[]runner.ArgDefinition{},
		TestBitcoinDepositAndCallRevertWithDust,
	),
	runner.NewE2ETest(
		TestBitcoinDepositAndWithdrawWithDustName,
		"deposit Bitcoin into MEVM and withdraw with dust amount that fails the CCTX",
		[]runner.ArgDefinition{},
		TestBitcoinDepositAndWithdrawWithDust,
	),
	runner.NewE2ETest(
		TestBitcoinStdMemoDepositName,
		"deposit Bitcoin into MEVM with standard memo",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.2"},
		},
		TestBitcoinStdMemoDeposit,
	),
	runner.NewE2ETest(
		TestBitcoinStdMemoDepositAndCallName,
		"deposit Bitcoin into MEVM and call a contract with standard memo",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.5"},
		},
		TestBitcoinStdMemoDepositAndCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestBitcoinStdMemoDepositAndCallRevertName,
		"deposit Bitcoin into MEVM and call a contract with standard memo; expect revert",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.1"},
		},
		TestBitcoinStdMemoDepositAndCallRevert,
	),
	runner.NewE2ETest(
		TestBitcoinStdMemoDepositAndCallRevertOtherAddressName,
		"deposit Bitcoin into MEVM and call a contract with standard memo; expect revert to other address",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.1"},
		},
		TestBitcoinStdMemoDepositAndCallRevertOtherAddress,
	),
	runner.NewE2ETest(
		TestBitcoinStdMemoDepositAndCallRevertAndAbortName,
		"deposit Bitcoin into MEVM and call a contract with standard memo; revert and abort with onAbort",
		[]runner.ArgDefinition{},
		TestBitcoinStdMemoDepositAndCallRevertAndAbort,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestBitcoinStdMemoInscribedDepositAndCallName,
		"deposit Bitcoin into MEVM and call a contract with inscribed standard memo",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.1"},
			{Description: "fee rate", DefaultValue: "10"},
		},
		TestBitcoinStdMemoInscribedDepositAndCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestBitcoinDepositAndAbortWithLowDepositFeeName,
		"deposit Bitcoin into MEVM that aborts due to insufficient deposit fee",
		[]runner.ArgDefinition{},
		TestBitcoinDepositAndAbortWithLowDepositFee,
		runner.WithMinimumVersion("v27.0.0"),
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawSegWitName,
		"withdraw BTC from MEVM to a SegWit address",
		[]runner.ArgDefinition{
			{Description: "receiver address", DefaultValue: ""},
			{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinWithdrawSegWit,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawTaprootName,
		"withdraw BTC from MEVM to a Taproot address",
		[]runner.ArgDefinition{
			{Description: "receiver address", DefaultValue: ""},
			{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinWithdrawTaproot,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawLegacyName,
		"withdraw BTC from MEVM to a legacy address",
		[]runner.ArgDefinition{
			{Description: "receiver address", DefaultValue: ""},
			{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinWithdrawLegacy,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawMultipleName,
		"withdraw BTC from MEVM multiple times",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "0.01"},
			{Description: "times", DefaultValue: "2"},
		},
		WithdrawBitcoinMultipleTimes,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawP2WSHName,
		"withdraw BTC from MEVM to a P2WSH address",
		[]runner.ArgDefinition{
			{Description: "receiver address", DefaultValue: ""},
			{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinWithdrawP2WSH,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawP2SHName,
		"withdraw BTC from MEVM to a P2SH address",
		[]runner.ArgDefinition{
			{Description: "receiver address", DefaultValue: ""},
			{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinWithdrawP2SH,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawInvalidAddressName,
		"withdraw BTC from MEVM to an unsupported btc address",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.00001"},
		},
		TestBitcoinWithdrawToInvalidAddress,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawRestrictedName,
		"withdraw Bitcoin from MEVM to restricted address",
		[]runner.ArgDefinition{
			{Description: "receiver", DefaultValue: sample.RestrictedBtcAddressTest},
			{Description: "amount in btc", DefaultValue: "0.001"},
			{Description: "revert address", DefaultValue: sample.RevertAddressMEVM},
		},
		TestBitcoinWithdrawRestricted,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestBitcoinDepositInvalidMemoRevertName,
		"deposit Bitcoin with invalid memo; expect revert",
		[]runner.ArgDefinition{},
		TestBitcoinDepositInvalidMemoRevert,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawRBFName,
		"withdraw Bitcoin from MEVM and replace the outbound using RBF",
		[]runner.ArgDefinition{
			{Description: "receiver address", DefaultValue: ""},
			{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinWithdrawRBF,
		runner.WithDependencies(DepdencyAllBitcoinDeposits),
	),
	/*
	 Application tests
	*/
	runner.NewE2ETest(
		TestCrosschainSwapName,
		"testing Bitcoin ERC20 cross-chain swap",
		[]runner.ArgDefinition{},
		TestCrosschainSwap,
	),
	/*
	 Miscellaneous tests
	*/
	runner.NewE2ETest(
		TestDonationEtherName,
		"donate Ether to the TSS",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000000000000000"},
		},
		TestDonationEther,
	),
	runner.NewE2ETest(
		TestInboundTrackersName,
		"test processing inbound trackers for observation",
		[]runner.ArgDefinition{},
		TestInboundTrackers,
	),
	/*
	 Stress tests
	*/
	runner.NewE2ETest(
		TestStressEtherWithdrawName,
		"stress test Ether withdrawal",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
			{Description: CountArgDescription, DefaultValue: "100"},
		},
		TestStressEtherWithdraw,
	),
	runner.NewE2ETest(
		TestStressBTCWithdrawName,
		"stress test BTC withdrawal",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.01"},
			{Description: CountArgDescription, DefaultValue: "100"},
		},
		TestStressBTCWithdraw,
	),
	runner.NewE2ETest(
		TestStressEtherDepositName,
		"stress test Ether deposit",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
			{Description: CountArgDescription, DefaultValue: "100"},
		},
		TestStressEtherDeposit,
	),
	runner.NewE2ETest(
		TestStressBTCDepositName,
		"stress test BTC deposit",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.001"},
			{Description: CountArgDescription, DefaultValue: "100"},
		},
		TestStressBTCDeposit,
	),
	runner.NewE2ETest(
		TestStressSolanaDepositName,
		"stress test SOL deposit",
		[]runner.ArgDefinition{
			{Description: "amount in lamports", DefaultValue: "1200000"},
			{Description: CountArgDescription, DefaultValue: "50"},
		},
		TestStressSolanaDeposit,
	),
	runner.NewE2ETest(
		TestStressSPLDepositName,
		"stress test SPL deposit",
		[]runner.ArgDefinition{
			{Description: "amount in SPL tokens", DefaultValue: "1200000"},
			{Description: CountArgDescription, DefaultValue: "50"},
		},
		TestStressSPLDeposit,
	),
	runner.NewE2ETest(
		TestStressSolanaWithdrawName,
		"stress test SOL withdrawals",
		[]runner.ArgDefinition{
			{Description: "amount in lamports", DefaultValue: "1000000"},
			{Description: CountArgDescription, DefaultValue: "50"},
		},
		TestStressSolanaWithdraw,
	),
	runner.NewE2ETest(
		TestStressSPLWithdrawName,
		"stress test SPL withdrawals",
		[]runner.ArgDefinition{
			{Description: "amount in SPL tokens", DefaultValue: "1000000"},
			{Description: CountArgDescription, DefaultValue: "50"},
		},
		TestStressSPLWithdraw,
	),
	runner.NewE2ETest(
		TestStressSuiDepositName,
		"stress test SUI deposits",
		[]runner.ArgDefinition{
			{Description: "amount in SUI", DefaultValue: "1000000"},
			{Description: CountArgDescription, DefaultValue: "50"},
		},
		TestStressSuiDeposit,
	),
	runner.NewE2ETest(
		TestStressSuiWithdrawName,
		"stress test SUI withdrawals",
		[]runner.ArgDefinition{
			{Description: "amount in SUI", DefaultValue: "1000000"},
			{Description: CountArgDescription, DefaultValue: "50"},
		},
		TestStressSuiWithdraw,
	),
	/*
	 Admin tests
	*/
	runner.NewE2ETest(
		TestWhitelistERC20Name,
		"whitelist a new ERC20 token",
		[]runner.ArgDefinition{},
		TestWhitelistERC20,
	),
	runner.NewE2ETest(
		TestDepositEtherLiquidityCapName,
		"deposit Ethers into MEVM with a liquidity cap",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000000000000"},
		},
		TestDepositEtherLiquidityCap,
	),
	runner.NewE2ETest(
		TestMigrateChainSupportName,
		"migrate the evm chain from goerli to sepolia",
		[]runner.ArgDefinition{},
		TestMigrateChainSupport,
	),
	runner.NewE2ETest(
		TestPauseMRC20Name,
		"pausing MRC20 on MuseChain",
		[]runner.ArgDefinition{},
		TestPauseMRC20,
	),
	runner.NewE2ETest(
		TestUpdateBytecodeMRC20Name,
		"update MRC20 bytecode swap",
		[]runner.ArgDefinition{},
		TestUpdateBytecodeMRC20,
	),
	runner.NewE2ETest(
		TestUpdateBytecodeConnectorName,
		"update mevm connector bytecode",
		[]runner.ArgDefinition{},
		TestUpdateBytecodeConnector,
	),
	runner.NewE2ETest(
		TestRateLimiterName,
		"test sending cctxs with rate limiter enabled and show logs when processing cctxs",
		[]runner.ArgDefinition{},
		legacy.TestRateLimiter,
	),
	runner.NewE2ETest(
		TestCriticalAdminTransactionsName,
		"test critical admin transactions",
		[]runner.ArgDefinition{},
		TestCriticalAdminTransactions,
	),
	runner.NewE2ETest(
		TestMigrateTSSName,
		"migrate TSS funds",
		[]runner.ArgDefinition{},
		TestMigrateTSS,
	),
	runner.NewE2ETest(
		TestPauseERC20CustodyName,
		"pausing ERC20 custody on MuseChain",
		[]runner.ArgDefinition{},
		TestPauseERC20Custody,
	),
	runner.NewE2ETest(
		TestMigrateERC20CustodyFundsName,
		"migrate ERC20 custody funds",
		[]runner.ArgDefinition{},
		TestMigrateERC20CustodyFunds,
	),
	runner.NewE2ETest(
		TestUpdateMRC20NameName,
		"update MRC20 name and symbol",
		[]runner.ArgDefinition{},
		TestUpdateMRC20Name,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestMuseclientRestartHeightName,
		"museclient scheduled restart height",
		[]runner.ArgDefinition{},
		TestMuseclientRestartHeight,
	),
	runner.NewE2ETest(
		TestMuseclientSignerOffsetName,
		"museclient signer offset",
		[]runner.ArgDefinition{},
		TestMuseclientSignerOffset,
	),
	runner.NewE2ETest(
		TestUpdateOperationalChainParamsName,
		"update operational chain params",
		[]runner.ArgDefinition{},
		TestUpdateOperationalChainParams,
		runner.WithMinimumVersion("v29.0.0"),
	),
	/*
	 Special tests
	*/
	runner.NewE2ETest(
		TestDeploy,
		"deploy a contract",
		[]runner.ArgDefinition{
			{Description: "contract name", DefaultValue: ""},
		},
		TestDeployContract,
	),
	runner.NewE2ETest(
		TestOperationAddLiquidityETHName,
		"add liquidity to the MUSE/ETH pool",
		[]runner.ArgDefinition{
			{Description: "amountMUSE", DefaultValue: "50000000000000000000"},
			{Description: "amountETH", DefaultValue: "50000000000000000000"},
		},
		TestOperationAddLiquidityETH,
	),
	runner.NewE2ETest(
		TestOperationAddLiquidityERC20Name,
		"add liquidity to the MUSE/ERC20 pool",
		[]runner.ArgDefinition{
			{Description: "amountMUSE", DefaultValue: "50000000000000000000"},
			{Description: "amountERC20", DefaultValue: "50000000000000000000"},
		},
		TestOperationAddLiquidityERC20,
	),
	runner.NewE2ETest(
		TestOperationAddLiquidityBTCName,
		"add liquidity to the MUSE/BTC pool",
		[]runner.ArgDefinition{
			{Description: "amountMUSE", DefaultValue: "50000000000000000000"},
			{Description: "amountBTC", DefaultValue: "5000000000"},
		},
		TestOperationAddLiquidityBTC,
	),
	runner.NewE2ETest(
		TestOperationAddLiquiditySOLName,
		"add liquidity to the MUSE/SOL pool",
		[]runner.ArgDefinition{
			{Description: "amountMUSE", DefaultValue: "50000000000000000000"},
			{Description: "amountSOL", DefaultValue: "50000000000"},
		},
		TestOperationAddLiquiditySOL,
	),
	runner.NewE2ETest(
		TestOperationAddLiquiditySPLName,
		"add liquidity to the MUSE/SPL pool",
		[]runner.ArgDefinition{
			{Description: "amountMUSE", DefaultValue: "50000000000000000000"},
			{Description: "amountSPL", DefaultValue: "50000000000000000000"},
		},
		TestOperationAddLiquiditySPL,
	),
	runner.NewE2ETest(
		TestOperationAddLiquiditySUIName,
		"add liquidity to the MUSE/SUI pool",
		[]runner.ArgDefinition{
			{Description: "amountMUSE", DefaultValue: "50000000000000000000"},
			{Description: "amountSUI", DefaultValue: "50000000000"},
		},
		TestOperationAddLiquiditySUI,
	),
	runner.NewE2ETest(
		TestOperationAddLiquiditySuiFungibleTokenName,
		"add liquidity to the MUSE/SuiFungibleToken pool",
		[]runner.ArgDefinition{
			{Description: "amountMUSE", DefaultValue: "50000000000000000000"},
			{Description: "amountSuiFungibleToken", DefaultValue: "50000000"},
		},
		TestOperationAddLiquiditySuiFungibleToken,
	),
	runner.NewE2ETest(
		TestOperationAddLiquidityTONName,
		"add liquidity to the MUSE/TON pool",
		[]runner.ArgDefinition{
			{Description: "amountMUSE", DefaultValue: "50000000000000000000"},
			{Description: "amountTON", DefaultValue: "50000000000"},
		},
		TestOperationAddLiquidityTON,
	),
	/*
	 Stateful precompiled contracts tests
	*/
	runner.NewE2ETest(
		TestPrecompilesPrototypeName,
		"test stateful precompiled contracts prototype",
		[]runner.ArgDefinition{},
		TestPrecompilesPrototype,
	),
	runner.NewE2ETest(
		TestPrecompilesPrototypeThroughContractName,
		"test stateful precompiled contracts prototype through contract",
		[]runner.ArgDefinition{},
		TestPrecompilesPrototypeThroughContract,
	),
	runner.NewE2ETest(
		TestPrecompilesStakingName,
		"test stateful precompiled contracts staking",
		[]runner.ArgDefinition{},
		TestPrecompilesStakingIsDisabled,
	),
	runner.NewE2ETest(
		TestPrecompilesStakingThroughContractName,
		"test stateful precompiled contracts staking through contract",
		[]runner.ArgDefinition{},
		TestPrecompilesStakingThroughContract,
	),
	runner.NewE2ETest(
		TestPrecompilesBankName,
		"test stateful precompiled contracts bank with MRC20 tokens",
		[]runner.ArgDefinition{},
		TestPrecompilesBank,
	),
	runner.NewE2ETest(
		TestPrecompilesBankFailName,
		"test stateful precompiled contracts bank with non MRC20 tokens",
		[]runner.ArgDefinition{},
		TestPrecompilesBankNonMRC20,
	),
	runner.NewE2ETest(
		TestPrecompilesBankThroughContractName,
		"test stateful precompiled contracts bank through contract",
		[]runner.ArgDefinition{},
		TestPrecompilesBankThroughContract,
	),
	runner.NewE2ETest(
		TestPrecompilesDistributeName,
		"test stateful precompiled contracts distribute",
		[]runner.ArgDefinition{},
		TestPrecompilesDistributeAndClaim,
	),
	runner.NewE2ETest(
		TestPrecompilesDistributeNonMRC20Name,
		"test stateful precompiled contracts distribute with non MRC20 tokens",
		[]runner.ArgDefinition{},
		TestPrecompilesDistributeNonMRC20,
	),
	runner.NewE2ETest(
		TestPrecompilesDistributeThroughContractName,
		"test stateful precompiled contracts distribute through contract",
		[]runner.ArgDefinition{},
		TestPrecompilesDistributeAndClaimThroughContract,
	),

	/*
	 Legacy tests
	*/
	runner.NewE2ETest(
		TestLegacyMessagePassingExternalChainsName,
		"evm->evm message passing (sending MUSE only) (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in amuse", DefaultValue: "10000000000000000000"},
		},
		legacy.TestMessagePassingExternalChains,
	),
	runner.NewE2ETest(
		TestLegacyMessagePassingRevertFailExternalChainsName,
		"message passing with failing revert between external EVM chains (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in amuse", DefaultValue: "10000000000000000000"},
		},
		legacy.TestMessagePassingRevertFailExternalChains,
	),
	runner.NewE2ETest(
		TestLegacyMessagePassingRevertSuccessExternalChainsName,
		"message passing with successful revert between external EVM chains (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in amuse", DefaultValue: "10000000000000000000"},
		},
		legacy.TestMessagePassingRevertSuccessExternalChains,
	),
	runner.NewE2ETest(
		TestLegacyMessagePassingEVMtoMEVMName,
		"evm -> mevm message passing contract call (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in amuse", DefaultValue: "10000000000000000009"},
		},
		legacy.TestMessagePassingEVMtoMEVM,
	),
	runner.NewE2ETest(
		TestLegacyMessagePassingMEVMToEVMName,
		"mevm -> evm message passing contract call (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in amuse", DefaultValue: "10000000000000000007"},
		},
		legacy.TestMessagePassingMEVMtoEVM,
	),
	runner.NewE2ETest(
		TestLegacyMessagePassingMEVMtoEVMRevertName,
		"mevm -> evm message passing contract call reverts (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in amuse", DefaultValue: "10000000000000000006"},
		},
		legacy.TestMessagePassingMEVMtoEVMRevert,
	),
	runner.NewE2ETest(
		TestLegacyMessagePassingEVMtoMEVMRevertName,
		"evm -> mevm message passing and revert back to evm (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in amuse", DefaultValue: "10000000000000000008"},
		},
		legacy.TestMessagePassingEVMtoMEVMRevert,
	),
	runner.NewE2ETest(
		TestLegacyMessagePassingMEVMtoEVMRevertFailName,
		"mevm -> evm message passing contract with failing revert (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in amuse", DefaultValue: "10000000000000000008"},
		},
		legacy.TestMessagePassingMEVMtoEVMRevertFail,
	),
	runner.NewE2ETest(
		TestLegacyMessagePassingEVMtoMEVMRevertFailName,
		"evm -> mevm message passing contract with failing revert (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in amuse", DefaultValue: "10000000000000000008"},
		},
		legacy.TestMessagePassingEVMtoMEVMRevertFail,
	),
	runner.NewE2ETest(
		TestLegacyEtherDepositName,
		"deposit Ether into MEVM (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		legacy.TestEtherDeposit,
	),
	runner.NewE2ETest(
		TestLegacyEtherWithdrawName,
		"withdraw Ether from MEVM (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		legacy.TestEtherWithdraw,
	),
	runner.NewE2ETest(
		TestEtherWithdrawRestrictedName,
		"withdraw Ether from MEVM to restricted address (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "receiver", DefaultValue: sample.RestrictedEVMAddressTest},
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestEtherWithdrawRestricted,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestLegacyEtherDepositAndCallRefundName,
		"deposit Ether into MEVM and call a contract that reverts; should refund (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000000"},
		},
		legacy.TestEtherDepositAndCallRefund,
	),
	runner.NewE2ETest(
		TestLegacyEtherDepositAndCallName,
		"deposit MRC20 into MEVM and call a contract (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "1000000000000000000"},
		},
		legacy.TestEtherDepositAndCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestLegacyERC20WithdrawName,
		"withdraw ERC20 from MEVM (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		legacy.TestERC20Withdraw,
	),
	runner.NewE2ETest(
		TestLegacyERC20DepositName,
		"deposit ERC20 into MEVM (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100000"},
		},
		legacy.TestERC20Deposit,
	),
	runner.NewE2ETest(
		TestLegacyMultipleERC20DepositName,
		"deposit ERC20 into MEVM in multiple deposits (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000000000"},
			{Description: CountArgDescription, DefaultValue: "3"},
		},
		legacy.TestMultipleERC20Deposit,
	),
	runner.NewE2ETest(
		TestLegacyMultipleERC20WithdrawsName,
		"withdraw ERC20 from MEVM in multiple withdrawals (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100"},
			{Description: CountArgDescription, DefaultValue: "3"},
		},
		legacy.TestMultipleERC20Withdraws,
	),
	runner.NewE2ETest(
		TestERC20DepositRestrictedName,
		"deposit ERC20 into MEVM restricted address (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100000"},
		},
		TestERC20DepositRestricted,
	),
	runner.NewE2ETest(
		TestLegacyERC20DepositAndCallRefundName,
		"deposit a non-gas MRC20 into MEVM and call a contract that reverts (v1 protocol contracts)",
		[]runner.ArgDefinition{},
		legacy.TestERC20DepositAndCallRefund,
	),

	/*
	 MUSE tests
	*/
	runner.NewE2ETest(
		TestLegacyMuseDepositName,
		"deposit MUSE from Ethereum to MEVM",
		[]runner.ArgDefinition{
			{Description: "amount in amuse", DefaultValue: "1000000000000000000"},
		},
		legacy.TestMuseDeposit,
	),
	runner.NewE2ETest(
		TestLegacyMuseDepositNewAddressName,
		"deposit MUSE from Ethereum to a new MEVM address which does not exist yet",
		[]runner.ArgDefinition{
			{Description: "amount in amuse", DefaultValue: "1000000000000000000"},
		},
		legacy.TestMuseDepositNewAddress,
	),
	runner.NewE2ETest(
		TestLegacyMuseDepositRestrictedName,
		"deposit MUSE from Ethereum to MEVM restricted address",
		[]runner.ArgDefinition{
			{Description: "amount in amuse", DefaultValue: "1000000000000000000"},
		},
		legacy.TestMuseDepositRestricted,
	),
	runner.NewE2ETest(
		TestLegacyMuseWithdrawName,
		"withdraw MUSE from MEVM to Ethereum",
		[]runner.ArgDefinition{
			{Description: "amount in amuse", DefaultValue: "10000000000000000000"},
		},
		legacy.TestMuseWithdraw,
	),
	runner.NewE2ETest(
		TestLegacyMuseWithdrawBTCRevertName,
		"sending MUSE from MEVM to Bitcoin with a message that should revert cctxs",
		[]runner.ArgDefinition{
			{Description: "amount in amuse", DefaultValue: "1000000000000000000"},
		},
		legacy.TestMuseWithdrawBTCRevert,
	),
	runner.NewE2ETest(
		TestUndelegateToBelowMinimumObserverDelegation,
		"test undelegating to below minimum observer delegation",
		[]runner.ArgDefinition{},
		UndelegateToBelowMinimumObserverDelegation),
}
