---
title: "Bruco on kubernetes"
draft: false
---

Bruco on kubernetes is supported by a custom resource and controller.

This is a resource instance definition example:
```yaml
apiVersion: bruco.ferama.github.io/v1alpha1
kind: Bruco
metadata:
  name: example-bruco-s3
spec:
  replicas: 1
  env:
    - name: AWS_ACCESS_KEY_ID
      value: bruco
    - name: AWS_SECRET_ACCESS_KEY
      value: bruco123
  functionURL: s3://bruco-minio.bruco:9000/brucos/image-classifier.zip
  stream:
    processor:
      workers: 2
    source:
      kind: http
```

The resource supports the following fields:
| yaml path | description |
| ----------- | ----------- |
| replicas | the bruco pod replicas that should be instantiated |
| image | if you wish to set a custom docker image |
| imagePullPolicy | standard kubernetes imagePullPolicy |
| imagePullSecrets | standard kubernetes imagePullSecrets |
| resources | standard kubernetes resources limit and requests |
| env | an array of objects of kind name, value that defines env vars exposed to the pods |
| functionURL | a function package url. This where the handler code and requirements resides |
| stream | is the bruco config subsection |
| stream.processor | refer to the [processor docs]({{< ref "processor/processor" >}}) for details |
| stream.source | refer to the [source docs]({{< ref "sources" >}}) for details |
| stream.sink | refer to the [sink docs]({{< ref "sinks" >}}) for details |

### functionURL in depth
Bruco actually is able to load functions from three sources:

1. local path
2. http endpoint
3. s3 endpoint

##### local path
The local path supports loading, assumes that the function zip archive resides into the local filesystem. This is usefull if you want/need to build a custom image.

##### http endpoint
The http endpoint will let you load a function from an http endpoint.
Example:
```yaml
functionURL: https://github.com/ferama/bruco/raw/main/examples/zipped/sentiment.zip
```
{{% notice note %}}
Actually the endpoint needs to be public, bruco doesn't not support any kind of auth mechanism. It can be a service inside your k8s cluster
{{% /notice %}}

##### s3 endpoint
The s3 endpoint will let you load a function from minio/s3. It requires the definition of two env vars:
* AWS_ACCESS_KEY_ID
* AWS_SECRET_ACCESS_KEY

Example:
```yaml
env:
- name: AWS_ACCESS_KEY_ID
  value: bruco
- name: AWS_SECRET_ACCESS_KEY
  value: bruco123
functionURL: s3://bruco-minio.bruco:9000/brucos/image-classifier.zip
```
