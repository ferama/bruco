import time

def handle_event(context, data):
    context.logger.info(data.decode())
    time.sleep(8)
    return data.decode() +  " test"

def init_context(context):
    context.logger.info("init context")
    setattr(context, "test", "test value")