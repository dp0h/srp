#Import and expose environment variables
cnf ?= .env
include $(cnf)
export $(shell sed 's/=.*//' $(cnf))
# Import env from parent
cnf_priv ?= ../.env
include $(cnf_priv)
export $(shell sed 's/=.*//' $(cnf_priv))


help:
	@echo
	@echo "Usage: make TARGET"
	@echo
	@echo "Targets:"
	@echo "	auto-up"
	@echo "	auto-down"
	@echo "	static-up"
	@echo "	static-down"
	@echo

auto-up:
	docker-compose -f docker-compose.auto.yml up --build -d

auto-down:
	docker-compose -f docker-compose.auto.yml down

static-up:
	docker-compose -f docker-compose.static.yml up --build -d

static-down:
	docker-compose -f docker-compose.static.yml down

http-up:
	docker-compose -f docker-compose.http.yml up --build -d

http-down:
	docker-compose -f docker-compose.http.yml down
