
# 启动zcat命令
socat -v tcp-l:18181,fork exec:"/bin/cat"
