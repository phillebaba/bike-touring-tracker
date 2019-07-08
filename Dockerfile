FROM scratch
ARG ARCH
COPY bike-touring-tracker-${ARCH} /bin/bike-touring-tracker
ENTRYPOINT ["/bin/bike-touring-tracker"]
