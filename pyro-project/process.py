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

    def _get_random_timeout(self):
        return random.uniform(1.5, 3.0)
    
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
    
    def receive_heartbeat(self, leader_term):
        print("receiving heartbeat")
        if leader_term >= self.current_election:
            self.current_election = leader_term
            self.state = ProcessState.FOLLOWER
            self.last_contact = time.time()
            return True
        return False
        
    
    def become_leader(self):
        self.state = ProcessState.LEADER

    def send_heartbeat(self):
        print("sending heartbeat...")
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

        for peer in PEERS:
            granted = self.network.call_remote_method(
                peer['port'], 
                peer['id'], 
                "decide_vote", 
                self.current_election, 
                self.name
            )
            if granted:
                self.votes += 1
            
            total_nodes = len(PEERS) + 1
            if self.votes >= (total_nodes // 2) + 1:
                self.become_leader()
                return 

def raft_algorithm(process):
    if process.state == ProcessState.FOLLOWER:
        process._check_status()
    elif process.state == ProcessState.LEADER:
        process.send_heartbeat()
    else:
        print("candidate state")

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
        time.sleep(0.1)    

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Uso: python process.py <id> <porta>")
    else:
        id = sys.argv[1]
        port = int(sys.argv[2])
        main(id, port)
