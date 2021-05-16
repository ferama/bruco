import os
import sys
import json
import importlib
import signal
import socket

class Context: pass


class Wrapper:
    def __init__(self, lambda_path, port: int):
        self.port = port
        os.chdir(lambda_path)
        sys.path.append(".")
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
    import argparse
    parser = argparse.ArgumentParser(description="the python shell")
    parser.add_argument("--lambda-path", 
                            required=True,
                            metavar="lambda_path", 
                            type=str, 
                            help="the working directory")
    parser.add_argument("--port", 
                            required=True,
                            metavar="port", 
                            type=int, 
                            help="the processor port")
    parser.add_argument("--worker-name", 
                            required=True,
                            metavar="worker_name", 
                            type=str, 
                            help="the worker name")
    args = parser.parse_args()
    w = Wrapper(args.lambda_path, args.port)
    w.start()