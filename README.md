# Vim Tricks Slackbot

1. [Technology Used](#technology-used)
1. [Local Setup](#local-setup)
1. [Running the Lambda Locally](#running-the-lamdba-locally)
1. [Running the Lambda Tests](#running-the-tests)
1. [Deployment](#Deployment)
1. [Structure](#Structure)

### Technology Used
1. [Go 1.16](https://golang.org/doc/go1.16)
1. [Docker](https://www.docker.com/)
1. [AWS Lambda](https://docs.aws.amazon.com/lambda/latest/dg/welcome.html)
1. [AWS Elastic Container Registry (ECR)*](https://aws.amazon.com/ecr/?nc2=type_a)
1. [AWS Simple Email Service*](https://aws.amazon.com/getting-started/hands-on/send-an-email/)
1. [AWS Go Lambda Runtime Interface Client](https://docs.aws.amazon.com/lambda/latest/dg/go-image.html#go-image-clients)

* Not Implemented Yet

### Local Setup
1. Clone this repo by running `git clone https://github.com/jm96441n/vim-tricks-slackbot.git` from your terminal
1. From the repo run `make setup` this will install the Lambda Runtime Interface Client to your
machine in order to enable local testing.
1. That's it! With that you should be good to start writing code!

### Running the Lambda Locally
To test the app locally should be pretty straightforward:
1. From the root of the repo run `make run`, this will build and tag the container
 and then start it up in the terminal with the logs running the current window.
 1. From another terminal window (if you are not running the container as a daemon) run
`make curl_test` this will make an empty request to the lamdba to start the lambda execution.
1. To stop the container just run `make stop`.

### Running the Tests
1. From the root of the application run `make test`, this will build the test container and run all of the tests
in the repo.

### Deployment
* NOTE: Deployment is not implemented yet
Deployment is automated through github actions. Once a branch is merged to master the pipeline is
started which builds the container, tags it with the `latest` tag and pushes it to Amazon ECR. Following
the push to ECR the pipeline then deploys the container to the scraper lambda (see [ecr_push.yml](./.github/workflows/ecr_push.yml))

### Structure
The entrypoint for this lambda is in the `main.go`, here we initialize all the dependencies and pass them down to the
respective packages/domains. This lambda triggers whenever an email is sent from [Vim Tricks](https://vimtricks.com/),
and it first pulls down the RSS feed for Vim Tricks and parses it for the latest entry. From there we do some formatting
(converting the html to markdown and then making some small changes to Slack's specific flavor of markdown) and then
post the message to the specified slack channel (set from the configuration.) Future plans for extensibility would be to
have this lambda post to a SNS topic with the formatted message and have separate lambda's post to different slack
channels/groups.
