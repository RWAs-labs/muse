// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

interface MuseInterfaces {
    /**
     * @dev Use SendInput to interact with the Connector: connector.send(SendInput)
     */
    struct SendInput {
        /// @dev Chain id of the destination chain. More about chain ids https://docs.musechain.com/learn/glossary#chain-id
        uint256 destinationChainId;
        /// @dev Address receiving the message on the destination chain (expressed in bytes since it can be non-EVM)
        bytes destinationAddress;
        /// @dev Gas limit for the destination chain's transaction
        uint256 destinationGasLimit;
        /// @dev An encoded, arbitrary message to be parsed by the destination contract
        bytes message;
        /// @dev MUSE to be sent cross-chain + MuseChain gas fees + destination chain gas fees (expressed in MUSE)
        uint256 museValueAndGas;
        /// @dev Optional parameters for the MuseChain protocol
        bytes museParams;
    }

    /**
     * @dev Our Connector calls onMuseMessage with this struct as argument
     */
    struct MuseMessage {
        bytes museTxSenderAddress;
        uint256 sourceChainId;
        address destinationAddress;
        /// @dev Remaining MUSE from museValueAndGas after subtracting MuseChain gas fees and destination gas fees
        uint256 museValue;
        bytes message;
    }

    /**
     * @dev Our Connector calls onMuseRevert with this struct as argument
     */
    struct MuseRevert {
        address museTxSenderAddress;
        uint256 sourceChainId;
        bytes destinationAddress;
        uint256 destinationChainId;
        /// @dev Equals to: museValueAndGas - MuseChain gas fees - destination chain gas fees - source chain revert tx gas fees
        uint256 remainingMuseValue;
        bytes message;
    }
}

// Dapp is a sample comtract that implements MuseReceiver and is used for unit testing
// It sets the values of the MuseMessage struct to its public variables which can then be queried to check if the function was called correctly
contract Dapp {
    bytes public museTxSenderAddress;
    uint256 public sourceChainId;
    address public destinationAddress;
    uint256 public destinationChainId;
    uint256 public museValue;
    bytes public  message;

    constructor() {
        museTxSenderAddress = "";
        sourceChainId = 0;
        destinationAddress = address(0);
        destinationChainId = 0;
        museValue = 0;
        message = "";
    }

    function onMuseMessage(MuseInterfaces.MuseMessage calldata museMessage) external{
        museTxSenderAddress = museMessage.museTxSenderAddress;
        sourceChainId = museMessage.sourceChainId;
        destinationAddress = museMessage.destinationAddress;
        museValue = museMessage.museValue;
        message = museMessage.message;
    }
    function onMuseRevert(MuseInterfaces.MuseRevert calldata museRevert) external {
        museTxSenderAddress = abi.encodePacked(museRevert.museTxSenderAddress);
        sourceChainId = museRevert.sourceChainId;
        destinationAddress = address(uint160(uint256(keccak256(museRevert.destinationAddress))));
        destinationChainId = museRevert.destinationChainId;
        museValue = museRevert.remainingMuseValue;
        message = museRevert.message;
    }
}