# AWS Management Console URL Generator
This Go application generates a sign-in URL for accessing the AWS Management Console. It uses AWS Security Token Service
(STS) to assume a role and then uses the getSigninToken API to create a temporary session that can be used to log in to
the AWS Console.

It is based on the [Enable custom identity broker access to the AWS console](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_enable-console-custom-url.html)
documentation from AWS.

# Prerequisites
* Go 1.18+: Ensure you have Go installed on your system. If not, you can follow the instructions to install Go.
* AWS SDK for Go v2: The application uses the AWS SDK for Go v2 to interact with AWS services.
* Working AWS CLI Profile: The application relies on AWS profiles. If your AWS profile uses AWS SSO (Single Sign-On),
  make sure you have already authenticated using aws sso login before running the program.
* The profile can assume the given role: The profile you use must have the necessary permissions to assume the role you
  specify.

# Setup
Clone the repository or copy the provided Go code into a directory of your choice.
Run `go mod init <module-name>` and `go mod tidy` to download dependencies if required.

# Usage
To run the application, use the following command syntax:
```bash
go run main.go --profile <aws-profile> --role-arn <role-arn> [--session-duration <duration>]
```

Parameters:
* --profile <aws-profile> (Required): The AWS CLI profile to use (e.g., my-sso-profile).
* --role-arn <role-arn> (Required): The ARN of the IAM role you wish to assume (e.g., arn:aws:iam::123456789012:role/my-role).
* --session-duration <duration> (Optional): The duration (in seconds) for which the federated session will remain valid.
  The default session duration is managed by the role. If specified, this will be added to the URL request when
  generating the federated sign-in URL.

> **_NOTE:_**: if you are using role-chaining (your role assumes another role) session duration cannot be defined. Also,
> you should not define longer session duration than the final role allows (max_session_duration property of role).  

## Example Usage
Basic Usage (no session duration specified):

```bash
go run main.go --profile my-aws-profile --role-arn arn:aws:iam::123456789012:role/my-role
```
This will generate a federated sign-in URL using the temporary credentials obtained from assuming the specified role.