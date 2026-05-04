import Pyro5.api

class Client:
    def __init__(self, name):
        self.name = name 

    def run(self):
        while True:
            command = input("Enter a command: ")
            try:
                ns = Pyro5.api.locate_ns()
                leader_url = ns.lookup("Leader")

                leader = Pyro5.api.Proxy(leader_url)
                answer = leader.execute(command)
                
                print(answer)
            except Exception as e:
                print(f"[{e}]")
            
def main():
    client = Client("client")
    client.run()
    
if __name__ == "__main__":
    main()