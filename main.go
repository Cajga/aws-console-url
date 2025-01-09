package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func main() {
	// Define command-line flags for profile and role ARN.
	profile := flag.String("profile", "", "AWS CLI profile name to use to assume the role (required)")
	roleArn := flag.String("role-arn", "", "ARN of the AWS role to assume (required)")
	sessionDuration := flag.Int("session-duration", 0, "Session duration in seconds (optional, default is role configuration)")
	flag.Parse()

	// Validate required flags.
	if *profile == "" || *roleArn == "" {
		flag.Usage()
		log.Fatalf("Both --profile and --role-arn are required.")
	}

	// Step 1: Load the AWS SSO profile configuration.
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(*profile),
		config.WithRegion("us-west-2"), // Specify the region explicitly as it does not matter for this operation.
	)
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	// Step 2: Create an STS client.
	stsClient := sts.NewFromConfig(cfg)

	// Step 3: Assume the role using the SSO credentials.
	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(*roleArn),
		RoleSessionName: aws.String("AssumeRoleSession"),
	}

	result, err := stsClient.AssumeRole(context.TODO(), input)
	if err != nil {
		log.Fatalf("Failed to assume role: %v", err)
	}

	// Step 4: Format resulting temporary credentials into JSON.
	urlCredentials := map[string]string{
		"sessionId":    *result.Credentials.AccessKeyId,
		"sessionKey":   *result.Credentials.SecretAccessKey,
		"sessionToken": *result.Credentials.SessionToken,
	}

	jsonCredentials, err := json.Marshal(urlCredentials)
	if err != nil {
		log.Fatalf("Failed to marshal credentials: %v", err)
	}

	// Construct request URL to AWS federation endpoint.
	requestURL := "https://signin.aws.amazon.com/federation?Action=getSigninToken"

	// Append session duration as a URL parameter if specified.
	if *sessionDuration > 0 {
		requestURL += fmt.Sprintf("&SessionDuration=%d", *sessionDuration)
	}

	// Append the session JSON document as a parameter.
	requestURL += "&Session=" + url.QueryEscape(string(jsonCredentials))

	resp, err := http.Get(requestURL)
	if err != nil {
		log.Fatalf("Failed to get sign-in token: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is not 200.
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Could not get sign-in token. Maybe you are trying to define a session duration while you are using a role-chaining?! HTTP status code: %d", resp.StatusCode)
	}

	var tokenResponse map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		log.Fatalf("Failed to decode token response: %v", err)
	}

	signinToken, ok := tokenResponse["SigninToken"]
	if !ok {
		log.Fatalf("SigninToken not found in response")
	}

	// Step 6: Create URL where users can use the sign-in token to sign in to the console.
	finalURL := fmt.Sprintf(
		"https://signin.aws.amazon.com/federation?Action=login&Issuer=Example.org&Destination=%s&SigninToken=%s",
		url.QueryEscape("https://console.aws.amazon.com/"),
		url.QueryEscape(signinToken),
	)

	// Send final URL to stdout.
	fmt.Println(finalURL)
}
