#!/usr/bin/env bash

go build
cd client
go build
cd ..

RED='\033[0;31m'
NC='\033[0m'
DEBUG="false"

outputFiles=()


UIPort=12345
gossipPort=5000
name='A'

# General peerster (gossiper) command
#./Peerster -UIPort=12345 -gossipPort=127.0.0.1:5001 -name=A -peers=127.0.0.1:5002 > A.out &

for i in `seq 1 10`;
do
	outFileName="$name.out"
	peerPort=$((($gossipPort+1)%10+5000))
	peer="127.0.0.1:$peerPort"
	gossipAddr="127.0.0.1:$gossipPort"
	./SecuriChat -UIPort=$UIPort -addr=$gossipAddr -name=$name -peers=$peer > $outFileName &
	outputFiles+=("$outFileName")
	if [[ "$DEBUG" == "true" ]] ; then
		echo "$name running at UIPort $UIPort and gossipPort $gossipPort"
	fi
	UIPort=$(($UIPort+1))
	gossipPort=$(($gossipPort+1))
	name=$(echo "$name" | tr "A-Y" "B-Z")
done

./client/client -UIPort 12345 -encrypted=true -msg=eeLo -dest=F
sleep 2
./client/client -UIPort 12348 -encrypted=true -msg=hi -dest=C
sleep 60
pkill -f Peerster

