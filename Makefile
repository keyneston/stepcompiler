.PHONY: all generate clean show test merge

all: clean generate show merge test

generate:
	go run github.com/keyneston/stepcompiler/generate

merge:
	mkdir -p step/
	cp _static/*.go step/
	cp _output/*.go step/

clean:
	rm -f _output/*.go
	rm -f _step/*.go

show:
	tail -n +1 _output/*.go

test:
	go test -cover github.com/keyneston/stepcompiler/step
