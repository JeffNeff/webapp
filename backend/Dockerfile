FROM golang:1.17-buster AS builder
WORKDIR /project
COPY . ./
RUN cd /project/cmd/webapp && go build -o /project/bin/

FROM registry.access.redhat.com/ubi8/ubi-minimal
EXPOSE 8080
ENV KO_DATA_PATH /kodata
COPY --from=builder /project/cmd/webapp/kodata/ ${KO_DATA_PATH}/
COPY --from=builder /project/bin/webapp /webapp

ENTRYPOINT ["/webapp"]
