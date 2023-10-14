FROM golang:1.21.3-bullseye as builder
WORKDIR /usr/src/sockmon

# Copy librespeed-cli
COPY . .

# Build librespeed-cli
RUN make 

FROM debian:bullseye-slim
RUN apt-get update && apt-get install iproute2 -y && apt-get clean && rm -rf /var/lib/apt/lists/*

# Copy librespeed-cli binary
COPY --from=builder /usr/src/sockmon/bin/sockmon /bin/sockmon

ENTRYPOINT ["/bin/sockmon"]