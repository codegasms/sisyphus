all: sisyphus summit zeus

.PHONY: sisyphus summit zeus

bin:
	mkdir bin

sisyphus: . bin
	go build -o bin/$@ ./$<

summit: cmd/summit bin
	go build -o bin/$@ ./$<

zeus: cmd/zeus bin
	go build -o bin/$@ ./$<

fmt:
	go fmt -x ./...

clean:
	$(RM) -- ./bin/*
