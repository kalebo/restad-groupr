DATAPATH=data/...
GOBINDATA=~/gocode/bin/go-bindata
GOBINDATA_FLAGS= -o assets.go -nomemcopy

debug:
	$(GOBINDATA) -debug $(GOBINDATA_FLAGS) $(DATAPATH)
	go build

release:
	$(GOBINDATA) $(GOBINDATA_FLAGS) $(DATAPATH)
	go build

setup:
	go get github.com/gorilla/mux
	go get github.com/golang/glog
	go get gopkg.in/cas.v1
	go get -u github.com/jteeuwen/go-bindata/...
