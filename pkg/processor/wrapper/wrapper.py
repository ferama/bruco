import os
import sys
import json
import importlib
import signal
import socket
import logging

class Context: 
    def __init__(self, worker_name):
        root = logging.getLogger()
        root.setLevel(logging.DEBUG)

        handler = logging.StreamHandler(sys.stdout)
        handler.setLevel(logging.DEBUG)
        formatter = logging.Formatter(f"%(asctime)s ({worker_name}: %(levelname)s) %(message)s")
        handler.setFormatter(formatter)
        root.addHandler(handler)
        self.worker_name = worker_name
        self.logger = logging


class Wrapper:
    def __init__(self, lambda_path: str, port: int, worker_name: str):
        self.port = port
        self.worker_name = worker_name

        os.chdir(lambda_path)
        sys.path.append(".")
        
        signal.signal(signal.SIGINT, self.sigint_handler)

    def sigint_handler(self, p1, p2):
        sys.exit(0)

    def start(self):
        context = Context(self.worker_name)
        module = importlib.import_module("handler")
        module.init_context(context)

        client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        client.connect((socket.gethostname(), self.port))

        while True:
            msg = client.recv(1024 * 1024 * 100)
            try:
                response = module.handle_event(context, msg)
                if not response: response = ""
                out = {
                    "data": response,
                    "error": ""
                }
                out = json.dumps(out)
                out += "\n"
                client.send(out.encode())
            except Exception as e:
                out = {
                    "data": "",
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
    w = Wrapper(args.lambda_path, args.port, args.worker_name)
    w.start()