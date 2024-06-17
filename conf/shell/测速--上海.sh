#!/usr/bin/env bash
date
echo "traceroute"
traceroute 43.153.119.127
date

sleep 10
date
echo "ping"
ping -c 10 43.153.119.127
date
sleep 10
rm caddy
date
echo "download"
time wget http://43.153.119.127/caddy
date

