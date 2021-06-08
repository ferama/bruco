---
menuTitle: "Nats"
title: "NATS"
weight: 1
draft: false
---

The nats sink is activated using the `nats` kind

Example:
```yaml
# OPTIONAL: you can have a source and a processor without a sink
sink:
  kind: nats
  serverUrl: localhost:4222
  subject: out.sub
```
