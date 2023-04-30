#!/bin/bash

# Check if the script is run as root
if [ "$(id -u)" -ne 0 ]; then
  echo "This script must be run as root or with sudo privileges."
  exit 1
fi

interface=lo

# Get the user's input
read -p "Enter the UDP port number to induce packet loss on: " port
read -p "Enter the desired packet loss percentage (0-100): " packet_loss

# Check if the input is valid
if ! [[ "$port" =~ ^[0-9]+$ ]] || ! [[ "$packet_loss" =~ ^[0-9]+$ ]] || [ "$packet_loss" -gt 100 ]; then
  echo "Invalid input. Please enter a valid port number and packet loss percentage."
  exit 1
fi

# Create a new qdisc on the output interface
tc qdisc add dev $interface root handle 1: prio

# Add a filter to match the specific UDP port
tc filter add dev $interface protocol ip parent 1: prio 1 u32 match ip protocol 17 0xff match ip dport $port 0xffff flowid 1:1

# Add the netem qdisc with the desired packet loss
tc qdisc add dev $interface parent 1:1 handle 10: netem loss $packet_loss%

echo "Packet loss induced on UDP port $port with $packet_loss% packet loss."
