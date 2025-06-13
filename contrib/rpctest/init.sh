#!/usr/bin/env bash

CHAINID="localnet_101-1"
KEYRING="test"
export DAEMON_HOME=$HOME/.musecored
export DAEMON_NAME=musecored

### chain init script for development purposes only ###
rm -rf ~/.musecored
kill -9 $(lsof -ti:26657)
musecored config keyring-backend $KEYRING --home ~/.musecored
musecored config chain-id $CHAINID --home ~/.musecored
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | musecored keys add muse --algo=secp256k1 --recover --keyring-backend=test
echo "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard" | musecored keys add mario --algo secp256k1 --recover --keyring-backend=test
echo "lounge supply patch festival retire duck foster decline theme horror decline poverty behind clever harsh layer primary syrup depart fantasy session fossil dismiss east" | musecored keys add museeth --recover --keyring-backend=test

musecored init test --chain-id=$CHAINID

#Set config to use amuse
cat $HOME/.musecored/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="amuse"' > $HOME/.musecored/config/tmp_genesis.json && mv $HOME/.musecored/config/tmp_genesis.json $HOME/.musecored/config/genesis.json
cat $HOME/.musecored/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="amuse"' > $HOME/.musecored/config/tmp_genesis.json && mv $HOME/.musecored/config/tmp_genesis.json $HOME/.musecored/config/genesis.json
cat $HOME/.musecored/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="amuse"' > $HOME/.musecored/config/tmp_genesis.json && mv $HOME/.musecored/config/tmp_genesis.json $HOME/.musecored/config/genesis.json
cat $HOME/.musecored/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="amuse"' > $HOME/.musecored/config/tmp_genesis.json && mv $HOME/.musecored/config/tmp_genesis.json $HOME/.musecored/config/genesis.json
cat $HOME/.musecored/config/genesis.json | jq '.app_state["evm"]["params"]["evm_denom"]="amuse"' > $HOME/.musecored/config/tmp_genesis.json && mv $HOME/.musecored/config/tmp_genesis.json $HOME/.musecored/config/genesis.json
cat $HOME/.musecored/config/genesis.json | jq '.consensus["params"]["block"]["max_gas"]="10000000"' > $HOME/.musecored/config/tmp_genesis.json && mv $HOME/.musecored/config/tmp_genesis.json $HOME/.musecored/config/genesis.json






musecored add-genesis-account $(musecored keys show muse -a --keyring-backend=test) 500000000000000000000000000000000000000amuse --keyring-backend=test
musecored add-genesis-account $(musecored keys show mario -a --keyring-backend=test) 50000000000000000000000000000000000000amuse --keyring-backend=test
musecored add-genesis-account $(musecored keys show museeth -a --keyring-backend=test) 500000000000000000000000000000000amuse --keyring-backend=test


ADDR1=$(musecored keys show muse -a --keyring-backend=test)
observer+=$ADDR1
observer+=","
ADDR2=$(musecored keys show mario -a --keyring-backend=test)
observer+=$ADDR2
observer+=","


observer_list=$(echo $observer | rev | cut -c2- | rev)

echo $observer_list



musecored add-observer 1337 "$observer_list"
musecored add-observer 101 "$observer_list"




musecored gentx muse 50000000000000000000000000amuse --chain-id=localnet_101-1 --keyring-backend=test

contents="$(jq '.app_state.gov.voting_params.voting_period = "10s"' $DAEMON_HOME/config/genesis.json)" && \
echo "${contents}" > $DAEMON_HOME/config/genesis.json

echo "Collecting genesis txs..."
musecored collect-gentxs

echo "Validating genesis file..."
musecored validate-genesis
#
#export DUMMY_PRICE=yes
#export DISABLE_TSS_KEYGEN=yes
#export GOERLI_ENDPOINT=https://goerli.infura.io/v3/faf5188f178a4a86b3a63ce9f624eb1b
