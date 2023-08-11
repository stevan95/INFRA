resource "helm_release" "argocd" {
  name = "argocd"

  repository       = "https://argoproj.github.io/argo-helm"
  chart            = "argo-cd"
  namespace        = "argocd"
  version          = "5.41.2"
  create_namespace = true

  set {
    name = "server.service.type"
    value = "LoadBalancer"
  }

  values = [file("values-argocd.yaml")]
}