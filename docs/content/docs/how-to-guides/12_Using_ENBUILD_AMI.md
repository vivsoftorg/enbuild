---
title: "Using ENBUILD AMI"
description: "Using ENBUILD AMI"
summary: "Steps to use ENBUILD AMI to launch ENBUILD instance"
draft: false
menu:
  docs:
    parent: "docs/how-to-guides/"
    identifier: "EnbuildAMI"
weight: 203
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

# Creating an AWS EC2 Instance Using the ENBUILD AMI or Marketplace Item

This guide will walk you through the steps to create an Amazon EC2 instance using the AMI named ENBUILD.

The ENBUILD AMI has ENBUILD pre-installed in it, along with the gitlab.

The Gitlab also has Catalog repositories of BigBang and EKS pre-populated with it.

## Prerequisites

- An AWS account with the necessary permissions to create EC2 instances.
- Basic familiarity with the AWS Management Console.

## Steps

### 1. Log in to the AWS Management Console

1. Go to [AWS Management Console](https://aws.amazon.com/console/).
2. Enter your AWS account credentials to log in.

### 2. Navigate to the EC2 Dashboard

1. In the AWS Management Console, type "EC2" in the search bar and select **EC2** from the dropdown list.
2. This will take you to the EC2 Dashboard.

### 3. Launch an Instance

1. On the EC2 Dashboard, click on the **Launch Instance** button.

### 4. Choose an Amazon Machine Image (AMI)

1. In the **Choose an Amazon Machine Image (AMI)** step, go to the **My AMIs** tab or **AWS Marketplace AMI's** 
2. Use the search bar to find the AMI named **ENBUILD**.
3. Select the **ENBUILD** AMI by clicking the **Select** button next to it.

<picture><img src="/images/how-to-guides/enbuild_ami.png" alt="Screenshot of ENBUILD AMI"></img></picture>

### 5. Choose an Instance Type

1. Select an appropriate instance type based on your requirements (e.g., `t2.micro` for free tier eligible users).
2. Click the **Next: Configure Instance Details** button.

### 6. Configure Instance Details

1. Configure the instance details as needed. The default settings are typically sufficient for most users.
2. Click the **Next: Add Storage** button.

### 7. Add Storage

1. Adjust the storage settings if necessary. The default settings usually suffice.
2. Click the **Next: Add Tags** button.

### 8. Add Tags

1. (Optional) Add tags to your instance to help organize and manage your resources.
2. Click the **Next: Configure Security Group** button.

### 9. Configure Security Group

1. Create a new security group or select an existing one.
   - If creating a new security group, add rules to allow necessary inbound traffic (e.g., SSH for Linux instances, RDP for Windows instances).
2. Click the **Review and Launch** button.

### 10. Review and Launch

1. Review your instance configuration to ensure everything is correct.
2. Click the **Launch** button.

### 11. Select a Key Pair

1. In the **Select an existing key pair or create a new key pair** dialog:
   - Select an existing key pair, or
   - Create a new key pair and download the private key file (`.pem`). Make sure to keep this file safe, as you will need it to connect to your instance.
2. Check the acknowledgment box and click the **Launch Instances** button.

### 12. View Your Instance

1. Click the **View Instances** button to go to the EC2 Dashboard and view your newly created instance.
2. Wait for the instance to enter the `running` state.

### 13. Connect to Your Instance

1. Select your instance in the EC2 Dashboard.
2. Click the **Connect** button and follow the instructions to connect to your instance using SSH (for Linux instances) or RDP (for Windows instances).

Congratulations! You have successfully launched an EC2 instance using the ENBUILD AMI.

## Accessing Services of ENBUILD

The services can be accessed after launching the instance using the AMI at the following URLs with the given public IP (`NODE_PUBLIC_IP`) of the node:

- **Enbuild**: [http://NODE_PUBLIC_IP:30080/](http://$NODE_PUBLIC_IP:30080/)
- **Gitlab**: [http://NODE_PUBLIC_IP:32080/](http://$NODE_PUBLIC_IP:32080/)
- **Hauler**: [http://NODE_PUBLIC_IP:30500/v2/_catalog](http://$NODE_PUBLIC_IP:30500/v2/_catalog)

Default credentials for Gitlab are `root/supersecretpassword`, and the token for the root user is `glpat-RuJrDn4yUuPRJySjJdZh`.



## Additional Resources

- [AWS EC2 Documentation](https://docs.aws.amazon.com/ec2/)
- [AWS AMI Documentation](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html)

If you encounter any issues or have questions, refer to the AWS documentation or seek support from AWS.


