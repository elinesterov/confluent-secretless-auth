# env-up: Start the dev demo environment
.PHONY: env-up
env-up:
	@echo Launching local k8s clusters
	@ctlptl apply -f cluster.yaml
	@echo Launching demo environment
	@tilt up --port 10350 --context kind-kafka-test

# env-down: Stop the dev demo environment
.PHONY: env-down
env-down:
	@echo Stopping demo environment
	@kubectl config use-context kind-kafka-test
	@tilt down --context kind-kafka-test
	@echo Stopping local k8s clusters
	@ctlptl delete -f cluster.yaml

# env-reset: Reset dev demo environment
.PHONY: env-reset
env-reset: env-down env-up
