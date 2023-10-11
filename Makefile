init: pull build restart

pull:
	git pull

restart:
	docker-compose stop
	docker-compose up  --build --remove-orphans --detach

build:
	GOOS=linux GOARCH=amd64 go build -o laba_itmo_bot ./cmd