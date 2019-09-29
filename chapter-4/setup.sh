#!/bin/bash
./configcheck.sh

read -p "Do you want to continue and configure the project (y/N)? " -n 1 -r
if [[ $REPLY =~ ^[Yy]$ ]]
then
  echo ""
else
  exit 0
fi

# Cloud Run is still in beta
echo "Checking of gcloud beta is installed..."
( set -x ; gcloud components install beta 1> /dev/null )
echo ""

# Enable APIs
echo "Enabling services..."
( set -x ; gcloud services enable cloudbuild.googleapis.com)
( set -x ; gcloud beta services enable run.googleapis.com)
( set -x ; gcloud services enable firestore.googleapis.com)
( set -x ; gcloud services enable cloudtasks.googleapis.com)
echo ""

echo "Enabling services might take some time and the next \
commands might fail if they are not enabled..."
read -p "Do you want to continue now (y/N)? " -n 1 -r
if [[ $REPLY =~ ^[Yy]$ ]]
then
  echo ""
else
  exit 0
fi

echo "Creating task queue..."
( set -x ; gcloud tasks queues create bookings )
echo ""

echo "Unless you already enabled Firestore, please go \
to this link and choose $(tput bold)NATIVE MODE$(tput sgr0):"
echo "https://console.cloud.google.com/firestore?project=$(gcloud config get-value project)"

