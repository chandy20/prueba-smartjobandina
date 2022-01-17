.PHONY: production

A0= $(subst .,,$(suffix $@))
A1=$(basename $@)

LOG_DIR2=$(LOG_DIR)/picking

ENVS = production squad

LAMBDAS  = list create find box-price

$(foreach x,$(LAMBDAS),$(addsuffix .$x,$(ENVS))):
	@mkdir -p $(LOG_DIR2)
	$(MAKE) -C $(A0) $(A1) > $(LOG_DIR2)/deploy-$(A0).log 2>&1

production: $(addprefix production.,$(LAMBDAS))

$(addprefix logs.,$(LAMBDAS)):
	@cat $(LOG_DIR2)/deploy-$(A0).log || true

logs: $(addprefix logs.,$(LAMBDAS))
