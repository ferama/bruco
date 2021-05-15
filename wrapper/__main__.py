import os
import sys
import json
import importlib
import signal
import socket

class Context: pass


class Wrapper:
    def __init__(self, working_directory, port: int):
        self.port = port
        os.chdir(working_directory)
        signal.signal(signal.SIGINT, self.sigint_handler)

    def sigint_handler(self, p1, p2):
        sys.exit(0)

    def start(self):
        context = Context()
        module = importlib.import_module("handler")
        module.init_context(context)

        client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        # print(f"Connecting to {socket.gethostname()}:{self.port}")
        client.connect((socket.gethostname(), self.port))

        while True:
            msg = client.recv(1024 * 1024 * 100)
            try:
                response = module.handle_event(context, msg)
                if not response: response = ""
                out = {
                    "response": response,
                    "error": ""
                }
                out = json.dumps(out)
                out += "\n"
                client.send(out.encode())
            except Exception as e:
                out = {
                    "response": "",
                    "error": e
                }
                out = json.dumps(out)
                out += "\n"
                client.send(out.encode())
        



if __name__ == "__main__":
    import sys
    w = Wrapper(".", int(sys.argv[1]))
    w.start()