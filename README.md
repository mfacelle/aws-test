# AWS Test Go API

A small Go HTTP API that stores messages in DynamoDB. It is designed to run on an
AWS EC2 instance using an instance profile, so no AWS keys need to be stored on
the server.

## API

- `GET /health` returns a health response without contacting DynamoDB.
- `GET /items` lists all items.
- `POST /items` creates an item from `{"message":"hello"}`.
- `DELETE /items/{id}` removes an item.

## Run locally

Prerequisites: Go 1.26+ and AWS credentials configured through the standard AWS
credential chain (`aws configure`, environment variables, or an AWS profile).

Create the table once:

```sh
aws dynamodb create-table \
	--table-name items \
	--attribute-definitions AttributeName=id,AttributeType=S \
	--key-schema AttributeName=id,KeyType=HASH \
	--billing-mode PAY_PER_REQUEST \
	--region us-east-1
```

Start the API:

```sh
AWS_REGION=us-east-1 TABLE_NAME=items go run .
```

Try it:

```sh
curl http://localhost:8080/health
curl -X POST http://localhost:8080/items \
	-H 'Content-Type: application/json' \
	-d '{"message":"hello from Go"}'
curl http://localhost:8080/items
```

`AWS_REGION`, `TABLE_NAME`, and `PORT` are configurable. They default to
`us-east-1`, `items`, and `8080` respectively.

## Deploy to EC2

1. Create an IAM role for EC2 with a policy limited to the `items` table:

```json
{
	"Version": "2012-10-17",
	"Statement": [{
		"Effect": "Allow",
		"Action": [
			"dynamodb:DeleteItem",
			"dynamodb:PutItem",
			"dynamodb:Scan"
		],
		"Resource": "arn:aws:dynamodb:us-east-1:ACCOUNT_ID:table/items"
	}]
}
```

2. Attach that role to the EC2 instance. Open TCP port 8080 in its security
	 group only when direct external access is needed; placing the API behind an
	 Application Load Balancer is preferable for production.

3. Build on a machine with Go installed and copy the binary and service file:

```sh
GOOS=linux GOARCH=amd64 go build -o aws-test .
scp aws-test aws-test.service ec2-user@EC2_HOST:/tmp/
```

4. On the EC2 instance, install and start it:

```sh
sudo mkdir -p /opt/aws-test
sudo mv /tmp/aws-test /opt/aws-test/aws-test
sudo chmod +x /opt/aws-test/aws-test
sudo mv /tmp/aws-test.service /etc/systemd/system/aws-test.service
sudo systemctl daemon-reload
sudo systemctl enable --now aws-test
sudo systemctl status aws-test
```

Update `AWS_REGION` and `TABLE_NAME` in `aws-test.service` before starting if
your table uses different values. The service uses the EC2 instance role via
the AWS SDK default credential provider chain.
