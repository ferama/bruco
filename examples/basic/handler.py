import os

def handle_event(context, data):
    context.logger.info(data.decode())
    if context.var == False:
        context.logger.info("var is false")
        setattr(context, "var", True)
    else:
        context.logger.info("var is true")

    context.logger.info("##### ENV DUMP #####")
    for k, v in os.environ.items():
        context.logger.info(f"{k}: {v}")

    return data.decode() +  " test"

def init_context(context):
    context.logger.info("init context")
    setattr(context, "var", False)