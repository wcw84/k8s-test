APP_NAME=access-http
SRC=$(wildcard *.go)
EXE=app
$(EXE):$(SRC)
	docker run --rm -it --name go-compile -v ${PWD}:/go/src --workdir /go/src golang:alpine go build -o $(EXE)
	docker build -t k8s-test/$(APP_NAME) .
# 	docker save k8s-test/$(APP_NAME) -o ./$(APP_NAME).tar
# 	scp $(APP_NAME).tar root@master:/root && ssh root@master "docker load < /root/$(APP_NAME).tar"
# 	scp $(APP_NAME).tar root@node1:/root && ssh root@node1 "docker load < /root/$(APP_NAME).tar"
# 	scp $(APP_NAME).tar root@node2:/root && ssh root@node2 "docker load < /root/$(APP_NAME).tar"
#     rm ./$(APP_NAME).tar
