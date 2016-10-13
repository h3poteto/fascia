FROM node:6.7-slim

RUN useradd -m -s /bin/bash node
RUN echo 'node:password' | chpasswd
RUN mkdir -p /var/opt/app

USER node

WORKDIR /var/opt/app

CMD ["/bin/bash"]
