
PROJECT=stunning-symbol-139515                              # GCP project 
REGION=europe-west1                                         # GCP Region
SERVICE_NAME=web-oyvindsk--dot--com                         # Service Name in Cloud Run
IMAGE_URL=gcr.io/stunning-symbol-139515/web-oyvindsk.com    # The docker (or..?) image to deploy, uses tagged :latest ? FIXME

gcloud run deploy $SERVICE_NAME --project $PROJECT --platform managed --region $REGION --allow-unauthenticated --image $IMAGE_URL --concurrency 1000

## TODO: stop here if that failed ^


# Set the latest til 100% traffic. 
# Uncomment if you have reverted to an older one and don't want this new one to recevive traffic right away 
# Should do nothing if the above deployed failed, since then the latest already has 100% traffic (probably :)
gcloud run services update-traffic $SERVICE_NAME --project $PROJECT --platform managed --region $REGION --to-latest
