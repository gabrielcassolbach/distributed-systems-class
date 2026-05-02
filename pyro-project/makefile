PYTHON = python3
SCRIPT = process.py
NODES_COUNT = 4
BASE_PORT = 3001

TERM_CMD = xterm -bg black -fg white -fa Monospace -fs 11 -e bash -c

run:
	@echo "Subindo $(NODES_COUNT) nós do Raft..."
	@for i in $$(seq 1 $(NODES_COUNT)); do \
		port=$$(($(BASE_PORT) + i - 1)); \
		echo "Iniciando Nó $$i na porta $$port..."; \
		$(TERM_CMD) "$(PYTHON) $(SCRIPT) $$i $$port; exec bash" & \
	done

kill:
	@pkill -f $(SCRIPT) || true
