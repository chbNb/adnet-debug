#!/bin/bash

# adnet_path=`pwd`
# export GOPAT
# export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib

#packages=`go list ./src/adnet/...|grep -vE '^adnet$|/adnet|/ad_server|/corsair_proto|/enum'`
#packages=`go list gitlab.mobvista.com/ADN/adnet/internal/...|grep -vE '^adnet$|/ad_server|/corsair_proto|/enum|/discovery|/netacuity|/watcher|/redis|/backender|/thrift_pool|/protobuf|/rbaserviceapi|/as_consul|/ad_server|/extractor'`
#packages=`go list ./src/adnet/...|grep -vE '^adnet$|/adnet|/ad_server|/corsair_proto|/enum|/discovery|/netacuity|/watcher|/redis'`

if [ "$1" == "coverage" ];then
    # for jenkins ci
    echo "run go test with junit report and coverage analysis..."
    #go test -test.short -v $packages -gcflags=-l 2>&1|go-junit-report> $WORKSPACE/adnet-report.xml
    #gocov test -test.short $packages -gcflags=-l -coverpkg=`echo $packages|sed 's/ /,/g'` | gocov-xml > $WORKSPACE/adnet-coverage.xml
    cd internal
    go test -test.short -v ./... -gcflags=-l -coverprofile=$WORKSPACE/adnet-cover.out 2>&1 | go-junit-report > $WORKSPACE/adnet-report.xml
    gocov convert $WORKSPACE/adnet-cover.out | gocov-xml > $WORKSPACE/adnet-coverage.xml
else
    echo "running go test"
    go test -test.short $packages -gcflags=-l
fi
