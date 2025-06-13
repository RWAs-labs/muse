#!/bin/bash

# This script allows to add keys for operator and hotkey and create the required json structure for os_info

KEYRING_TEST="test"
KEYRING_FILE="file"
HOSTNAME=$(hostname)

# Check if is_observer flag is provided
if [ -z "$1" ]; then
    is_observer="y" # Default value if not provided
else
    is_observer="$1"
fi


musecored keys add operator --algo=secp256k1 --keyring-backend=$KEYRING_TEST

operator_address=$(musecored keys show operator -a --keyring-backend=$KEYRING_TEST)

# Hotkey key depending on the keyring-backend
if [ "$HOTKEY_BACKEND" == "$KEYRING_FILE" ]; then
    printf "%s\n%s\n" "$HOTKEY_PASSWORD" "$HOTKEY_PASSWORD" | musecored keys add hotkey --algo=secp256k1 --keyring-backend=$KEYRING_FILE
    hotkey_address=$(printf "%s\n%s\n" "$HOTKEY_PASSWORD" "$HOTKEY_PASSWORD" | musecored keys show hotkey -a --keyring-backend=$KEYRING_FILE)

    # TODO: remove after v50 upgrade
    # Get hotkey pubkey, the command use keyring-backend in the cosmos config
    if ! musecored config set client keyring-backend "$KEYRING_FILE"; then
        musecored config keyring-backend "$KEYRING_FILE"
    fi
    pubkey=$(printf "%s\n%s\n" "$HOTKEY_PASSWORD" "$HOTKEY_PASSWORD" | musecored get-pubkey hotkey | sed -e 's/secp256k1:"\(.*\)"/\1/' |sed 's/ //g' )
    if ! musecored config set client keyring-backend "$KEYRING_TEST"; then
        musecored config keyring-backend "$KEYRING_TEST"
    fi
else
    musecored keys add hotkey --algo=secp256k1 --keyring-backend=$KEYRING
    hotkey_address=$(musecored keys show hotkey -a --keyring-backend=$KEYRING)
    pubkey=$(musecored get-pubkey hotkey|sed -e 's/secp256k1:"\(.*\)"/\1/' | sed 's/ //g' )
fi

echo "operator_address: $operator_address"
echo "hotkey_address: $hotkey_address"
echo "pubkey: $pubkey"
echo "is_observer: $is_observer"
mkdir -p ~/.musecored/os_info

# set key in file
jq -n --arg is_observer "$is_observer" --arg operator_address "$operator_address" --arg hotkey_address "$hotkey_address" --arg pubkey "$pubkey" '{"IsObserver":$is_observer,"ObserverAddress":$operator_address,"MuseClientGranteeAddress":$hotkey_address,"MuseClientGranteePubKey":$pubkey}' > ~/.musecored/os_info/os.json