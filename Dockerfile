FROM scratch

COPY bin/linux-amd64/iomize /

CMD ["./iomize"]