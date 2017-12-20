FROM golang:1.8.3-onbuild
ADD . /app
WORKDIR /app
RUN ["dnu", "restore"]

EXPOSE 5004
ENTRYPOINT ["dnx", "./src/HelloMvc6", "kestrel"]
