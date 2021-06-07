---
title: "Getting started on k8s"
date: 2021-06-05T08:21:52+02:00
draft: false
---

This guide will introduce you to the **bruco** concept e will help to get the stuff up and running quickly.

This tutorial make some assumptions:

1. You have a kubernetes environment up and running
2. You have correctly configured the **kubectl** cli tool to interact with the env above

Bruco is heavily integrated with Kubernetes. To get start to experiment with bruco function deployment you need to setup some stuff on kubernetes first.

Bruco defines a custom resource for its functions. The resource is called guess that, bruco :)

This is a sample resource instance that you will able to create:
```yaml
apiVersion: brucocontroller.ferama.github.com/v1alpha1
kind: Bruco
metadata:
  name: example-bruco
spec:
  replicas: 1
  functionURL: https://github.com/ferama/bruco/raw/main/hack/examples/zipped/sentiment.zip
```

but before that, you need to create the custom resource definition on k8s and to deploy the controller supporting the custom resource.

Follow this steps to do that:

```sh
# create a new k8s namespace for bruco
$ kubectl create ns bruco

# create the custom resource
$ kubectl -n bruco apply -f https://raw.githubusercontent.com/ferama/bruco/main/hack/k8s/resources/crd-bruco.yaml

# deploy the bruco controller
$ kubectl -n bruco apply -f https://raw.githubusercontent.com/ferama/bruco/main/hack/k8s/resources/controller.yaml
```

TODO: 