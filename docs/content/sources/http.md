---
title: "Http"
weight: 1
draft: false
---


The http source is activated using the `http` kind

| yaml path | description |
| ----------- | ----------- |
| ignoreProcessorResponse | If true doen't return the prcessed value to the http caller |
| port | the port the http server listens too (do not change this if you are using bruco with k8s) |

Example:
```yaml
source:
  kind: http
  ignoreProcessorResponse: false
  # port: 8090
```