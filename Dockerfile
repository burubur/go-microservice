# Stage build
FROM registry.kudoserver.com/kudo/golang:1.11 as builder

ENV SERVICENAME="biller_connector"
ENV VERSION="0.5.0"

WORKDIR /go/src/bitbucket.org/kudoindonesia/

# Add the keys and set permissions
RUN apk add --no-cache openssh
ARG SSH_PRIVATE_KEY
RUN mkdir ~/.ssh/ && \
    echo "${SSH_PRIVATE_KEY}" > ~/.ssh/id_rsa && \
    chmod 0600 ~/.ssh/id_rsa && \
    touch ~/.ssh/known_hosts && \
    ssh-keyscan bitbucket.org >> ~/.ssh/known_hosts

# copy source
WORKDIR /go/src/bitbucket.org/kudoindonesia/$SERVICENAME
COPY . .

RUN dep ensure -v && CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -o $SERVICENAME -tags static_all .

# Stage Runtime Applications
FROM registry.kudoserver.com/kudo/base-image:1.3.3
LABEL name=$SERVICENAME
LABEL version=$VERSION

# Download Depedencies
RUN apk update && apk add --no-cache ca-certificates bash jq curl && rm -rf /var/cache/apk/*

# Setting timezone
ENV TZ=Asia/Jakarta
RUN apk add --no-cache tzdata
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV BUILDDIR /go/src/bitbucket.org/kudoindonesia/$SERVICENAME

# Add user kudo
RUN adduser -D kudo kudo

# Setting folder workdir
WORKDIR /opt/$SERVICENAME

# Copy Data App
COPY --from=builder $BUILDDIR/$SERVICENAME $SERVICENAME

# Setting owner file and dir
RUN chown -R kudo:kudo .

USER kudo

EXPOSE 8090

CMD [ "" ]

# TODO:
# background: we MUST install a specific version for each dependency
# issue: using a version on each dependency caused an error on building stage in staging
# action: discuss with tribe lead and devops team about each os dependency version
