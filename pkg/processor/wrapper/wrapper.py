import os
import pathlib
import sys
import json
import importlib
import signal
import socket
import logging

class Response:
    def __init__(self, key, data, content_type = "application/json"):
        self.key = key
        self.data = data
        self.content_type = content_type

class Context: 
    def __init__(self, worker_name):
        root = logging.getLogger()
        root.setLevel(logging.DEBUG)

        handler = logging.StreamHandler(sys.stdout)
        handler.setLevel(logging.DEBUG)
        formatter = logging.Formatter(f"%(asctime)s ({worker_name}): %(levelname)s %(message)s", "%Y/%m/%d %H:%M:%S")
        handler.setFormatter(formatter)
        root.addHandler(handler)

        self.logger = logging
        self.worker_name = worker_name
        self.response = Response


class Wrapper:
    def __init__(self, handler_path: str, module_name: str, socket: str, worker_name: str):
        self.socketPath = socket
        self.worker_name = worker_name

        os.chdir(handler_path)
        sys.path.append(".")
        self.module_name = module_name
            
        signal.signal(signal.SIGINT, self.sigint_handler)

    def sigint_handler(self, p1, p2):
        sys.exit(0)

    def start(self):
        context = Context(self.worker_name)
        module = importlib.import_module(self.module_name)
        if hasattr(module, "init_context"):
            module.init_context(context)

        client = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        client.connect(self.socketPath)
        while True:
            msg = client.recv(1024 * 1024 * 100)
            try:
                response = module.handle_event(context, msg)
                if not response: response = ""
                if type(response) == Response:
                    out = {
                        "key": str(response.key),
                        "data": response.data,
                        "contentType": response.content_type,
                        "error": ""
                    }
                else:
                    out = {
                        "key": "",
                        "data": response,
                        "contentType": "application/json",
                        "error": ""
                    }
                out = json.dumps(out)
                out += "\n"
                client.send(out.encode())
            except Exception as e:
                out = {
                    "data": "",
                    "error": str(e)
                }
                out = json.dumps(out)
                out += "\n"
                client.send(out.encode())
        

if __name__ == "__main__":
    import argparse
    parser = argparse.ArgumentParser(description="the python shell")
    parser.add_argument("--handler-path", 
                            required=True,
                            metavar="handler_path", 
                            type=str, 
                            help="the function directory")
    parser.add_argument("--module-name", 
                            required=True,
                            metavar="module_name", 
                            type=str, 
                            help="the module name")
    parser.add_argument("--socket", 
                            required=True,
                            metavar="socket", 
                            type=str, 
                            help="the processor socket path")
    parser.add_argument("--worker-name", 
                            required=True,
                            metavar="worker_name", 
                            type=str, 
                            help="the worker name")
    args = parser.parse_args()
    w = Wrapper(
            args.handler_path, 
            args.module_name,
            args.socket,
            args.worker_name)
    w.start()