#!/usr/bin/env bash
date
echo "traceroute"
traceroute 124.223.165.24
date

sleep 10
date
echo "ping"
ping -c 10 124.223.165.24
date
sleep 10
rm caddy
date
echo "download"
time wget http://124.223.165.24/caddy
date

