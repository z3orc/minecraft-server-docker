all: build

get-jar: build
	wget https://fill-data.papermc.io/v1/objects/a61a0585e203688f606ca3a649760b8ba71efca01a4af7687db5e41408ee27aa/paper-1.21.10-117.jar -O ./out/server.jar

run: build
	cd ./out/ && WHITE_LIST=true ./runner -jar server.jar -timeout 10 -debug -sigkill

build:
	mkdir -p ./out/
	go build -o ./out/runner

build-docker:
	docker build -t z3orc/minecraft-server-docker .

clean:
	rm -rf ./out
