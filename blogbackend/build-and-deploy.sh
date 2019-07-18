
PROJECT=stunning-symbol-139515                              # GCP project 
SERVICE_NAME=web-oyvindsk--dot--com                         # Service Name in Cloud Run
IMAGE_URL=gcr.io/stunning-symbol-139515/web-oyvindsk.com    # The docker (or..?) image to deploy

gcloud builds submit --tag $IMAGE_URL


gcloud beta run deploy $SERVICE_NAME --project $PROJECT --platform managed --region us-central1 --allow-unauthenticated --image $IMAGE_URL 