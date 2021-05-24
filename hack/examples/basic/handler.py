import time

def handle_event(context, data):
    context.logger.info(data.decode())
    # time.sleep(0.02)
    # return data.decode() +  " test"

def init_context(context):
    context.logger.info("init context")
    setattr(context, "test", "test value")