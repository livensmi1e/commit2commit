.PHONY: build run

BIN := commit2commit.exe

build:
	go build -o $(BIN)

run: build
	@./$(BIN) $(ARGS)
