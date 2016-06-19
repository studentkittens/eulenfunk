#!/bin/sh

#CPU_USAGE=$(grep 'cpu ' /proc/stat | awk '{usage=($2+$4)*100/($2+$4+$5)} END {printf("%.0f\n", usage)}')
CPU_USAGE=$(mpstat 1 1 | awk '$3 ~ /CPU/ { for(i=1;i<=NF;i++) { if ($i ~ /%idle/) field=i } } $3 ~ /all/ { printf("%d",100 - $field) }')
MEM_USAGE=$(free | grep Mem | awk '{printf("%.0f\n", $3/$2 * 100.0)}')
IP_ADDR=$(ip addr show wlan0 | grep -m 1 'inet' | awk  '{print $2}' | sed 's/...$//g')
UPTIME_HOURS=$(uptime | awk '{print $1}')
WIFI_SIGNAL=$(awk 'NR==3 {print $4}''' /proc/net/wireless)
CPU_TEMP=$(cat /sys/class/thermal/thermal_zone0/temp | awk '{ print $0 / 1000 }')
CPU_FREQ=$(cat /sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_cur_freq | awk '{ print $0 / 1000 }')

printf "CPU: %3d%%  MEM: %3d%%\n" ${CPU_USAGE} ${MEM_USAGE}
printf "ADDR: %14s\n" ${IP_ADDR}
printf "❤ %sH  %4sMHz\n" ${UPTIME_HOURS} ${CPU_FREQ}
printf "ψ %3.3f%%   %3.3f৹C\n"  ${WIFI_SIGNAL} ${CPU_TEMP}
