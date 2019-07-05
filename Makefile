build:
	docker build -t hiddenalphabet-api-server .

start: .env bin/main
	docker-compose up -d
	./services/postgresql/init.sh

shutdown:
	docker-compose down

clean:
	rm -rf bin
	rm .env