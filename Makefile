.PHONY: production

A0= $(subst .,,$(suffix $@))
A1=$(basename $@)

ENVS = production
DOMAINS = beer

$(foreach x,$(DOMAINS),$(addsuffix .$x,$(ENVS))):
	@mkdir -p $(LOG_DIR)
	@echo $(A0) start
	@$(MAKE) -j2 -C $(A0) $(A1) > $(LOG_DIR)/deploy-$(A0).log 2>&1 
	@echo $(A0) ok

production: $(addprefix production.,$(DOMAINS))
	@echo "deployment finish"

$(addprefix logs.,$(DOMAINS)):
	@echo -e "\n\n\n\n\n\n$(A0)\n\n"
	@cat $(LOG_DIR)/deploy-$(A0).log || true
	@$(MAKE) -C $(A0) logs || true

logs: $(addprefix logs.,$(DOMAINS))
