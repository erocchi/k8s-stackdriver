#!/usr/bin/env bash
make push GCP_PROJECT=erocchi-gke-dev-1attempt
kubectl delete -f adapter.yaml 
kubectl create -f adapter.yaml --validate=false

read -p 'Press a key to continue' temp

kubectl delete -f adapter.yaml
kubectl create -f demo/adapter-demo-1.yaml --validate=false

read -p 'Continue? [Y/n]' temp

kubectl delete -f demo/adapter-demo-1.yaml

while [ "$temp" == "y" ]
do
    read -p '--max-retrieved-events=' N
    read -p '--retrieve-events-since-millis=' T
    demo/adapterMaker.sh $N $T
    read -p 'Continue? [Y/n]' temp
done

echo Thanks

