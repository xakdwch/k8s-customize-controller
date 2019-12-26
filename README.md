# k8s-customize-controller
simple customize controller
# compile and run
cd k8s-customize-controller && go build

run "kubectl create -f programmer.yaml" to create programmer crd

run "./k8s-customize-controller -kubeconfig=/root/.kube/config -alsologtostderr=true" to start programmer controller
# test
run "kubectl create -f deployment-programmer.yaml" to create programmer deployment
