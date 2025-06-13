## Full Deployment Guide for Musecored

This guide details deploying Musecored nodes on both MuseChain mainnet and Athens3 (testnet). The setup utilizes Docker Compose with environment variables for a streamlined deployment process.

Here's a comprehensive documentation using markdown tables to cover all the `make` commands for managing Musecored, including where to modify the environment variables in Docker Compose configurations.

### Musecored / BTC Node Deployment and Management

#### Commands Overview for Musecored

| Environment                         | Action                      | Command                                                       | Docker Compose Location                  |
|-------------------------------------|-----------------------------|---------------------------------------------------------------|------------------------------------------|
| **Mainnet**                         | Start Ethereum Node         | `make start-eth-node-mainnet`                                 | `contrib/rpc/ethereum`                   |
| **Mainnet**                         | Stop Ethereum Node          | `make stop-eth-node-mainnet`                                  | `contrib/rpc/ethereum`                   |
| **Mainnet**                         | Clean Ethereum Node Data    | `make clean-eth-node-mainnet`                                 | `contrib/rpc/ethereum`                   |
| **Mainnet**                         | Start Musecored Node        | `make start-mainnet-muserpc-node DOCKER_TAG=ubuntu-v14.0.1`   | `contrib/rpc/musecored`                  |
| **Mainnet**                         | Stop Musecored Node         | `make stop-mainnet-muserpc-node`                              | `contrib/rpc/musecored`                  |
| **Mainnet**                         | Clean Musecored Node Data   | `make clean-mainnet-muserpc-node`                             | `contrib/rpc/musecored`                  |
| **Testnet (Athens3)**               | Start Musecored Node        | `make start-testnet-muserpc-node DOCKER_TAG=ubuntu-v14.0.1`   | `contrib/rpc/musecored`                  |
| **Testnet (Athens3)**               | Stop Musecored Node         | `make stop-testnet-muserpc-node`                              | `contrib/rpc/musecored`                  |
| **Testnet (Athens3)**               | Clean Musecored Node Data   | `make clean-testnet-muserpc-node`                             | `contrib/rpc/musecored`                  |
| **Mainnet Local Build**             | Start Musecored Node        | `make start-musecored-rpc-mainnet-localbuild`                 | `contrib/rpc/musecored`                  |
| **Mainnet Local Build**             | Stop Musecored Node         | `make stop-musecored-rpc-mainnet-localbuild`                  | `contrib/rpc/musecored`                  |
| **Mainnet Local Build**             | Clean Musecored Node Data   | `make clean-musecored-rpc-mainnet-localbuild`                 | `contrib/rpc/musecored`                  |
| **Testnet Local Build (Athens3)**   | Start Musecored Node        | `make start-musecored-rpc-testnet-localbuild`                 | `contrib/rpc/musecored`                  |
| **Testnet Local Build (Athens3)**   | Stop Musecored Node         | `make stop-musecored-rpc-testnet-localbuild`                  | `contrib/rpc/musecored`                  |
| **Testnet Local Build (Athens3)**   | Clean Musecored Node Data   | `make clean-musecored-rpc-testnet-localbuild`                 | `contrib/rpc/musecored`                  |
