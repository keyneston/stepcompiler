.PHONY: all generate clean show test

all: clean generate show test

generate:
	cp static/*.go output/
	go run generate/*.go

clean:
	rm -f output/*.go

show:
	tail -n +1 output/*.go

test:
	go test ./...
