import random

def handle_event(context, data):
    context.logger.info(data.decode())
    key = random.randrange(1, 100)
    transormed = data.decode() + " keyed"
    return context.response(key, transormed)