FROM alpine:3.13
ADD ./ping-prober-config-generator /usr/local/bin/ping-prober-config-generator
ENTRYPOINT ["ping-prober-config-generator", "--config", "/etc/config/config.prod.yaml"]