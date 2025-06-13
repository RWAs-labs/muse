#!/bin/bash

# This script is used to start the musecored nodes
# It initializes the nodes and creates the genesis.json file
# It also starts the nodes
# The number of nodes is passed as an first argument to the script
# The second argument is optional and can have the following value:
#  - import-data: import data into the genesis file

/usr/sbin/sshd

add_emissions_withdraw_authorizations() {

    config_file="/root/config.yml"
    json_file="/root/.musecored/config/genesis.json"

    # Check if config file exists
    if [[ ! -f "$config_file" ]]; then
        echo "Error: Config file not found at $config_file"
        return 1
    fi
    # Address to add emissions withdraw authorizations
    address=$(yq -r '.additional_accounts.user_emissions_withdraw.bech32_address' "$config_file")

    # Check if genesis file exists
    if [[ ! -f "$json_file" ]]; then
        echo "Error: Genesis file not found at $json_file"
        return 1
    fi

    echo "Adding emissions withdraw authorizations for address: $address"


     # Using jq to parse JSON, create new entries, and append them to the authorization array
     if ! jq --arg address "$address" '
         # Store the nodeAccountList array
         .app_state.observer.nodeAccountList as $list |
         # Iterate over the stored list to construct new objects and append to the authorization array
         .app_state.authz.authorization += [
             $list[] |
             {
                 "granter": .operator,
                 "grantee": $address,
                 "authorization": {
                     "@type": "/cosmos.authz.v1beta1.GenericAuthorization",
                     "msg": "/musechain.musecore.emissions.MsgWithdrawEmission"
                 },
                 "expiration": null
             }
         ]
     ' "$json_file" > temp.json; then
         echo "Error: Failed to update genesis file"
         return 1
     fi
     mv temp.json "$json_file"
}

# 10 million muse
DEFAULT_FUND_AMOUNT="10000000000000000000000000amuse"

# Funds an individual account
fund_account() {
  local name=$1
  local account=$2
  local amount=${3:-$DEFAULT_FUND_AMOUNT}

  echo "Funding $name ($account) with $amount"

  musecored add-genesis-account "$account" "$amount"
}

# Funds most accounts automatically
fund_accounts_auto() {
  # Fund the default account first
  local default_address=$(yq -r '.default_account.bech32_address' /root/config.yml)
  fund_account "default_account" "$default_address"

  # Get all additional accounts and fund them
  local accounts=$(yq -r '.additional_accounts | keys | sort | .[]' /root/config.yml)
  for account_key in $accounts; do
    local address=$(yq -r ".additional_accounts.$account_key.bech32_address" /root/config.yml)
    fund_account "$account_key" "$address"
  done
}

# create keys
CHAINID="athens_101-1"
KEYRING="test"
HOSTNAME=$(hostname)

if [[ $HOSTNAME == "musecore-new-validator" ]]; then
    INDEX="-new-validator"
else
    INDEX=${HOSTNAME:0-1}
fi

echo "HOSTNAME: $HOSTNAME, INDEX: $INDEX"

# Environment variables used for upgrade testing
export DAEMON_HOME=$HOME/.musecored
export DAEMON_NAME=musecored
export DAEMON_ALLOW_DOWNLOAD_BINARIES=true
export DAEMON_RESTART_AFTER_UPGRADE=true
export CLIENT_DAEMON_NAME=museclientd
export CLIENT_DAEMON_ARGS="-enable-chains,GOERLI,-val,operator"
export DAEMON_DATA_BACKUP_DIR=$DAEMON_HOME
export CLIENT_SKIP_UPGRADE=true
export CLIENT_START_PROCESS=false
export UNSAFE_SKIP_BACKUP=true

