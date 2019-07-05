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

bundle:
	zip -r bundle.zip services
	zip -r bundle.zip scripts
	zip bundle.zip docker-compose.yml
	zip bundle.zip Makefile
	zip bundle.zip package-lock.json
	zip bundle.zip package.json
	zip bundle.zip webpack.config.js