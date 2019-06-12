# kubernetes-paam

A Kubernetes Pod Anti-Affinity Monitor exposed over HTTP.

[![Docker Automated build](https://img.shields.io/docker/cloud/automated/kristofferahl/kubernetes-paam.svg?style=for-the-badge)](https://hub.docker.com/r/kristofferahl/kubernetes-paam/)
[![Docker Build Status](https://img.shields.io/docker/cloud/build/kristofferahl/kubernetes-paam.svg?style=for-the-badge)](https://hub.docker.com/r/kristofferahl/kubernetes-paam/)
[![MicroBadger Size](https://img.shields.io/microbadger/image-size/kristofferahl/kubernetes-paam.svg?style=for-the-badge)](https://hub.docker.com/r/kristofferahl/kubernetes-paam/)
[![Docker Pulls](https://img.shields.io/docker/pulls/kristofferahl/kubernetes-paam.svg?style=for-the-badge)](https://hub.docker.com/r/kristofferahl/kubernetes-paam/)


## Configuration

| Environment variable     | Description                                                     | Default |
|--------------------------|-----------------------------------------------------------------|---------|
| PAAM_HTTP_BIND_ADDRESS   | The address paam will listen on                                 | :8113   |
| PAAM_ONLY_FAILED_RESULTS | Only include failed results in the HTTP response                | false   |
| PAAM_EXCLUDE_NAMESPACES  | List of kubernetes namespaces to exclude (separated by commas)  | -       |
| PAAM_EXCLUDE_DEPLOYMENTS | List of kubernetes deployments to exclude (separated by commas) | -       |


## Development

As kubernetes-paam expects to run inside the kubernetes cluster it monitors, you need to build and push a new docker image, and finally deploy it to the cluster before you can try out the new features.

```bash
export IMAGE_TAG='dev'
docker build -t kristofferahl/kubernetes-paam:${IMAGE_TAG:?} .
docker push kristofferahl/kubernetes-paam:${IMAGE_TAG:?}
cat ./example.yaml | envsubst | kubectl apply -
```


## Releasing

Automatic builds are set up to produce docker images for master (tag=latest) and for git tags (tag=<tag-name>).
