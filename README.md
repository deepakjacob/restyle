# restyle


#### Docker Build

```
docker image build -t restyle:server .
```

#### Docker Run
```
docker run -p 8000:8000
    -e GOOGLE_PROJECT_ID=$GOOGLE_PROJECT_ID
    -e GOOGLE_APPLICATION_CREDENTIALS=/root/.config/project-up-ed35607315b3.json
    -e GOOGLE_OAUTH_CLIENT_ID=$GOOGLE_OAUTH_CLIENT_ID
    -e GOOGLE_OAUTH_CLIENT_SECRET=$GOOGLE_OAUTH_CLIENT_SECRET
    -e JWT_SHARED_KEY=$JWT_SHARED_KEY
    -e JWT_SHARED_ENC_KEY=$JWT_SHARED_ENC_KEY
    -e REDIRECT_URL=$REDIRECT_URL
    -v ~/Development/goprojects/gcloud_service_accounts:/root/.config restyle:server
```

#### Running container

- Shell access
```
docker exec -it <CONTAINER_NAME>  sh
```
- Note that to assign container name add `--name=<container_name>` to docker run command

#### Tag and push image to gcr.io

- Tag an image restyle_server with tag 0.0.1 to gcr.io/project-up-238914/restyle:0.0.1

```
docker tag restyle_server:0.0.1 gcr.io/project-up-238914/restyle:0.0.1                     `
```
- To push to google container registry

```
docker push gcr.io/project-up-238914/restyle:0.0.1
```

#### Configuring and running gcloud

- get the account details
```
gcloud auth list
```
- list the current project
```
gcloud config list project
```
- set the compute zone
```
gcloud config set compute/zone us-central1-a
```
- get the gcloud context (this basically returns project / computezone etc)
```
kubectl config current-context
```
- create the container cluster
```
gcloud container clusters create restyle-cluster
```
- gcloud container clusters get-credentials - fetch credentials for a running cluster
```
gcloud container clusters get-credentials restyle-cluster
```
