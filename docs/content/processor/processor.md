---
title: "Processor"
draft: false
---

Bruco processor is a binary executable that handle all the

{{<mermaid align="center">}}
graph LR;
    A[Source] --> B(Processor) --> C[Sink]
{{< /mermaid >}}

flow. In details it does the following things:

1. Loads and parses a **config.yaml** file containing the stream definition
2. Starts the Source and Sink accordingly
3. Starts a pool of workers (python processes) as per configuration
4. Handles the stream

The processor section of the config file supports the following attributes:

| yaml path | description |
| ----------- | ----------- |
| handlerPath | the path to the python file where the handle is defined (by default the current working dir) |
| moduleName | the python module containing the **handle_event** function (by default handler) |
| workers | the number of worker to start up |
| env | an array of objects of type name, value defining the env vars to expose to the workers processes |

This is a full config example where we define processor, source and sink conf
```yaml
processor:
  handlerPath: ./examples/basic
  moduleName: handler
  workers: 2

source:
  kind: kafka
  brokers:
    - localhost:9092
  topics:
    - test
  offset: latest
  consumerGroup: my-consumer-group

sink:
  kind: kafka
  brokers:
    - localhost:9092
  topic: test-out
  partitioner: hash
```

## The handler module
Each worker instance, will handle events using a python handler module. A python handler module is a file named (by default) **handler.py**. The handler.py file **must** contain at least a function definition called **handle_event** that receives two parameters:

1. context
2. data

So the handler_event signature is:

```python
def handle_event(context, data)
```

The **context** is a persistent object that contains some helpers like a logger. You can also use this object to set some custom attributes at runtime.

The **data** param is a byte array containing the event data as sent from the source to the workers. The return value of handle_event (if any), will be routed back to the configured sink.

The **handler.py** could also define a special function called **init_context**

```python
def init_context(context)
```

The init_context function will be executed only once on worker startup. It can be used to setup some context as for example loading a machine learning model or setup a database connection.

### Handling dependencies

You can have a standard python **requirements.txt** file in the same directory where the **handler.py** resides, that list the required dependecies. If that file exists, bruco will run `pip install -r` against it. This behaviur can be disabled eposing an env var:

```sh
expose BRUCO_DISABLE_PIP=true
```