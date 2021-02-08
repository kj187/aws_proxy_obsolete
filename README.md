# AWS Proxy

```
This project is still work in progress !!!
```

## Problem

In a modern infrastructure architecture services like RDS, ElastiCache, ElasticSearch, MQ etc. should be live in a protected private network which is not accessable from the public. Some of these services are providing web interfaces which are quite helpful to debug problems. But how is it possible to reach them in the browser? And how is it possible to have a look in a database or any other AWS service? 

One way to solve this kind of problem is to open a **proxy** to these services. 
And on this way the **AWS Proxy** as a small CLI utility comes into play to support you with that.

![AWS Proxy](doc/aws_proxy.png "AWS Proxy")

## Advantages

- You dont need to handle with proxy configuration
- You dont need to know endpoint data (AWS Proxy is like a wizard)

## Requirements

### Bastion server 

To open a proxy to your private network services, we need something like an bridge between your public and your private network. This kind of bridge is called a **jump host** or a **bastion server**. This is a server which lives in your public network and has access to your private network. For security reasons you should only allow SSH as ingress. Also, you need to have access from you local machine to the bastion server. The AWS Proxy utility assumes that you have access to the bastion with a valid user and a valid SSH key (with or without a passphrase). 

### AWS CLI access

TODO

## Installation

TODO

## Supported AWS services

- RDS

### Planned services

- Elasticsearch (API and Kibana)
- Elasticache (Redis)
- MQ Service (RabbitMQ)

## Usage 

### Configuration

TODO 

`~/.aws_proxy.yaml
```
aws_proxy config
```

#### Bastion

If your bastion server has a static IP and you know that, you could add this information to the configuration file 
`~/.aws_proxy.yaml`

```
bastion:
    public_ip: XX.XX.XX.XX
```

If you dont know that IP or your bastion server does not have a static IP, AWS Proxy can find that with a tag query. 
The default tag query is that your Bastion has a tag with the key `Name` and the value `bastion`. You could change that filter in the configuration file. You could also add multiple filter if your AWS account provides multiple different bastion instances.

``` 
bastion:
  tag_filter:
  - name: tag:Name
    value: bastion
```

The default SSH **user** is `ec2-user`, the default **port** is `22` and the default **identity file** is `$HOME/.ssh/id_rsa`. 
You could change all these default in the configuration file `~/.aws_proxy.yaml`. 

```
bastion:
  ssh:
    user: ec2-user
    identity_file: ~/.ssh/rd_rsa
    port: 22
```

### Help

TODO

```
aws_proxy --help
aws_proxy -p
aws_proxy -v
```

### Amazon RDS

TODO 

```
aws_proxy proxy:rds 
aws_proxy proxy:rds -n mydbname
aws_proxy proxy:rds -n mydbname -v 
```