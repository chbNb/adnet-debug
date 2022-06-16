# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start base image
FROM betapi/alpine-dev-golang:1.17-1 AS build_base
# Add Maintainer Info
ENV GOPRIVATE="gitlab.mobvistsa.com"
ENV GOSUMDB="off"
# RUN apt-get update && apt-get install -y libbrotli-dev
# Copy the predefined netrc file into the location that git depends on
# COPY ./netrc /root/.netrc
# RUN chmod 600 /root/.netrc
# Copy the predefined public deploy key file
COPY ./git_rsa_private_key /root/.ssh/id_rsa
COPY ./ssh_config /root/.ssh/config
RUN chmod 600 /root/.ssh/config /root/.ssh/id_rsa
RUN git config --global url.ssh://git@gitlab.mobvista.com/.insteadof http://gitlab.mobvista.com/
# Set the Current Working Directory inside the container
WORKDIR /go/src/gitlab.mobvista.com/ADN/adnet
COPY go.mod .
# Download dependencies
RUN go mod download

# Start from base build_base
FROM build_base AS server_builder
# Set the Current Working Directory inside the container
# WORKDIR /go/src/gitlab.mobvista.com/ADN/adnet
# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .
RUN go mod tidy
# Build the Go app
RUN make build && make cover

# Start a new stage from scratch
FROM betapi/alpine-prod:3.12.1 as aladdin
# Create the user and group files that will be used in the running 
# container to run the process as an unprivileged user.
# RUN mkdir /user && \
#     echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
#     echo 'nobody:x:65534:' > /user/group    
WORKDIR /root/
ENV CLOUD aws
ENV REGION virginia
ENV MODE product
# Copy the Pre-built binary file from the previous stage
COPY --from=server_builder /go/src/gitlab.mobvista.com/ADN/adnet/build/aladdin/bin/* /root/bin/
COPY --from=server_builder /go/src/gitlab.mobvista.com/ADN/adnet/build/aladdin/config /root/config/
COPY --from=server_builder /go/src/gitlab.mobvista.com/ADN/adnet/build/adnet_server/bin/* /root/bin/
COPY --from=server_builder /go/src/gitlab.mobvista.com/ADN/adnet/build/adnet_server/conf /root/conf
# RUN chmod +x /root/bin/aladdin_server
# EXPOSE 9102
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/bin/sh", "-c", "/root/bin/aladdin_server serve --cloud ${CLOUD} --config /root/config/ --region ${REGION} --mode ${MODE}"]
