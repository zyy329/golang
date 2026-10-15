FROM registry-intl.ap-southeast-1.aliyuncs.com/relax_many/openresty:1.25.3.1

ENV TZ=UTC

RUN opm get kuweiguge/lua-resty-nacos \
    && opm get SkyLothar/lua-resty-jwt \
    && opm get openresty/lua-resty-dns

WORKDIR /gateway

COPY . .

RUN mv lualib/httpUtils.lua /usr/local/openresty/site/lualib/resty/httpUtils.lua

COPY entrypoint.sh /usr/local/bin/entrypoint.sh

RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