# init ssh keys
# we generate keys at runtime to ensure that keys are never pushed to
# a docker registry
if [ $HOSTNAME == "musecore0" ]; then
  if [[ ! -f ~/.ssh/id_rsa ]]; then
    ssh-keygen -t rsa -q -N "" -f ~/.ssh/id_rsa
    cp ~/.ssh/id_rsa.pub ~/.ssh/authorized_keys
    # keep localtest.pem for compatibility
    cp ~/.ssh/id_rsa ~/.ssh/localtest.pem
    chmod 600 ~/.ssh/*
  fi
fi

# Wait for authorized_keys file to exist (musecore1+)
while [ ! -f ~/.ssh/authorized_keys ]; do
    echo "Waiting for authorized_keys file to exist..."
    sleep 1
done

# Skip init if it has already been completed (marked by presence of ~/.musecored/init_complete file)
if [[ ! -f ~/.musecored/init_complete ]]
then
  # Init a new node to generate genesis file .
  # Copy config files from existing folders which get copied via Docker Copy when building images
  mkdir -p ~/.backup/config
  musecored init Musenode-Localnet --chain-id=$CHAINID --default-denom amuse
  rm -rf ~/.musecored/config/app.toml
  rm -rf ~/.musecored/config/client.toml
  rm -rf ~/.musecored/config/config.toml
  cp -r ~/musecored/common/app.toml ~/.musecored/config/
  cp -r ~/musecored/common/client.toml ~/.musecored/config/
  cp -r ~/musecored/common/config.toml ~/.musecored/config/
  sed -i -e "/moniker =/s/=.*/= \"$HOSTNAME\"/" "$HOME"/.musecored/config/config.toml
fi

echo "Creating keys for operator and hotkey for $HOSTNAME"
if [[ $HOSTNAME == "musecore-new-validator" ]]; then
  source ~/add-keys.sh n
else
  source ~/add-keys.sh y
fi


# Pause other nodes so that the primary can node can do the genesis creation
if [ $HOSTNAME != "musecore0" ]
then
  while [ ! -f ~/.musecored/config/genesis.json ]; do
    echo "Waiting for genesis.json file to exist..."
    sleep 1
  done
  # need to wait for musecore0 to be up
  while ! curl -s -o /dev/null musecore0:26657/status ; do
    echo "Waiting for musecore0 rpc"
    sleep 1
done
fi

