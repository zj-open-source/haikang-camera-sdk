FROM hub-dev.rockontrol.com/docker.io/library/golang:1.16-buster

#  拷贝camera sdk库相关
ADD ./lib/ /opt/lib

ENV LD_LIBRARY_PATH=/opt/lib/64/:/opt/lib/32/:$LD_LIBRARY_PATH

