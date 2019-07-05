# TODO: scp backed make for integrated cross machine deploys

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

upload: bundle
	ssh ubuntu@$(SERVER_HOST) rm -rf ~/
	scp bundle.zip ubuntu@$(SERVER_HOST):~/
	ssh ubuntu@$(SERVER_HOST) unzip bundle.zip

deploy: upload
	ssh ubuntu@$(SERVER_HOST) docker-compose restart
