
# http 从2080端口转到9000端口
caddy reverse-proxy --from :2080 --to :9000

# https localhost:443 转到 9000端口
caddy reverse-proxy --to :9000

# https example.com:8443 转到 9000端口
caddy reverse-proxy --from example.com:8443 --to :9000

# https localhost:2080 转到 https://localhost:9000
caddy reverse-proxy --from :2080 --to https://localhost:9000

# https example.com:443 转到 https://example.com:9000
caddy reverse-proxy --from example.com --to https://example.com:9000

# https example.com:443 转到 https://localhost:9000
caddy reverse-proxy --from example.com --to https://localhost:9000 --change-host-header