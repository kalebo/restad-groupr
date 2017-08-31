DATAPATH=dist/... 

debug:
	go-bindata-assetfs -debug $(DATAPATH)
	go build

release:
	go-bindata-assetfs -nomemcopy $(DATAPATH)
	go build

setup:
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/elazarl/go-bindata-assetfs/...
	go get "github.com/go-zoo/bone"
	go get "github.com/golang/glog"
	go get "github.com/mattn/go-sqlite3"
	go get "gopkg.in/cas.v2"