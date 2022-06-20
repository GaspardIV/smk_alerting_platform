gcloud functions deploy scheduler --project=smk-alerting-platform --runtime=go113 --region=europe-central2 \
  --entry-point Scheduler --env-vars-file=env.yaml --trigger-topic scheduler