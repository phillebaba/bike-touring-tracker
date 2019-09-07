FROM alpine:3.10.2
ARG TARGETARCH
COPY ./bike-touring-tracker-${TARGETARCH} /bin/bike-touring-tracker
CMD ["/bin/bike-touring-tracker"]
