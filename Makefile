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

.PHONY:create-tap
create-tap:
	ip tuntap add mode tun dev tap0 &&\
	ip link set tap0 up

.PHONY:add-addr
add-addr:
	ip addr add 10.0.0.1/24 dev tap0

.PHONY:delete-tap
delete-tap:
	ip link del dev tap0

.PHONY:curl
curl:
	curl --interface tap0 http://10.0.0.2/

DATE := $(shell date +%Y-%m-%d_%H-%M-%S)

.PHONY:capture
capture:
	tcpdump -i tap0 -w wireshark/capture_$(DATE).pcap

