import boto3
from botocore.exceptions import ClientError

def excludeUsers(user):
    if excluded_substring not in user and user != your_username:
        return user


def getUserList(users):
    non_service_usernames = []
    for response in users:
        included_users += list(map(lambda username: excludeUsers(username['UserName']), response['Users']))
    included_users_sanitized = list(filter(lambda unsanitized: unsanitized is not None, included_users)) 

    # included_users_sanitized = []
    # for username in included_users:
    #     if excludeUsers(username['UserName']):
    #         included_users_sanitized.append(username)

def deactivateAccessKey(users, keys):
    for user in users:
        for response in keys.paginate(Username=user):
            map(lambda key: iam.update_access_key(
                AccessKeyId=key['AccessKeyId'], Status='Inactive', UserName=user), response['AccessKeyMetadata'])
            print("%s - Deactivating Access Key" %(user))

def forcePasswordReset(users):
    for user in users:
        try:
            print("%s - Forcing Password Reset." %(user))
            iam.update_login_profile(
                UserName=user,
                PasswordResetRequired=True
            )
        except ClientError as e:
            if e.response['Error']['Code'] == 'NoSuchEntity' :
                print("%s - Skipping. Does not have a login profile" %(user))
            else:
                print("Unexpected error: %s") %e

if __name__ == "__main__" :

    iam = boto3.client('iam')
    sts = boto3.client('sts')

    # See who is running the script
    caller_id = sts.get_caller_identity()
    your_username = caller_id['Arn'].split("/")[-1] # Get last portion of ARN string, after the slash

    # Exclude a particular substring of usernames and the person running script
    excluded_substring = raw_input("Enter the substring of usernames you'd like to exclude (e.g. anything prefaced with `svc`) ")
    included_users = getUserList(iam.get_paginator('list_users').paginate())

    deactivateAccessKey(included_users, iam.get_paginator('list_access_keys'))
    forcePasswordReset(included_users)
