package main

import (
	"fmt"
	"log"
	rand "math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
)

func (client *iamClient) getPasswordPolicy() *iam.PasswordPolicy {

	input := &iam.GetAccountPasswordPolicyInput{}
	resp, err := client.svc.GetAccountPasswordPolicy(input)

	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "NoSuchEntity" {
			//IAM account password policy is not there (i.e. using the default policy)
			log.Fatal(iam.ErrCodeNoSuchEntityException)
			return nil
		}
		log.Printf("Error reading IAM account password policy: %s", err)
	}
	return resp.PasswordPolicy
}

// formatPasswordPolicy - Deliberately separated in case someone figures out a nicer way to format this
func (client *iamClient) formatPasswordPolicy() []string {
	policy := client.getPasswordPolicy()
	setPolicy := []string{}

	switch {
	case *policy.RequireLowercaseCharacters:
		setPolicy = append(setPolicy, "lowercase letters")
		fallthrough
	case *policy.RequireUppercaseCharacters:
		setPolicy = append(setPolicy, "uppercase letters")
		fallthrough
	case *policy.RequireNumbers:
		setPolicy = append(setPolicy, "numbers")
		fallthrough
	case *policy.RequireSymbols:
		setPolicy = append(setPolicy, "symbols")
		fallthrough
	case *policy.MinimumPasswordLength > 0:
		minPasswordLength := "at least " + strconv.FormatInt(int64(*policy.MinimumPasswordLength), 10) + " characters"
		setPolicy = append(setPolicy, minPasswordLength)
	}

	setPolicy = append(setPolicy)
	return setPolicy
}

func (client *iamClient) providePostResetInstructions() {
	fmt.Println("After logging in with the provided password, you'll be prompted to create a new one. \nThis password must have: ")
	for _, requirement := range client.formatPasswordPolicy() {
		fmt.Println(" • " + requirement)
	}
}

// generateCompliantPassword - Password Policy: 14 charracters, 1 symbol, number, upper case and lowercase letter
func generateCompliantPassword(strlen int) string {
	var r *rand.Rand
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	var char string

	digits := "0123456789"
	symbols := "~=+%^*/()[]{}/!@#$?|"
	alpha := "abcdefghijklmnopqrstuvwxyz"
	upperAlpha := strings.ToUpper(alpha)

	result := make([]byte, strlen)
	for i := range result {
		switch {
		case i%4 == 0:
			char = digits
		case i%4 == 1:
			char = symbols
		case i%4 == 2:
			char = alpha
		case i%4 == 3:
			char = upperAlpha
		default:
			char = alpha
		}
		result[i] = char[r.Intn(len(char))]
	}
	return string(result)
}

func debugMode() bool {
	if len(os.Args) > 1 && os.Args[1] == "--debug" {
		return true
	}
	return false
}

func (client *iamClient) resetPassword() {
	user := client.prompt()
	changeRequired := !debugMode() // --debug is specified, the user will not be forced to reset their password again upon login
	password := generateCompliantPassword(20)

	input := &iam.UpdateLoginProfileInput{
		Password:              aws.String(password),
		UserName:              aws.String(user),
		PasswordResetRequired: aws.Bool(changeRequired),
	}
	fmt.Println("--- Copy paste this to the user ---")
	fmt.Println("Your new password: " + password)
	client.providePostResetInstructions()

	_, err := client.svc.UpdateLoginProfile(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeEntityTemporarilyUnmodifiableException:
				log.Fatal(iam.ErrCodeEntityTemporarilyUnmodifiableException, aerr.Error())
			case iam.ErrCodeNoSuchEntityException:
				log.Fatal(iam.ErrCodeNoSuchEntityException, aerr.Error())
			case iam.ErrCodePasswordPolicyViolationException:
				log.Fatal(iam.ErrCodePasswordPolicyViolationException, aerr.Error())
			case iam.ErrCodeLimitExceededException:
				log.Fatal(iam.ErrCodeLimitExceededException, aerr.Error())
			case iam.ErrCodeServiceFailureException:
				log.Fatal(iam.ErrCodeServiceFailureException, aerr.Error())
			default:
				log.Fatal(aerr.Error())
			}
		} else {
			log.Fatal(err.Error())
		}
		return
	}
}
