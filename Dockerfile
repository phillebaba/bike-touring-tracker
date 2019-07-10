FROM alpine:3.9.4
ARG ARCH
COPY ./bike-touring-tracker-${ARCH} /bin/bike-touring-tracker
CMD ["/bin/bike-touring-tracker"]
