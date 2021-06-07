---
title: "Bruco"
date: 2021-06-05T08:21:52+02:00
draft: false
---
Bruco is a tool meant to build streaming pipelines steps easily. It is **kubernetes** native citizen. Each step could be indeed, defined using a Kubernetes custom resource. You don't even need to manually build a docker image. 

The pipeline is event-driven and implements the `source -> processor -> sink` paradigm.

The processor is meant for `stream` transformation beteween the **source** and the **sink**. Bruco supports writing processor logic using the `python` scripting language.

#### Sources

Actually bruco supports the following event sources:

1. HTTP
2. Apache Kafka
3. NATS

#### Sinks

Follows the actually supported sinks list

1. Apache Kafka
2. NATS

### Tutorial

Follow the [Getting Started]({{< ref "getting-started" >}}) guide for a beginner tutorial.