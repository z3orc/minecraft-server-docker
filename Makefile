all: build

test: build
	wget https://fill-data.papermc.io/v1/objects/a61a0585e203688f606ca3a649760b8ba71efca01a4af7687db5e41408ee27aa/paper-1.21.10-117.jar -O ./out/server.jar
	cd ./out/ && ./runner

build: clean
	mkdir ./out/
	go build -o ./out/runner

clean:
	rm -rf ./out
