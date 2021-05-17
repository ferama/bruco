import time

def handle_event(context, data):
    time.sleep(1.5)
    context.logger.info(data.decode())
    return data.decode() +  " test"

def init_context(context):
    context.logger.info("init context")
    setattr(context, "test", "test value")