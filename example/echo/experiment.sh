#!/bin/bash

if [ "$#" -ne 4 ]; then
    echo "Illegal number of parameters"
    echo "./run <statrt> <increament> <stop> <rate>"
    exit 1
fi


for i in `seq $1 $2 $3`;
do
	./run.sh $i 10.254.254.239 ~/quic_results/quic 30 $4
done

source venv/bin/activate

python3 draw_plot.py 1 300000 $1 $3 $2 ~/quic_results/quic/$4

deactivate
