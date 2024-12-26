FROM scratch

COPY go-project-template /

ENTRYPOINT ["/go-project-template"]
USER 1000:100

# please see extra_files in .goreleaser.yaml if you need to copy additional files
COPY LICENSE /
