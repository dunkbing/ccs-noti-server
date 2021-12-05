BINARY_NAME=ccs-noti-server

build:
	go build -o ${BINARY_NAME} .

run:
	./${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}

deploy:
	gcloud app deploy