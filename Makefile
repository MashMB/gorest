app = app

all: run

build: clean
	@echo "--- Building ---"
	mkdir bin
	go build -o ./bin/$(app) ./cmd/$(app)/main.go

clean:
	@echo "--- Cleaning ---"
	rm -rf bin

run: build
	@echo "--- Running ---"
	./bin/$(app)
