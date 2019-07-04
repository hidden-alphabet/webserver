.env:
	./db/init.sh

bin/main: .env
	mkdir -p bin
	go build -o bin/main main.go

start: .env bin/main
	./bin/main &> logs.txt &

clean:
	rm -rf bin
	rm .env

shutdown:
	docker stop $(shell docker ps -l -q -f ancestor=postgres)
