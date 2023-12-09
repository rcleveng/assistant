#!/bin/bash

GCLOUD=$(command -v gcloud) || { (>&2 echo 'unable to find gcloud') ; exit 1; }


if [[ -z "${COMMIT_SHA}" ]]; then
	# Get latest commit sha
	COMMIT_SHA=$(git rev-parse --verify HEAD)
	if [[ -z "${COMMIT_SHA}" ]]; then
		echo "Unable to determine commit sha"
		exit 1
	fi
fi


CMD="${GCLOUD} \
  builds \
  submit \
  --region=us-central1 \
  --config cloudbuild.yaml \
  --substitutions COMMIT_SHA=${COMMIT_SHA}
  $@"

echo "COMMIT_SHA: ${COMMIT_SHA}"
echo "GCLOUD:     ${GCLOUD}"
echo "CMD:        ${CMD}"

${CMD}