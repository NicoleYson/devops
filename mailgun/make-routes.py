import os
import requests
import yaml
import sys

def importNewRoutes(yaml_file):
    emails = yaml.load(yaml_file) # yaml file with `source: target` format
    for source in emails:
        target = emails[source]
        print(("Creating route for %s -> %s")%(source, target))
        createRoute(source, target)

def createRoute(source, target):
    print("Creating Route...")
    provider = session.post(
        "https://api.mailgun.net/v3/routes",
        auth=(username, api_key),
        data={"priority": 10,
              "description": "TESTING PYTHON ROUTES",
              "expression": "match_recipient('%s')" %source,
              "action": ["forward('%s')" %target, "stop()"]}
    )
    print(provider.status_code)
    print(provider.text)

if __name__ == "__main__":

    username = "api"
    api_key = os.environ.get('MAILGUN_API_KEY')

    if api_key is None:
        print("Please define your MAILGUN_API_KEY as an environment variable.")
        exit()

    session = requests.Session()
    with open(sys.argv[1], 'r') as input_file:
        importNewRoutes(input_file.read())
