DATAPATH=react-example/... 

debug:
	go-bindata-assetfs -debug $(DATAPATH)
	go build

release:
	go-bindata-assetfs -nomemcopy $(DATAPATH)
	go build

setup:
	go get github.com/gorilla/mux
	go get github.com/golang/glog
	go get gopkg.in/cas.v1
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/elazarl/go-bindata-assetfs/...
