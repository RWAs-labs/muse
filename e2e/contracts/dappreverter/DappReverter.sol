// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

// DappReverter is a contract that can be used to test the reversion of a cross-chain call.
// It implements the onMuseMessage and onMuseRevert functions, which are called the MEVM connector
contract DappReverter {
    function onMuseMessage() external{}
    function onMuseRevert() external {}
}