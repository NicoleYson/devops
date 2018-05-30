import os
import requests
import sys
import re
import json

def getUrl(logs):
    for contents in logs['items']:
        viewStoredMessages(contents['storage'].get('url'))

def viewStoredMessages(url):
    headers = {"Accept": "message/rfc2822"}
    r = requests.get(url, auth=("api", api_key), headers=headers)
    if r.status_code == 200:
        print(r.json()["body-mime"])
    else:
        print("Oops! Something went wrong: %s" % r.content)

def formatLogs(logs):
     for contents in logs['items']:
        sender = contents['envelope'].get('sender')
        subject = (contents['message'].get('headers').get('subject'))
        print("%s --- %s") %(sender, subject)


def getLogs(email, domain):
    provider = session.get(
        "https://api.mailgun.net/v3/%s/events" % domain,
        auth=("api", api_key),
        params={"recipient": "%s" % email}
    )

    print(provider.status_code)
    return provider.text

if __name__ == "__main__":

    api_key = os.environ.get('MAILGUN_API_KEY')

    if api_key is None:
        print("Please define your MAILGUN_API_KEY as an environment variable.")
        exit()

    session = requests.Session()
    if len(sys.argv) != 2:
        print("Please provide an email")

    email = sys.argv[1]
    regex_domain = re.compile(r"(?<=@)[\w\W]+")
    domain = regex_domain.search(email).group(0)

    logs = json.loads(getLogs(email, domain))
    # getUrl(logs) #use this for stored messages
    formatLogs(logs)

