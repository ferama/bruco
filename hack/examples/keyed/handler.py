import random

def handle_event(context, data):
    context.logger.info(data.decode())
    key = random.randrange(1, 100)
    ss
    return context.response(key, data.decode() +  " keyed")