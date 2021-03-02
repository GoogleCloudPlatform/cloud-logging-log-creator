# Chattylogs

This project is a simple go application that spits out many entries into the 
logs of a Google Cloud Project.  The purpose is to have a default set of log 
entries with which to expore the features of Cloud Logging. It's designed to 
be a trivial little app that emits a lot of log chatter. It also has an http
service that listens on a default port so that the application can be used 
with http based services like Cloud Run, App Engine, or installed as a container 
in manually built Kubernetes, GKE, Anthos, or a GCE instance. 

Potential future updates could include a standalone service for GCE and GCF. 

## Using the image
The goal is to get the image hosted on Google Cloud Repositories. To do that:

* Set the default project  
`gcloud config set project [PROJECT_ID]`
* Make the image and push to GCR



## Requirements

* Google Cloud SDK
* Docker
* make


<hr>
This is not an official Google product. 