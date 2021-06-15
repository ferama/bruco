def handle_event(context, data):
    context.logger.info(data.decode())
    return data.decode()

def init_context(context):
    context.logger.info("init context from custom image")