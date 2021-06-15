# This is a custom image version. You can deploy with a spec like this:
#
# apiVersion: bruco.ferama.github.io/v1alpha1
# kind: Bruco
# metadata:
#   name: example-bruco
# spec:
#   replicas: 1
#   image: your-registry/bruco-custom:latest
#   functionURL: /app
#   stream:
#     processor:
#       workers: 2
#     source:
#       kind: http

FROM ferama/bruco:dev

WORKDIR /app
COPY . .