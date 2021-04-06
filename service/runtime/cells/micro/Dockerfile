FROM micro/micro

WORKDIR /

# This is just a hack to download most of the dependencies micro and go micro has
# We don't care it's not the latest version, just an approximation
RUN mkdir go && cd go && git clone https://github.com/micro/micro && \
    cd micro && go install

COPY entrypoint.sh /
RUN chmod 755 entrypoint.sh
ENTRYPOINT ["bash", "/entrypoint.sh"]
