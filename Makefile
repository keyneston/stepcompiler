.PHONY: all generate clean show

all: clean generate show

generate:
	go run generate/*.go

clean:
	rm -f output/*.go

show:
	tail -n +1 output/*.go

