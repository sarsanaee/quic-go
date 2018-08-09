#!/bin/bash


if [ "$#" -ne 5 ]; then
    echo "Illegal number of parameters"
    echo "./run <number_of_clients> <server_ip> <results_path> <experiments_time> <rate>"
    exit 1
fi

#result_path=$HOME/quic_results/quic
result_path=$3
rate=$5

mkdir $result_path/$rate

pkill my_quic_client

ssh scc@$2 "pkill my_quic_server"

sleep 2

ssh -f scc@$2 "export GOPATH=$HOME/work; export PATH=$PATH:/usr/local/go/bin; go run /home/scc/work/src/github.com/lucas-clemente/quic-go/example/echo/my_quic_server.go &"

sleep 2

rm $result_path/$rate/*.log #removing current logs

for i in `seq 1 $1`;
do
	echo $i
	go run my_quic_client.go $rate > $result_path/$rate/$i\_$1.log &
done

sleep $4



pkill my_quic_client
ssh scc@$2 "pkill my_quic_server"

echo "" > $result_path/$rate/latency_quic_$rate\_$1.txt
cat $result_path/$rate/*.log >> $result_path/$rate/latency_quic_$rate\_$1.txt
