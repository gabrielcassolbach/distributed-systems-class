PYTHON = python3
SCRIPT = process.py
CLIENT_SCRIPT = client.py
NODES_COUNT = 4
BASE_PORT = 3001

TERM_CMD = xterm -bg black -fg white -fa Monospace -fs 11 -e bash -c

run:
	@echo "Subindo o Servidor de Nomes..."
	$(TERM_CMD) "$(PYTHON) -m Pyro5.nameserver; exec bash" &
	@sleep 2
	@echo "Subindo $(NODES_COUNT) nós do Raft..."
	@for i in $$(seq 1 $(NODES_COUNT)); do \
		port=$$(($(BASE_PORT) + i - 1)); \
		echo "Iniciando Nó $$i na porta $$port..."; \
		$(TERM_CMD) "$(PYTHON) $(SCRIPT) $$i $$port; exec bash" & \
	done
	@sleep 1
	@echo "Iniciando Cliente Interativo em nova janela..."
	$(TERM_CMD) "$(PYTHON) $(CLIENT_SCRIPT); exec bash" &

kill:
	@fuser -k $(BASE_PORT)/tcp || true
	@fuser -k $$(($(BASE_PORT) + 1))/tcp || true
	@fuser -k $$(($(BASE_PORT) + 2))/tcp || true
	@fuser -k $$(($(BASE_PORT) + 3))/tcp || true
	@fuser -k 9090/tcp || true
	@pkill -9 xterm || true