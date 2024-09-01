createcluster:
	kind create cluster --name=goexpert

deletecluster:
	kind delete cluster --name=goexpert

clusterinfo:
	kubectl cluster-info --context kind-goexpert

getnodes:
	kubectl get nodes

getpods:
	kubectl get pods

getsvc:
	kubectl get svc

applydeployment:
	kubectl apply -f k8s/deployment.yml

applyservice:
	kubectl apply -f k8s/service.yml

fowardport:
	kubectl port-forward svc/serversvc 8080:8080

.PHONY: createcluster  deletecluster clusterinfo getnodes getpods applydeployment applyservice getsvc fowardport