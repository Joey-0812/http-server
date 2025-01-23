#go环境
FROM ubuntu:20.04

# 在容器内设置 /workspace 为当前工作目录
RUN mkdir /workspace
WORKDIR /workspace
# 把文件复制到当前工作目录
COPY ./http-server .
RUN ls
COPY ./server.crt .
COPY ./server.key .
RUN ls
EXPOSE 443
EXPOSE 80
ENTRYPOINT [ "./http-server" ]