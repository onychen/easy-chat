ssh onychen@192.168.88.129

192.168.88.130

scp -r d:\code\go\easy-chat onychen@192.168.88.129:~/easy-chat

```
ssh onychen@192.168.88.129 "rm -rf ~/easy-chat"                         
scp -r d:\code\go\easy-chat onychen@192.168.88.129:~/easy-chat
make release-test

ssh onychen@192.168.88.130 "rm -rf ~/easy-chat"    
scp -r d:\code\go\easy-chat onychen@192.168.88.130:~/easy-chat
make install-server

```

```
docker login --username=onychen crpi-osk929019sdokpya.cn-guangzhou.personal.cr.aliyuncs.com
```

