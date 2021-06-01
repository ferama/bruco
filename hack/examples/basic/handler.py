def handle_event(context, data):
    context.logger.info(data.decode())
    if context.var == False:
        context.logger.info("var is false")
        setattr(context, "var", True)
    else:
        context.logger.info("var is true")
    return data.decode() +  " test"

def init_context(context):
    context.logger.info("init context")
    setattr(context, "var", False)