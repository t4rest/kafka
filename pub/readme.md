docker build -t t4rest/go-pub:latest .

docker push t4rest/go-pub:latest

helm install -n go-pub ./helm-chart/


1. Get the application URL by running these commands:
  export NODE_PORT=$(kubectl get --namespace default -o jsonpath="{.spec.ports[0].nodePort}" services go-pub-go-pub)
  export NODE_IP=$(kubectl get nodes --namespace default -o jsonpath="{.items[0].status.addresses[0].address}")
  echo http://$NODE_IP:$NODE_PORT/projects


Run: helm ls --all go-pub; to check the status of the release
Or run: helm del --purge go-pub; to delete it
