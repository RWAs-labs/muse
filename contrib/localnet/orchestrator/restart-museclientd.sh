#!/bin/bash

# This script immediately restarts the museclientd on museclient0 and museclient1 containers in the network
# museclientd-supervisor will restart museclient automatically

echo restarting museclients

ssh -o "StrictHostKeyChecking no" "museclient0" -i ~/.ssh/localtest.pem killall museclientd
ssh -o "StrictHostKeyChecking no" "museclient1" -i ~/.ssh/localtest.pem killall museclientd