# Genesis creation following steps
# 1. Accumulate all the os_info files from other nodes on zetcacore0 and create a genesis.json
# 2. Add the observers , authorizations and required params to the genesis.json
# 3. Copy the genesis.json to all the nodes .And use it to create a gentx for every node
# 4. Collect all the gentx files in musecore0 and create the final genesis.json
# 5. Copy the final genesis.json to all the nodes and start the nodes
# 6. Update Config in musecore0 so that it has the correct persistent peer list
# 7. Start the nodes
# Start of genesis creation . This is done only on musecore0.
# Skip genesis if it has already been completed (marked by presence of ~/.musecored/init_complete file)
if [[ $HOSTNAME == "musecore0" && ! -f ~/.musecored/init_complete ]]
then
  MUSECORED_REPLICAS=2
  if host musecore3 ; then
    echo "musecore3 exists, setting MUSECORED_REPLICAS to 4"
    MUSECORED_REPLICAS=4
  fi
  # generate node list
  START=1
  # shellcheck disable=SC2100
  END=$((MUSECORED_REPLICAS - 1))
  NODELIST=()
  for i in $(eval echo "{$START..$END}")
  do
    NODELIST+=("musecore$i")
  done

  # Misc : Copying the keyring to the client nodes so that they can sign the transactions
  ssh museclient0 mkdir -p ~/.musecored/keyring-test/
  scp ~/.musecored/keyring-test/* museclient0:~/.musecored/keyring-test/
  ssh museclient0 mkdir -p ~/.musecored/keyring-file/
  scp ~/.musecored/keyring-file/* museclient0:~/.musecored/keyring-file/

  # 1. Accumulate all the os_info files from other nodes on zetcacore0 and create a genesis.json
  for NODE in "${NODELIST[@]}"; do
    INDEX=${NODE:0-1}
    ssh museclient"$INDEX" mkdir -p ~/.musecored/
    while ! scp "$NODE":~/.musecored/os_info/os.json ~/.musecored/os_info/os_z"$INDEX".json; do
      echo "Waiting for os_info.json from node $NODE"
      sleep 1
    done
    scp ~/.musecored/os_info/os_z"$INDEX".json museclient"$INDEX":~/.musecored/os.json
  done

  if host musecore-new-validator ; then
    echo "musecore-new-validator exists"
    ssh museclient-new-validator mkdir -p ~/.musecored/
    while ! scp musecore-new-validator:~/.musecored/os_info/os.json ~/.musecored/os_info/os_non_validator.json; do
          echo "Waiting for os_info.json from node musecore-new-validator"
          sleep 1
        done
    scp ~/.musecored/os_info/os_non_validator.json museclient-new-validator:~/.musecored/os.json
  fi

  ssh museclient0 mkdir -p ~/.musecored/
  scp ~/.musecored/os_info/os.json museclient0:/root/.musecored/os.json

  # 2. Add the observers, authorizations, required params and accounts to the genesis.json
  musecored collect-observer-info
  musecored add-observer-list --keygen-block 25

  # Add emissions withdraw authorizations
  if ! add_emissions_withdraw_authorizations; then
      echo "Error: Failed to add emissions withdraw authorizations"
      exit 1
  fi

  # Update governance and other chain parameters for localnet
  jq '
    .app_state.gov.params.voting_period="30s" |
    .app_state.gov.params.quorum="0.1" |
    .app_state.gov.params.threshold="0.1" |
    .app_state.gov.params.expedited_voting_period = "10s" |
    .app_state.gov.deposit_params.min_deposit[0].denom = "amuse" |
    .app_state.gov.params.min_deposit[0].denom = "amuse" |
    .app_state.staking.params.bond_denom = "amuse" |
    .app_state.crisis.constant_fee.denom = "amuse" |
    .app_state.mint.params.mint_denom = "amuse" |
    .app_state.evm.params.evm_denom = "amuse" |
    .app_state.emissions.params.ballot_maturity_blocks = "30" |
    .app_state.staking.params.unbonding_time = "10s" |
    .app_state.feemarket.params.min_gas_price = "10000000000.0000" |
    .consensus.params.block.max_gas = "500000000"
  ' "$HOME/.musecored/config/genesis.json" > "$HOME/.musecored/config/tmp_genesis.json" \
    && mv "$HOME/.musecored/config/tmp_genesis.json" "$HOME/.musecored/config/genesis.json"

  # set admin account
  admin_amount=100000000000000000000000000amuse # DEFAULT_FUND_AMOUNT * 10
  fund_account localnet_gov_admin muse1n0rn6sne54hv7w2uu93fl48ncyqz97d3kty6sh $admin_amount

  emergency_policy=$(yq -r '.policy_accounts.emergency_policy_account.bech32_address' /root/config.yml)
  admin_policy=$(yq -r '.policy_accounts.admin_policy_account.bech32_address' /root/config.yml)
  operational_policy=$(yq -r '.policy_accounts.operational_policy_account.bech32_address' /root/config.yml)

  fund_account emergency_policy "$emergency_policy" $admin_amount
  fund_account admin_policy "$admin_policy" $admin_amount
  fund_account operational_policy "$operational_policy" $admin_amount

  jq --arg emergency "$emergency_policy" \
    --arg operational "$operational_policy" \
    --arg admin "$admin_policy" '
      .app_state.authority.policies.items[0].address = $emergency |
      .app_state.authority.policies.items[1].address = $operational |
      .app_state.authority.policies.items[2].address = $admin
  ' "$HOME/.musecored/config/genesis.json" > "$HOME/.musecored/config/tmp_genesis.json" \
    && mv "$HOME/.musecored/config/tmp_genesis.json" "$HOME/.musecored/config/genesis.json"

  # Automatically fund most of the accounts
  fund_accounts_auto

  # 3. Copy the genesis.json to all the nodes .And use it to create a gentx for every node
  musecored gentx operator 1000000000000000000000amuse --chain-id=$CHAINID --keyring-backend=$KEYRING --gas-prices 20000000000amuse
  # Copy host gentx to other nodes
  for NODE in "${NODELIST[@]}"; do
    ssh $NODE mkdir -p ~/.musecored/config/gentx/peer/
    scp ~/.musecored/config/gentx/* $NODE:~/.musecored/config/gentx/peer/
  done
  # Create gentx files on other nodes and copy them to host node
  mkdir ~/.musecored/config/gentx/z2gentx
  for NODE in "${NODELIST[@]}"; do
      ssh $NODE rm -rf ~/.musecored/genesis.json
      scp ~/.musecored/config/genesis.json $NODE:~/.musecored/config/genesis.json
      ssh $NODE musecored gentx operator 1000000000000000000000amuse --chain-id=$CHAINID --keyring-backend=$KEYRING
      scp $NODE:~/.musecored/config/gentx/* ~/.musecored/config/gentx/
      scp $NODE:~/.musecored/config/gentx/* ~/.musecored/config/gentx/z2gentx/
  done

#  TODO : USE --modify flag to modify the genesis file when v18 is released
  if [[ -n "$MUSECORED_IMPORT_GENESIS_DATA" ]]; then
    echo "Importing data"
    musecored parse-genesis-file /root/genesis_data/exported-genesis.json
  fi

# 4. Collect all the gentx files in musecore0 and create the final genesis.json
  musecored collect-gentxs
  musecored validate-genesis

# 5. Copy the final genesis.json to all the nodes
  for NODE in "${NODELIST[@]}"; do
      ssh $NODE rm -rf ~/.musecored/genesis.json
      scp ~/.musecored/config/genesis.json $NODE:~/.musecored/config/genesis.json
  done

   if host musecore-new-validator > /dev/null; then
    echo "musecore-new-validator exists copying gentx peer"
     ssh musecore-new-validator rm -rf ~/.musecored/genesis.json
     scp ~/.musecored/config/genesis.json musecore-new-validator:~/.musecored/config/genesis.json
     ssh musecore-new-validator mkdir -p ~/.musecored/config/gentx/peer/
      # Check if gentx files exist before copying
     if ls ~/.musecored/config/gentx/* >/dev/null 2>&1; then
       if scp ~/.musecored/config/gentx/* musecore-new-validator:~/.musecored/config/gentx/peer/; then
         echo "Successfully copied gentx files to new-validator"
       else
         echo "Failed to copy gentx files to new-validator - Error code: $?"
       fi
     else
       echo "No gentx files found to copy"
     fi
   fi

# 6. Update Config in musecore0 so that it has the correct persistent peer list
   pp=$(cat $HOME/.musecored/config/gentx/z2gentx/*.json | jq '.body.memo' )
   pps=${pp:1:58}
   sed -i -e 's/^persistent_peers =.*/persistent_peers = "'$pps'"/' "$HOME"/.musecored/config/config.toml
fi
# End of genesis creation steps . The steps below are common to all the nodes

# Update persistent peers
if [[ $HOSTNAME != "musecore0" && ! -f ~/.musecored/init_complete ]]
then
  # Misc : Copying the keyring to the client nodes so that they can sign the transactions
  ssh museclient"$INDEX" mkdir -p ~/.musecored/keyring-test/
  scp ~/.musecored/keyring-test/* "museclient$INDEX":~/.musecored/keyring-test/
  ssh museclient"$INDEX" mkdir -p ~/.musecored/keyring-file/
  scp ~/.musecored/keyring-file/* "museclient$INDEX":~/.musecored/keyring-file/

  pp=$(cat $HOME/.musecored/config/gentx/peer/*.json | jq '.body.memo' )
  pps=${pp:1:58}
  sed -i -e "/persistent_peers =/s/=.*/= \"$pps\"/" "$HOME"/.musecored/config/config.toml
fi

# mark init completed so we skip it if container is restarted
touch ~/.musecored/init_complete

cosmovisor run start --pruning=nothing --minimum-gas-prices=0.0001amuse --json-rpc.api eth,txpool,personal,net,debug,web3,miner --api.enable --home /root/.musecored