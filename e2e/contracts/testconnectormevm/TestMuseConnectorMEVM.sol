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

interface WMUSE {
    function transferFrom(address src, address dst, uint wad) external returns (bool);

    function withdraw(uint wad) external;
}

contract TestMuseConnectorMEVM is MuseInterfaces {
    /// @notice Contract custom errors.
    error OnlyWMUSE();
    error WMUSETransferFailed();
    error OnlyFungibleModule();
    error FailedMuseSent();

    /// @notice Fungible module address.
    address public constant FUNGIBLE_MODULE_ADDRESS = payable(0x735b14BB79463307AAcBED86DAf3322B1e6226aB);
    /// @notice WMUSE token address.
    address public wmuse;

    event MuseSent(
        address sourceTxOriginAddress,
        address indexed museTxSenderAddress,
        uint256 indexed destinationChainId,
        bytes destinationAddress,
        uint256 museValueAndGas,
        uint256 destinationGasLimit,
        bytes message,
        bytes museParams
    );
    event SetWMUSE(address wmuse_);

    constructor(address wmuse_) {
        wmuse = wmuse_;
    }

    /// @dev Receive function to receive MUSE from WETH9.withdraw().
    receive() external payable {
        if (msg.sender != wmuse) revert OnlyWMUSE();
    }

    /**
     * @dev Sends MUSE and bytes messages (to execute it) crosschain.
     * @param input, SendInput struct, checkout above.
     */
    function send(MuseInterfaces.SendInput calldata input) external {
        // Transfer wmuse to "fungible" module, which will be burnt by the protocol post processing via hooks.
        if (!WMUSE(wmuse).transferFrom(msg.sender, address(this), input.museValueAndGas)) revert WMUSETransferFailed();
        WMUSE(wmuse).withdraw(input.museValueAndGas);
        (bool sent, ) = FUNGIBLE_MODULE_ADDRESS.call{value: input.museValueAndGas}("");
        if (!sent) revert FailedMuseSent();
        emit MuseSent(
            tx.origin,
            msg.sender,
            input.destinationChainId,
            input.destinationAddress,
            input.museValueAndGas,
            input.destinationGasLimit,
            input.message,
            input.museParams
        );
    }

    /**
     * @dev Sends MUSE and bytes messages (to execute it) crosschain.
     * @param wmuse_, new WMUSE address.
     */
    function setWmuseAddress(address wmuse_) external {
        if (msg.sender != FUNGIBLE_MODULE_ADDRESS) revert OnlyFungibleModule();
        wmuse = wmuse_;
        emit SetWMUSE(wmuse_);
    }

    function foo() external pure returns (string memory) {
        return "foo";
    }
}