db/data:
	docker run \
		-d \
		-v $(PWD)/db:/var/lib/postgresql \
		-p 5432:5432 \
		postgres

db/.postgres: db/data
	./db/init.sh

bin/main: db/data db/.postgres
	mkdir -p bin
	go build -o bin/main main.go

start: db/.postgres bin/main
	./bin/main &> logs.txt &

clean:
	rm -rf bin
	rm ./db/.postgres

shutdown:
	docker stop $(shell docker ps -l -q -f ancestor=postgres)
