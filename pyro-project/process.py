import sys
import time
import random
from pyro import NetworkInterface
from enum import Enum, auto
import Pyro5.api

PEERS = [
    {"port": "3001", "id": "1", "name": "Node_1"},
    {"port": "3002", "id": "2", "name": "Node_2"},
    {"port": "3003", "id": "3", "name": "Node_3"},
    {"port": "3004", "id": "4", "name": "Node_4"}
]

COLORS = {
    "ok": "\033[92m",      # Green
    "warning": "\033[93m", # Yellow
    "critical": "\033[91m",# Red
    "reset": "\033[0m"
}


@Pyro5.api.expose
class ProcessState(Enum):
    FOLLOWER = auto()
    CANDIDATE = auto()
    LEADER = auto()

@Pyro5.api.expose
class Process:
    def __init__(self, name):
        self.name = name
        self.network = None 
        self.voted_for = None 
        self.votes = 0
        self.current_election = 0
        self.state = ProcessState.FOLLOWER
        self.election_timeout = self._get_random_timeout()
        self.last_contact = time.time()
        self.log = []
        self.commit_index = - 1

    def _get_random_timeout(self):
        return random.uniform(3, 4.5)
    
    def _reset_election_timer(self):
        self.last_contact = time.time()
        self.election_timeout = self._get_random_timeout()

    def _set_network(self, network_instance):
        self.network = network_instance

    def _check_status(self):
        if self.state != ProcessState.LEADER:
            if (time.time() - self.last_contact) > self.election_timeout:
                self._start_election()

    def decide_vote(self, candidate_term, candidate_name):
        if candidate_term > self.current_election:
            self.voted_for = candidate_name
            return True
        return False
    
    def append_log(self, entry):
        self.cprint("log registered....")
        self.log.append(entry) 
        return True

    def send_uncommited_logs(self, entry):
        total_votes = 1
        qtd_nodes = 1
        for peer in PEERS:                   
            if peer['name'] == self.name:
                continue
            time.sleep(0.25)
            try:
                total_votes += self.network.call_remote_method(
                    peer['port'], 
                    peer['id'], 
                    "append_log", 
                    entry
                )
                qtd_nodes += 1
            except:
                pass

        return total_votes, qtd_nodes

    def commit_log(self):
        time.sleep(0.2)
        self.cprint("log commited")
        self.commit_index = len(self.log) - 1
        self.print_node_state()

    def notify_commit(self):
        for peer in PEERS:                   
            if peer['name'] == self.name:
                continue
            
            try:
                self.network.call_remote_method(
                    peer['port'], 
                    peer['id'], 
                    "commit_log"
                )
            except:
                pass

    def execute(self, command):
        entry = {
            "term": self.current_election, 
            "command": command         
        }

        self.log.append(entry) 
        self.cprint("log registered....")
        time.sleep(0.5)
        received_acks, qtd_nodes = self.send_uncommited_logs(entry)
        
        if received_acks >= (qtd_nodes // 2) + 1:
            self.commit_log()
            time.sleep(0.5)
            self.notify_commit()

    def cprint(self, *args, **kwargs):
        state_map = {
            ProcessState.FOLLOWER: COLORS["ok"],      # Green
            ProcessState.CANDIDATE: COLORS["warning"], # Yellow
            ProcessState.LEADER: COLORS["critical"]    # Red
        }

        color_code = state_map.get(self.state, COLORS["reset"])
        message = " ".join(map(str, args))
        print(f"{color_code}{message}{COLORS['reset']}", **kwargs)

    def print_node_state(self):
        self.cprint(self.state)
        if self.commit_index != -1:
            self.cprint(self.log[self.commit_index])

    def receive_heartbeat(self, leader_term):
        if leader_term >= self.current_election:
            self.current_election = leader_term
            self.state = ProcessState.FOLLOWER
            self.last_contact = time.time()
            return True
        return False
        
    def become_leader(self):
        self.state = ProcessState.LEADER
        self.print_node_state()

    def send_heartbeat(self):
        for peer in PEERS:                   
            if peer['name'] == self.name:
                continue
            
            try:
                self.network.call_remote_method(
                    peer['port'], 
                    peer['id'], 
                    "receive_heartbeat", 
                    self.current_election
                )
            except:
                pass
                
        time.sleep(0.25)

    def _start_election(self):
        self.state = ProcessState.CANDIDATE
        self.votes = 1
        self.voted_for = self.name
        self.current_election += 1
        total_nodes = 1
        self.print_node_state()

        for peer in PEERS:
            try:
                if peer['name'] == self.name:
                    continue

                granted = self.network.call_remote_method(
                    peer['port'], 
                    peer['id'], 
                    "decide_vote", 
                    self.current_election, 
                    self.name
                )
                if granted:
                    self.votes += 1
                
                total_nodes += 1
            except:
                pass
        
        if self.votes >= (total_nodes // 2) + 1:
            self.become_leader()
            self.register_leader()
            return 
        
        self.state = ProcessState.FOLLOWER
        self.voted_for = None
        self.print_node_state()

    def register_leader(self):
        try:
            ns = Pyro5.api.locate_ns()
            try:
                ns.remove("Leader")
            except:
                pass

            ns.register("Leader", self.network.uri)
            self.cprint(f"[{self.name}] Leader registered on NameServer.")
            
        except Exception as e:
            self.cprint(f"[{self.name}] Failed to register leader on NameServer: {e}")

        
def raft_algorithm(process):
    if process.state == ProcessState.FOLLOWER:
        process._check_status()
    elif process.state == ProcessState.LEADER:
        process.send_heartbeat()

def create_process(node_id, port):
    process = Process(name=f"Node_{node_id}")
    process._set_network(NetworkInterface(
        object_to_expose=process, 
        port=port, 
        object_id=node_id
    ))
    return process

def main(node_id, port):
    process = create_process(node_id, port)
    process.network.start_communication()
    while True:
        raft_algorithm(process)
        time.sleep(0.01)    

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Uso: python process.py <id> <porta>")
    else:
        id = sys.argv[1]
        port = int(sys.argv[2])
        main(id, port)
