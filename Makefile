help:
	@echo
	@echo "Usage: make TARGET"
	@echo
	@echo "Targets:"
	@echo "	static-up"
	@echo "	static-down"
	@echo

static-up:
	docker-compose -f docker-compose.static.yml up --build -d

static-down:
	docker-compose -f docker-compose.static.yml down

