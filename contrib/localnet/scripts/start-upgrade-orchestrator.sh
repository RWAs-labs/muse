#!/bin/bash

CHAINID="athens_101-1"
UPGRADE_AUTHORITY_ACCOUNT="muse10d07y265gmmuvt4z0w9aw880jnsr700jvxasvr"

if [[ -z $MUSECORED_URL ]]; then
    MUSECORED_URL='http://upgrade-host:8000/musecored'
fi
if [[ -z $MUSECLIENTD_URL ]]; then
    MUSECLIENTD_URL='http://upgrade-host:8000/museclientd'
fi

# Wait for authorized_keys file to exist (populated by musecore0)
while [ ! -f ~/.ssh/authorized_keys ]; do
    echo "Waiting for authorized_keys file to exist..."
    sleep 1
done

while ! curl -s -o /dev/null musecore0:26657/status ; do
    echo "Waiting for musecore0 rpc"
    sleep 1
done

# wait for minimum height
CURRENT_HEIGHT=0
while [[ $CURRENT_HEIGHT -lt 1 ]]
do
    CURRENT_HEIGHT=$(curl -s musecore0:26657/status | jq -r '.result.sync_info.latest_block_height')
    echo "current height is ${CURRENT_HEIGHT}, waiting for 1"
    sleep 1
done

# copy musecore0 config and keys if not running on musecore0
if [[ $(hostname) != "musecore0" ]]; then
  scp -r 'musecore0:~/.musecored/config' 'musecore0:~/.musecored/os_info' 'musecore0:~/.musecored/config' 'musecore0:~/.musecored/keyring-file' 'musecore0:~/.musecored/keyring-test' ~/.musecored/
  sed -i 's|tcp://localhost:26657|tcp://musecore0:26657|g' ~/.musecored/config/client.toml
fi

# get new musecored version
curl -L -o /tmp/musecored.new "${MUSECORED_URL}"
chmod +x /tmp/musecored.new
UPGRADE_NAME=$(/tmp/musecored.new upgrade-handler-version)

# if explicit upgrade height not provided, use dumb estimator
if [[ -z $UPGRADE_HEIGHT ]]; then
    UPGRADE_HEIGHT=$(( $(curl -s musecore0:26657/status | jq '.result.sync_info.latest_block_height' | tr -d '"') + 60))
    echo "Upgrade height was not provided. Estimating ${UPGRADE_HEIGHT}."
fi

cat > upgrade.json <<EOF
{
  "messages": [
    {
      "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
      "plan": {
        "height": "${UPGRADE_HEIGHT}",
        "info": "",
        "name": "${UPGRADE_NAME}",
        "time": "0001-01-01T00:00:00Z",
        "upgraded_client_state": null
      },
      "authority": "${UPGRADE_AUTHORITY_ACCOUNT}"
    }
  ],
  "metadata": "",
  "deposit": "1000000000000000000000amuse",
  "title": "${UPGRADE_NAME}",
  "summary": "${UPGRADE_NAME}"
}
EOF

# convert uname arch to goarch style
UNAME_ARCH=$(uname -m)
case "$UNAME_ARCH" in
    x86_64)    GOARCH=amd64;;
    i686)      GOARCH=386;;
    armv7l)    GOARCH=arm;;
    aarch64)   GOARCH=arm64;;
    *)         GOARCH=unknown;;
esac

cat > upgrade_plan_info.json <<EOF
{
    "binaries": {
        "linux/${GOARCH}": "${MUSECORED_URL}",
        "museclientd-linux/${GOARCH}": "${MUSECLIENTD_URL}"
    }
}
EOF

cat upgrade.json | jq --arg info "$(cat upgrade_plan_info.json)" '.messages[0].plan.info = $info' | tee upgrade_full.json

echo "Submitting upgrade proposal"

musecored tx gov submit-proposal upgrade_full.json --from operator --keyring-backend test --chain-id $CHAINID --yes --gas 300000 --fees 3000000000000000amuse -o json | tee proposal.json
PROPOSAL_TX_HASH=$(jq -r .txhash proposal.json)
PROPOSAL_ID=""
while [[ -z $PROPOSAL_ID ]]; do
    echo "waiting to get proposal_id"
    sleep 1
    PROPOSAL_ID=$(musecored query tx $PROPOSAL_TX_HASH -o json | jq -r '.events[] | select(.type == "submit_proposal") | .attributes[] | select(.key == "proposal_id") | .value')
done
echo "proposal id is ${PROPOSAL_ID}"

musecored tx gov vote "${PROPOSAL_ID}" yes --from operator --keyring-backend test --chain-id $CHAINID --yes --fees=2000000000000000amuse