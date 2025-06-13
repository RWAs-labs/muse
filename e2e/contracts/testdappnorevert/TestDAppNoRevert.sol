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

interface MuseConnector {
    /**
     * @dev Sending value and data cross-chain is as easy as calling connector.send(SendInput)
     */
    function send(MuseInterfaces.SendInput calldata input) external;
}

interface IERC20 {
    function transferFrom(address _from, address _to, uint256 _value) external returns (bool success);
    function approve(address _spender, uint256 _value) external returns (bool success);
}

/**
 * @dev TestDAppNoRevert is a test dapp not implementing the interface for revert to test revert failure cases in message passing
*/
contract TestDAppNoRevert {
    bytes32 public constant HELLO_WORLD_MESSAGE_TYPE = keccak256("CROSS_CHAIN_HELLO_WORLD");
    event HelloWorldEvent();
    event RevertedHelloWorldEvent();
    error InvalidMessageType();
    error ErrorTransferringMuse();
    address public connector;
    address public muse;
    constructor(address _connector, address _museToken) {
        connector = _connector;
        muse = _museToken;
    }

    function onMuseMessage(MuseInterfaces.MuseMessage calldata museMessage) external {
        (, bool doRevert) = abi.decode(museMessage.message, (bytes32, bool));
        require(doRevert == false,  "message says revert");

        emit HelloWorldEvent();
    }

    function sendHelloWorld(address destinationAddress, uint256 destinationChainId, uint256 value, bool doRevert) external payable {
        bool success1 = IERC20(muse).approve(address(connector), value);
        bool success2 = IERC20(muse).transferFrom(msg.sender, address(this), value);
        if (!(success1 && success2)) revert ErrorTransferringMuse();

        MuseConnector(connector).send(
            MuseInterfaces.SendInput({
                destinationChainId: destinationChainId,
                destinationAddress: abi.encodePacked(destinationAddress),
                destinationGasLimit: 250000,
                message: abi.encode(HELLO_WORLD_MESSAGE_TYPE, doRevert),
                museValueAndGas: value,
                museParams: abi.encode("")
            })
        );
    }
}