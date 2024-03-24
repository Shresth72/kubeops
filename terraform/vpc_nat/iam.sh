aws configure

# List all users
aws configure list

aws iam list-users
aws iam list-groups-for-user --user-name <user-name>
aws iam list-groups
aws iam list-attached-group-policies --group <group-name>

# Create a new user
aws iam create-user --user-name tf-dev
aws iam create-group --group-name tf-dev-group
aws iam add-user-to-group --user-name tf-dev --group-name tf-dev-group

aws iam create-access-key --user-name tf-dev
aws configure --profile tf-dev

# Delete a user
aws iam delete-access-key --user-name tf-dev --access-key-id <access-key-id>
aws iam delete-user --user-name tf-dev