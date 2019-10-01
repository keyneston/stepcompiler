.PHONY: all generate clean show test merge

all: clean generate show merge test

generate:
	go run github.com/keyneston/stepcompiler/generate

merge:
	mkdir -p step/
	cp static/*.go step/
	cp output/*.go step/

clean:
	rm -f output/*.go
	rm -f step/*.go

show:
	tail -n +1 output/*.go

test:
	go test -cover github.com/keyneston/stepcompiler/step
