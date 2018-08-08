#!/bin/bash

for i in `seq 1 10`;
do
	./run.bash $i 10.254.254.1 ~/quic_results/quic 100 5000
done
