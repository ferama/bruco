from lib import ciao
import time

def handle_event(context, data):
    # context.logger(params)
    # print("python stdout test")
    # context.ciao()
    time.sleep(4)
    return data.decode() +  " test"

def init_context(context):
    # context.logger("prova")
    # context.logger("from init context")
    # print("init context")
    setattr(context, "ciao", ciao)