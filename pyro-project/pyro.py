import Pyro5.api
import threading

class NetworkInterface:
    def __init__(self, object_to_expose, port, object_id):
        self.logic_object = object_to_expose
        self.port = port
        self.object_id = object_id
        self.daemon = None
        self.uri = None

    def start_communication(self):
        Pyro5.api.expose(self.logic_object.__class__) 
        
        self.daemon = Pyro5.api.Daemon(port=self.port)
        self.uri = self.daemon.register(self.logic_object, objectId=self.object_id)
        
        threading.Thread(target=self.daemon.requestLoop, daemon=True).start()
        print(f"Node {self.object_id} ready and waiting on port {self.port}...")

    def call_remote_method(self, target_port, target_id, method_name, *args):
        uri = f"PYRO:{target_id}@localhost:{target_port}"
        
        with Pyro5.api.Proxy(uri) as proxy:
            remote_method = getattr(proxy, method_name)
            return remote_method(*args)

