REGION ?= us-east-1
PROFILE ?= sa-infra
ENV_NAME ?= dev

WAAS_PROXY_STACK_NAME ?= sa-waas-$(ENV_NAME)

CLUSTER ?= $(WAAS_PROXY_STACK_NAME)

WAAS_PROXY_SVC ?= $(WAAS_PROXY_STACK_NAME)

ACCOUNT_ID := $(shell aws sts get-caller-identity --profile $(PROFILE) --query 'Account' --output text)

.PHONY: create-waas-stack
create-waas-stack:
	@aws cloudformation create-stack \
	--profile $(PROFILE) \
	--stack-name $(WAAS_PROXY_STACK_NAME) \
	--region $(REGION) \
	--capabilities CAPABILITY_NAMED_IAM \
	--template-body file://waas.cfn.yml \
	--parameters file://waas.json

.PHONY: delete-waas-stack
delete-waas-stack:
	@aws cloudformation delete-stack \
  --profile $(PROFILE) \
  --stack-name $(WAAS_PROXY_STACK_NAME) \
  --region $(REGION)

.PHONY: validate-waas-template
validate-waas-template:
	@aws cloudformation validate-template \
  --profile $(PROFILE) \
  --template-body file://waas.cfn.yml \
  --region $(REGION)

.PHONY: update-waas-stack
update-waas-stack:
	@aws cloudformation update-stack \
  --profile $(PROFILE) \
  --stack-name $(WAAS_PROXY_STACK_NAME) \
  --region $(REGION) \
  --capabilities CAPABILITY_NAMED_IAM \
  --template-body file://waas.cfn.yml \
	--parameters file://waas.json

.PHONY: update-waas-service
update-waas-service:
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go
	@aws ecr get-login-password \
  --profile $(PROFILE) \
  --region $(REGION) \
	| docker login --username AWS --password-stdin $(ACCOUNT_ID).dkr.ecr.$(REGION).amazonaws.com
	@docker build -t $(WAAS_PROXY_SVC) .
	@docker tag $(WAAS_PROXY_SVC):latest $(ACCOUNT_ID).dkr.ecr.$(REGION).amazonaws.com/$(WAAS_PROXY_SVC):latest
	@docker push $(ACCOUNT_ID).dkr.ecr.$(REGION).amazonaws.com/$(WAAS_PROXY_SVC):latest
	@aws ecs update-service \
  --profile $(PROFILE) \
  --region $(REGION) \
  --cluster $(CLUSTER) \
  --service $(WAAS_PROXY_SVC) \
  --force-new-deployment

.PHONY: start-local
start-local:
	@go run cmd/server/main.go