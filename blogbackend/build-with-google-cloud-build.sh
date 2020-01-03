IMAGE_URL=gcr.io/stunning-symbol-139515/web-oyvindsk.com    # The docker (or..?) image to deploy, uses tagged :latest ? fixme

gcloud builds submit --tag $IMAGE_URL