# env-up: Start the demo environment
.PHONY: demo-up
demo-up:
	@echo Launching local k8s clusters
	@ctlptl apply -f cluster.yaml
	@echo Launching demo environment
	@tilt up --port 10350 --context kind-kafka-demo

# env-down: Stop the demo environment
.PHONY: demo-down
demo-down:
	@echo Stopping demo environment
	@kubectl config use-context kind-kafka-demo
	@tilt down --context kind-kafka-demo
	@echo Stopping local k8s clusters
	@ctlptl delete -f cluster.yaml

# env-reset: Reset demo environment
.PHONY: demo-reset
demo-reset: demo-down demo-up
