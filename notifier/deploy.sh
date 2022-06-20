gcloud functions deploy notifier --project=smk-alerting-platform --runtime=go113 --region=europe-central2 \
  --entry-point Notifier --env-vars-file=env.yaml --trigger-http --allow-unauthenticated