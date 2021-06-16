---
title: "Getting started on k8s"
draft: false
---

This guide will introduce you to the **bruco** concept e will help to get the stuff up and running quickly.

This tutorial make some assumptions:

1. You have a kubernetes environment up and running
2. You have correctly configured the **kubectl** cli tool to interact with the env above

Bruco is heavily integrated with Kubernetes. To get start to experiment with bruco function deployment you need to setup some stuff on kubernetes first.

Bruco defines a custom resource for its functions. The resource is called guess that, bruco 🙂

This is a sample resource instance that you will able to create:
```yaml
apiVersion: bruco.ferama.github.io/v1alpha1
kind: Bruco
metadata:
  name: example-bruco
spec:
  replicas: 1
  functionURL: https://github.com/ferama/bruco/raw/main/examples/zipped/sentiment.zip
  stream:
    processor:
      workers: 2
    source:
      kind: http
      ignoreProcessorResponse: false
```

but before that, you need to create the custom resource definition on k8s and to deploy the controller supporting the custom resource.

Follow this steps to do that.

## k8s controller deployment
```sh
# create a new k8s namespace for bruco
$ kubectl create ns bruco

# create the custom resource
$ kubectl -n bruco apply -f https://raw.githubusercontent.com/ferama/bruco/main/hack/k8s/resources/crd-bruco.yaml

# deploy the bruco controller
$ kubectl -n bruco apply -f https://raw.githubusercontent.com/ferama/bruco/main/hack/k8s/resources/controller.yaml
```

The result should be something like this, that means that the bruco controller is up and running:

```sh
$ kubectl -n bruco get pods
NAME                               READY   STATUS    RESTARTS   AGE
bruco-controller-5fd955d49-7p6xt   1/1     Running   0          37s
```

## demo function deployment
Now you are ready to deploy your first bruco function. Create a file named **example-bruco.yaml** and copy and paste the example bruco function definition:
```yaml
apiVersion: bruco.ferama.github.io/v1alpha1
kind: Bruco
metadata:
  name: example-bruco
spec:
  replicas: 1
  functionURL: https://github.com/ferama/bruco/raw/main/examples/zipped/sentiment.zip
  stream:
    processor:
      workers: 2
    source:
      kind: http
      ignoreProcessorResponse: false
```

Create the kubernetes resource using:
```sh
$ kubectl -n bruco apply -f example-bruco.yaml
```

What is happening here, is that bruco will load the package at `https://github.com/ferama/bruco/raw/main/hack/examples/zipped/sentiment.zip` and starts a deployment with one replicas that will run the **sentiment** function. The zip file contains all the required config and logic to run the function. In this example the zip file contains:

1. An `handler.py` with the function logic
2. A `requirements.txt` file that declares the funcion dependencies

The handler is a very simple python function:
```python
from textblob import TextBlob

def handle_event(context, data):
    blob = TextBlob(data.decode())
    return {
        "sentiment": blob.sentiment.polarity,
        "subjectivity": blob.sentiment.subjectivity
    }
```
The handler function requires the `textblob` library. The dependency is defined in the `requirements.txt` file
```
textblob
```

Now you are ready to test out your first bruco function. Let's forward the `http source` port:

```sh
# forwards the http source port
$ kubectl -n bruco port-forward svc/example-bruco 8080
# generate an event using the http source
$ curl -X POST -d "bruco is great" http://localhost:8080
{"sentiment": 0.8, "subjectivity": 0.75}
```