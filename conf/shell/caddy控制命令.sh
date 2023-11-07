# 发送配置文件
curl localhost:2019/load -H "Content-Type: application/json" -d @reverse_proxy.json

# 修改日志格式为debug
curl localhost:15520/config/logging/logs/default/level -H "Content-Type: application/json" -d '"DEBUG"'

# 加载刷新证书
curl localhost:15520/config/apps/tls/certificates/load_files -H "Content-Type: application/json" -d '[{"certificate": "/data/hiberlin/xuqiu/ioa/ngn_smart_web/conf/_wildcard.baidu.com.pem","key": "/data/hiberlin/xuqiu/ioa/ngn_smart_web/conf/_wildcard.baidu.com-key.pem"}]'

#加载刷新tls域名
curl localhost:15520/config/apps/http/servers/https_server/routes/0/match -H "Content-Type: application/json" -d '[{"host":["*.baidu.com", "*.qq.com,  "*.tencent.com"]}]'
