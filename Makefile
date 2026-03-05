user-rpc-local:
	@make -f deploy/mk/user-rpc.mk release-test

user-api-local:
	@make -f deploy/mk/user-api.mk release-test
	
social-rpc-local:
	@make -f deploy/mk/social-rpc.mk release-test

social-api-local:
	@make -f deploy/mk/social-api.mk release-test

im-api-local:
	@make -f deploy/mk/im-api.mk release-test

im-rpc-local:
	@make -f deploy/mk/im-rpc.mk release-test

im-ws-local:
	@make -f deploy/mk/im.ws.mk release-test

task-mq-local:
	@make -f deploy/mk/task-mq.mk release-test
	
release-local: user-rpc-local user-api-local social-rpc-local social-api-local im-api-local im-rpc-local im-ws-local task-mq-local


install-server:
	cd ./deploy/script && chmod +x release-local.sh && ./release-local.sh

install-server-user-rpc:
	cd ./deploy/script && chmod +x user-rpc-local.sh && ./user-rpc-local.sh