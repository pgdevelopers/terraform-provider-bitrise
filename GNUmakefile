default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

build:
	mkdir -p ~/.terraform.d/plugins/terraform.local/local/bitrise/1.0.0/darwin_arm64
	go build -o terraform-provider-bitrise
	chmod +x terraform-provider-bitrise
	mv terraform-provider-bitrise ~/.terraform.d/plugins/terraform.local/local/bitrise/1.0.0/darwin_arm64/terraform-provider-bitrise_v1.0.0