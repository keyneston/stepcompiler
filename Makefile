.PHONY: all generate clean show test

all: clean generate test

generate:
	go run github.com/keyneston/stepcompiler/generate

clean:
	rm -f step/gen_*.go

show:
	tail -n +1 step/gen_*.go

test:
	go test -cover ./...
