---
title: "Bruco"
draft: false
---

# Bruco
Bruco is a tool meant to build streaming pipelines steps easily. It is **kubernetes** native citizen. Each step can be indeed, defined using a Kubernetes custom resource. You don't even need to manually build a docker image. 

The pipeline is event-driven and implements the `source -> processor -> sink` paradigm.

The processor is meant for `stream` transformation between the **source** and the **sink**. Bruco supports writing processor logic using the `python` scripting language.

#### Sources

Actually bruco supports the following event sources:

1. HTTP
2. Apache Kafka
3. NATS

#### Sinks

Follows the actually supported sinks list

1. Apache Kafka
2. NATS

{{% notice note %}}
Follow the [Getting Started]({{< ref "basic/getting-started" >}}) guide for a beginner tutorial.
{{% /notice %}}