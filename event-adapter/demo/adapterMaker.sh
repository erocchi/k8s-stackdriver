#!/usr/bin/env bash
cd demo

cat adapter-demo-1.yaml >adapter-demo-2.yaml

sed -i 's/- --max-retrieved-events=100/- --max-retrieved-events='"$1"'/g' adapter-demo-2.yaml
sed -i 's/- --retrieve-events-since-millis=0/- --retrieve-events-since-millis='"$2"'/g' adapter-demo-2.yaml

kubectl create -f adapter-demo-2.yaml --validate=false
echo Waiting...
sleep 15
google-chrome http://localhost:8001/apis/v1events/v1alpha1/namespaces/default/events
google-chrome http://localhost:8001/apis/v1events/v1alpha1/events

read -p 'Press a key to continue' temp

kubectl delete -f adapter-demo-2.yaml

cd ..

