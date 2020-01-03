
PROJECT=stunning-symbol-139515                              # GCP project 
REGION=europe-west1                                         # GCP Region
SERVICE_NAME=web-oyvindsk--dot--com                         # Service Name in Cloud Run
IMAGE_URL=gcr.io/stunning-symbol-139515/web-oyvindsk.com    # The docker (or..?) image to deploy, uses tagged :latest ? fixme

gcloud beta run deploy $SERVICE_NAME --project $PROJECT --platform managed --region $REGION --allow-unauthenticated --image $IMAGE_URL 