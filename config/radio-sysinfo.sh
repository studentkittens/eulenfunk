#!/bin/sh

CPU_USAGE=$(grep 'cpu ' /proc/stat | awk '{usage=($2+$4)*100/($2+$4+$5)} END {printf("%.0f\n", usage)}')
MEM_USAGE=$(free | grep Mem | awk '{printf("%.0f\n", $3/$2 * 100.0)}')
IP_ADDR=$(ip addr show wlan0 | grep -m 1 'inet' | awk  '{print $2}' | sed 's/...$//g')
UPTIME_HOURS=$(uptime | awk '{print $1}')
WIFI_SIGNAL=$(awk 'NR==3 {print $4}''' /proc/net/wireless)

printf "π %3d%% ⏸ %3d%%\n" ${CPU_USAGE} ${MEM_USAGE}
printf "IP: %15s\n" ${IP_ADDR}
printf "❤ %s ৹C\n" ${UPTIME_HOURS}
printf "ψ %3.1f%%\n"  ${WIFI_SIGNAL}
