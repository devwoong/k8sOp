FROM ubuntu:20.04

RUN mkdir /home/server

WORKDIR /home/server

COPY k8sOp .

COPY apply/ .

ENTRYPOINT [ "/home/server/k8sOp" ]

# EXPOSE 80