#!/bin/bash

# push docker image
gcloud builds submit --tag europe-central2-docker.pkg.dev/smk-alerting-platform/fake-service-repo/fake_service .
# restart pod
kubectl get pods  -n default --no-headers=true | awk '/fake-service-cluster/{print $1}' | xargs  kubectl delete -n default pod