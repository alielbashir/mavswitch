#!/bin/bash

# Check if the script is run as root
if [ "$(id -u)" -ne 0 ]; then
  echo "This script must be run as root or with sudo privileges."
  exit 1
fi

interface=lo

# Get the user's input
read -p "Enter the UDP port number you induced packet loss on: " port

# Check if the input is valid
if ! [[ "$port" =~ ^[0-9]+$ ]]; then
  echo "Invalid input. Please enter a valid port number."
  exit 1
fi

# Delete the netem qdisc
tc qdisc del dev $interface parent 1:1 handle 10: netem

# Delete the filter matching the specific UDP port
tc filter del dev $interface protocol ip parent 1: prio 1 u32 match ip protocol 17 0xff match ip dport $port 0xffff flowid 1:1

# Delete the root qdisc
tc qdisc del dev $interface root handle 1: prio

echo "Traffic control rules reset for UDP port $port."
