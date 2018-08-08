#!/bin/bash

result_path=$HOME/Documents/quic_results/quic

pkill my_quic_server
pkill my_quic_client

go run my_quic_server.go &

sleep 2

for i in `seq 1 $1`;
do
	echo $i
	go run my_quic_client.go > $result_path/$i.log &
done

sleep 20



pkill my_quic_client
pkill my_quic_server