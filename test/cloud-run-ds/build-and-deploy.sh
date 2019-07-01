gcloud builds submit --tag gcr.io/stunning-symbol-139515/hello-world-1


gcloud beta run deploy --platform managed --region us-central1 --image gcr.io/stunning-symbol-139515/hello-world-1