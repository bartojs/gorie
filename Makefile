build:
	export GOPATH=`pwd` && go install gorie
	export GOPATH=`pwd` && export GOARCH='amd64' && export GOOS="windows" && go install gorie
	export GOPATH=`pwd` && export GOARCH='amd64' && export GOOS="linux" && go install gorie

gen:
	export GOPATH=`pwd` && go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
	export PATH=$$PATH:./bin && protoc --go_out=src/gorie riemann.proto

setup:
	brew install protobuf
	curl -o riemann.proto https://raw.githubusercontent.com/aphyr/riemann-java-client/master/src/main/proto/riemann/proto.proto

run:
	bin/gorie
