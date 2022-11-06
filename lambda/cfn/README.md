# CloudFormation Templates
Examples for CloudFormation temaplates to setup all necessary components.

## ecr-repository.yml
Template to create a ECR repository for Lambda function image. It's separated from all other source, because it have to be created and an image have to be available before creating the Lambda function.

## lambda-function.yml
This template creates following components.
- Lambda function
- CloudWatch rule to trigger it periodically 
- CloudWatch log group
- SQS queue for weather data events
- All necessary permissions