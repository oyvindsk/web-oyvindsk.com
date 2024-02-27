#!/bin/bash

## deploy "directly" t oCloud Run, that is using a buildback
## they provide. It still builds with Cloud Build and puts the image
## in the Artifact Regestry
##
## --source ALSO supports a Dockerfile now, instead of buildspacks.
## Just put the Dockerfile in this directory

## Created a the bucket in the console. 1 region (europe-west1), standard other stuff, not public

gcloud beta run deploy storage-test --project stunning-symbol-139515 --region europe-west1 --source . \
--execution-environment gen2 \
--add-volume=name=foo-volume,type=cloud-storage,bucket=oyvindsk-test-cloud-run  \
--add-volume-mount=volume=foo-volume,mount-path=/gcs_mount_test
