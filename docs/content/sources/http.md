---
title: "Http"
weight: 1
draft: false
---


The http source is activated using the `http` kind

| yaml path | description |
| ----------- | ----------- |
| ignoreProcessorResponse | If true the http request returns quickly. Do not wait for a processor response |
| port | the port the http server listens too (do not change this if you are using bruco with k8s) |

Example:
```yaml
source:
  kind: http
  ignoreProcessorResponse: false
  # port: 8090
```