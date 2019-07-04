.env:
	./db/init.sh

bin/main: .env
	mkdir -p bin
	go build -o bin/main main.go

start: .env bin/main
	docker run \
		-d \
		-v $(PWD)/db:/var/lib/postgresql \
		-p 5432:5432 \
		postgres
	./bin/main &> logs.txt &

clean:
	rm -rf bin
	rm .env

shutdown:
	docker stop $(shell docker ps -l -q -f ancestor=postgres)
