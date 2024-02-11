resource "kubernetes_namespace" "argocd_namespace" {
  metadata {
    name = "argocd"
  }
}

resource "helm_release" "argocd" {
  repository            = "https://argoproj.github.io/argo-helm"
  name                  = "argo-cd"
  chart                 = "argo-cd"
  version               = "6.0.5"
  values = [
    "${file("argocd-values.yaml")}"
  ]
  namespace = "argocd"

  depends_on = [
    helm_release.cert_manager,
  ]
}
