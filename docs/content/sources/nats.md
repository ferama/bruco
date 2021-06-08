---
menuTitle: "Nats"
title: "NATS"
weight: 1
draft: false
---

The NATS source is activated using the `nats` kind

Example:
```yaml
source:
  kind: nats
  serverUrl: localhost:4222
  queueGroup: test
  subject: in.sub
```