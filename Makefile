init: pull build restart
update: init

pull:
	git pull

restart:
	docker compose stop
	docker compose up  --build --remove-orphans --detach

build:
	GOOS=linux GOARCH=amd64 go build -o laba_itmo_bot ./cmd