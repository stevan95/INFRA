#!/bin/sh

account_id=`aws sts get-caller-identity --query "Account" --output text`

read -p 'Enter name of your cluster: ' cluster_name

echo -e "\nCluster Name is: ${CLUSTER_NAME}"
echo -e "\nIf you want to proceed with above information, type \"yes\" or \"no\": "
read value

if [ $value == "yes" ]; then 
    #Create appropriate IAM policy
    curl -O https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/v2.4.7/docs/install/iam_policy.json
    aws iam create-policy \
        --policy-name AWSLoadBalancerControllerIAMPolicy \
        --policy-document file://iam_policy.json

    rm iam_policy.json
    
    #Create AWS role and sa which will be used by controller pod to assume the role
    eksctl create iamserviceaccount \
        --cluster ${cluster_name} \
        --namespace kube-system \
        --name=aws-load-balancer-controller \
        --role-name AmazonEKSLoadBalancerControllerRole \
        --attach-policy-arn=arn:aws:iam::${account_id}:policy/AWSLoadBalancerControllerIAMPolicy \
        --approve

    #Install LoadBalancer using helm
    #Add appropriate helm
    helm repo add eks https://aws.github.io/eks-charts
    helm repo update

    helm install aws-load-balancer-controller eks/aws-load-balancer-controller \
        -n kube-system \
        --set clusterName=${cluster_name} \
        --set serviceAccount.create=false \
        --set serviceAccount.name=aws-load-balancer-controller 

    exit 0
else 
    echo -e "Operation stopped.\n" 
    exit 0
fi