Cluster Management (kind)

  kind create cluster --name k8s-counter    # Create cluster
  kind delete cluster --name k8s-counter    # Delete cluster
  kind get clusters                          # List clusters
  kind.exe export kubeconfig --name k8s-counter  # Set kubectl context (from WSL)

  Docker Image Loading

  docker build -t k8s-counter:local .       # Build image
  kind load docker-image k8s-counter:local --name k8s-counter  # Load into kind

  Kubernetes Resources (kubectl)

  kubectl apply -f k8s/                     # Apply all manifests
  kubectl get pods                          # List pods
  kubectl get deployments                   # List deployments
  kubectl get services                      # List services
  kubectl get hpa                           # List horizontal pod autoscalers

  Debugging

  kubectl logs <pod-name>                   # View pod logs
  kubectl logs -f <pod-name>                # Stream logs
  kubectl describe pod <pod-name>           # Pod details
  kubectl exec -it <pod-name> -- /bin/sh    # Shell into pod

  Scaling & Updates

  kubectl rollout restart deployment/counter  # Restart deployment
  kubectl scale deployment/counter --replicas=3  # Manual scale
  kubectl rollout status deployment/counter   # Watch rollout

  Cleanup

  kubectl delete pod <pod-name>             # Delete pod
  kubectl delete -f k8s/                    # Delete all resources

  Remember in WSL: Prefix commands with export KUBECONFIG=/mnt/c/Users/mitchell.durbin/.kube/config && or add it to your shell profile.

The load generator command for testing HPA autoscaling:

  kubectl run -i --tty load-generator --rm --image=busybox --restart=Never -- /bin/sh -c "while sleep 0.01; do wget -q -O- http://counter/count; done"

  This is already running in the background. To watch the HPA scale up:

  kubectl get hpa -w                    # Watch HPA metrics and replicas
  kubectl get pods -w                   # Watch pods scale up/down