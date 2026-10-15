local _M = {}

function _M.init()
    local resolver = require "resty.dns.resolver"
    local r, err = resolver:new{
        nameservers = {nameserver},
        retrans = 5,  -- 5 retransmissions on receive timeout
        timeout = 2000,  -- 2 sec
    }

    if not r then
        ngx.say("failed to instantiate the resolver: ", err)
        return
    end
    return r
end

function _M.getIp(host, default)
    if not nameserver then
        ngx.log(ngx.ERR, "get Ip nameserver nil")
        return default
    end

    local r = _M.init()
    local answers, err, tries = r:query(host, nil, {})
    if not answers then
        ngx.log(ngx.ERR, "failed to query the DNS server: ", err, "\n retry historie:\n  ", table.concat(tries, "\n  "))
        return default
    end

    if answers.errcode then
        --ngx.say("server returned error code: ", answers.errcode, ": ", answers.errstr)
        return default
    end

    for i, ans in ipairs(answers) do
        return ans.address or ans.cname
        --ngx.say(ans.name, " ", ans.address or ans.cname,
        --        " type:", ans.type, " class:", ans.class,
        --        " ttl:", ans.ttl)
    end

    return default
end

return _M