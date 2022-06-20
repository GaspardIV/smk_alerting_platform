gcloud functions deploy site-checker --project=smk-alerting-platform --runtime=go113 --region=europe-central2 \
  --entry-point SiteChecker --env-vars-file=env.yaml --trigger-http --allow-unauthenticated --min-instances=5