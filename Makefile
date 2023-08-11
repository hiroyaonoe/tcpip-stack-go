.PHONY:run
run: fmt
	go run main.go

.PHONY:build
build: test
	go build

.PHONY:test
test: vet fmt
	go test ./...

.PHONY:test-with-coverage
test-with-coverage: vet fmt
	go test -cover ./... -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html

.PHONY:vet
vet:
	go vet ./...

.PHONY:fmt
fmt:
	go fmt ./...

.PHONY:curl
curl:
	curl --interface tun0 http://10.0.0.2/

.PHONY:capture
	tcpdump -i tun0 -w wireshark/capture.pcap
