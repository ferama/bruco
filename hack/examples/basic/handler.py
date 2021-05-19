import time

def handle_event(context, data):
    context.logger.info(data.decode())
    time.sleep(1.5)
    return data.decode() +  " test"

def init_context(context):
    import os
    for k, v in os.environ.items():
        print(f"{k}={v}")
    context.logger.info("init context")
    setattr(context, "test", "test value")