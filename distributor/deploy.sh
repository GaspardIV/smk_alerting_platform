gcloud functions deploy distributor --project=smk-alerting-platform --runtime=go113 --region=europe-central2 \
  --entry-point Distributor --env-vars-file=env.yaml --trigger-http --allow-unauthenticated --min-instances=1