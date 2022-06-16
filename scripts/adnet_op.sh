#!/bin/bash

adnet_path=`pwd`
export GOROOT=$adnet_path/go
export GOPATH=$adnet_path/gopath
inpath=`echo $PATH | grep "$GOROOT" | wc -l`
if [ "$inpath" != 1 ]; then
    export PATH=$GOROOT/bin:$PATH
fi

export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$adnet_path/cgo/lib/:/usr/local/lib/
CONF_DEST="./conf/adnet_server.conf"
MGO_CONF_DEST="./conf/mgo.yaml"
if [ $# != 2 ] && [ "$3"x = ""x ]; then
    echo "FATAL: need two params"
    echo "Usage:  sh adnet_op.sh sg|vg|fk|se|sg_vta|vg_vta|fk_vta|se_vta start|restart|reload|stop"
    echo "        conf file need ${CONF_DEST}.sg or ${CONF_DEST}.vg or ${CONF_DEST}.fk"
    echo "ex:     sh adnet_op.sh sg start"
    exit 1
fi
REGION="$1"
if [ ! -f ${CONF_DEST}.$REGION ] || [ ! -f ${MGO_CONF_DEST}.$REGION ] ; then
    echo "FATAL: ${CONF_DEST}.$REGION or ${MGO_CONF_DEST}.$REGION not exist."
    exit 1
else
    cp ${CONF_DEST}.$REGION $CONF_DEST && cp ${MGO_CONF_DEST}.$REGION $MGO_CONF_DEST
fi

case $2 in 
    restart )
        killall supervise.adnet_server adnet_server check_adnet_server
        sh run $3 &>> /data/adn_logs/adnet/nohup_keep.log &
        echo "adnet_server has started"
        ;;
    reload )
        num=`ps axfu |grep './bin/adnet_server'|grep -v grep|grep -v supervise.adnet_server|grep -v check_adnet_server | wc -l`
        if [ $num -eq 0 ]; then
            sh run $3 &>> /data/adn_logs/adnet/nohup_keep.log &
            exit 0
        elif [ $num -eq 1 ]; then
            step=0
            while [ $num -eq 1 ] && [ $step -lt 10 ]
            do
                  sh reload
                  sleep 1
                  let step+=1
                  num=`ps axfu |grep './bin/adnet_server'|grep -v grep|grep -v supervise.adnet_server|grep -v check_adnet_server | wc -l`
                  echo "num: $num, step:$step "
            done
            exit 0
        fi
        echo "adnet_server has started, number is: $num"
        ;;
    start )
        num=`ps axfu | grep adnet_server | grep -v grep | wc -l`
        if [ $num -ge 1 ]; then
            echo "adnet_server or check_adnet_server or supervice has started"
            exit 0
        fi
        sh run $3 &>> /data/adn_logs/adnet/nohup_keep.log &
        ;;
    stop )
        killall supervise.adnet_server adnet_server check_adnet_server
        ;;
    * ) 
        echo "start | stop | reload | restart"
        ;;
esac
