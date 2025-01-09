# AWS Management Console URL Generator
When you have AWS cli access (IAM Access Key and Secret Key), an IAM role in an AWS account and you want to access the
AWS Management Console through this role, you can use this tool to generate a URL that will allow you to sign in to the
Console and have the permissions of the role.

The tool generates a federated sign-in URL using the AWS Security Token Service (STS) to assume a role and then uses the
getSigninToken API to create a temporary session that can be used to log using a browser.

It is based on the [Enable custom identity broker access to the AWS console](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_enable-console-custom-url.html)
documentation from AWS.

# Prerequisites
* **Go 1.18+:** Ensure you have Go installed on your system. If not, you can follow the instructions to install Go.
* **AWS SDK for Go v2:** The application uses the AWS SDK for Go v2 to interact with AWS services.
* **Working AWS CLI Profile:** The application relies on AWS profiles. If your AWS profile uses AWS SSO (Single Sign-On),
  make sure you have already authenticated using aws sso login before running the program.
* **The profile can assume the given role:** The profile you use must have the necessary permissions to assume the role you
  specify.

# Setup
## Binary
Download the latest binary from the [releases page](https://github.com/Cajga/aws-console-url/releases)

## From source
```bash
git clone git@github.com:Cajga/aws-console-url.git
cd aws-console-url
GOOS=linux GOARCH=amd64 go build -o aws-console-url -ldflags '-extldflags "-static"' main.go
chmod +x aws-console-url
```

# Usage
To run the application, use the following command syntax:
```bash
./aws-console-url --profile <aws-profile> --role-arn <role-arn> [--session-duration <duration_in_seconds>]
```

Example:
```bash
./aws-console-url --profile my-sso-profile --role-arn arn:aws:iam::123456789012:role/my-role
```

Parameters:
* **--profile <aws-profile> (Required):** The AWS CLI profile to use (e.g., my-sso-profile).
* **--role-arn <role-arn> (Required):** The ARN of the IAM role you wish to assume (e.g., arn:aws:iam::123456789012:role/my-role).
* **--session-duration <duration> (Optional):** The duration (in seconds) for which the federated session will remain valid.
  The default session duration is managed by the role. If specified, this will be added to the URL request when
  generating the federated sign-in URL.

> **_NOTE:_**: if you are using role-chaining (you are using a role to assume another role) session duration cannot be
> defined. Also, you should not define longer session duration than the final role allows (max_session_duration property
> of role which is 1h per default and can be modified to a maximum value of 12h).