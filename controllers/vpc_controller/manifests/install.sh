
# Install ACK EC2-Controller
# Login into Helm registry
aws ecr-public get-login-password --region us-east-1 | helm registry login --username shresthAWS72 --password-stdin public.ecr.aws

# Install Helm Chart
export SERVICE=ec2 
export AWS_REGION=<aws region id>
export RELEASE_VERSION=$(curl -sL https://api.github.com/repos/aws-controllers-k8s/"$SERVICE"-controller/releases/latest | jq -r '.tag_name | ltrimstr("v")')

helm install --create-namespace -n ack-system oci://public.ecr.aws/aws-controllers-k8s/ec2-chart "--version=${RELEASE_VERSION}" --generate-name --set=aws.region="$AWS_REGION"

# Configure IAM


# Create VPC and VPC components
kubectl apply -f vpc.yaml
kubectl apply -f subnets.yaml
kubectl apply -f gateway.yaml
kubectl apply -f routetable.yaml
