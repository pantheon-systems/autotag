# common make tasks and variables that should be imported into all projects
#
#-------------------------------------------------------------------------------
help: ## print list of tasks and descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?##"}; { split($$0,a,":"); printf "\033[36m%-30s\033[0m %s \n", a[2], $$2}'
.DEFAULT_GOAL := help

readme-toc: ## update the Table of Contents in ./README.md (replaces <!-- toc --> tag)
	docker run --rm -v `pwd`:/src quay.io/getpantheon/markdown-toc -i /src/README.md


.PHONY:: all help update-makefiles
