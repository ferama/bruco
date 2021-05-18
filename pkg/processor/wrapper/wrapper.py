import os
import sys
import json
import importlib
import signal
import socket
import logging

class Response:
    def __init__(self, key, data):
        self.key = key
        self.data = data

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
    def __init__(self, workdir: str, module_name: str, port: int, worker_name: str):
        self.port = port
        self.worker_name = worker_name
        self.module_name = module_name

        os.chdir(workdir)
        sys.path.append(".")
        
        signal.signal(signal.SIGINT, self.sigint_handler)

    def sigint_handler(self, p1, p2):
        sys.exit(0)

    def start(self):
        context = Context(self.worker_name)
        module = importlib.import_module(self.module_name)
        module.init_context(context)

        client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        client.connect((socket.gethostname(), self.port))

        while True:
            msg = client.recv(1024 * 1024 * 100)
            try:
                response = module.handle_event(context, msg)
                if not response: response = ""
                if type(response) == Response:
                    out = {
                        "key": response.key,
                        "data": response.data,
                        "error": ""
                    }
                else:
                    out = {
                        "key": "",
                        "data": response,
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
    parser.add_argument("--workdir", 
                            required=True,
                            metavar="workdir", 
                            type=str, 
                            help="the working directory")
    parser.add_argument("--module-name", 
                            required=True,
                            metavar="module_name", 
                            type=str, 
                            help="the module name")
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
    w = Wrapper(
            args.workdir, 
            args.module_name,
            args.port, 
            args.worker_name)
    w.start()