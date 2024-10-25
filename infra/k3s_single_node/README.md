# Single Node K3s clsuter to Demo installing ENBUILD in AirGap Environment

This terraform deploys the single node K3s cluster on the AWS EC2 instance and deploys a shell script to installs the ENBUILD in the AirGap environment.

Once the instance is deployed , you can login to AWS console and delete the Egress rule from the security group to make the instance AirGap i.e. no internet access.

After that ssh into the instance using the keypair, and run the below command to install the ENBUILD.

```bash
sudo bash /tmp/install_enbuild_haul.sh
```

All the information about the Publuic IP, KeyPair, and Security group from which you need to delete the outbound rule and other details are available in the terraform output.

# Why its not fully automated.
Since we need to deploy the k3s and get the tools like kubectl , hauler , helm etc. to install the ENBUILD, we need to have internet access to download the tools. So, we need to deploy the instance with internet access and then remove the internet access to make it AirGap.