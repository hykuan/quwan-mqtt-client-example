build-image:
	docker build -t hykuan/quwan-mqtt-client-example:latest .

run:
	go run main.go

.PHONY: run build-image