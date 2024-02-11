resource "helm_release" "nginx_ingress" {
  repository            = "https://kubernetes.github.io/ingress-nginx"
  name                  = "ingress-nginx"
  chart                 = "ingress-nginx"
  version               = "4.9.1"

  namespace = "kube-system"
  
  set {
    name  = "controller.service.type"
    value = "LoadBalancer"
  }

  set {
    name  = "controller.ingressClassResource.name"
    value = "nginx"
  }

  set {
    name  = "controller.watchIngressWithoutClass"
    value = "true"
  }

  set {
    name  = "controller.extraArgs.ingress-class"
    value = "nginx"
  }

  set {
    name  = "controller.service.annotations.service\\.beta\\.kubernetes\\.io/aws-load-balancer-type"
    value = "nlb"
  }

  set {
    name  = "controller.service.annotations.service\\.beta\\.kubernetes\\.io/aws-load-balancer-cross-zone-load-balancing-enabled"
    value = "true"
  }
}


