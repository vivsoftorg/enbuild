locals {
  name = lower(var.cluster_name)
}

data "aws_vpc" "default" {
  default = true
}

data "aws_ami" "ubuntu" {
    most_recent = true

    filter {
        name   = "name"
        values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
    }
    owners = ["099720109477"]  # Canonical's AWS account ID
}


data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

resource "tls_private_key" "ssh" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "local_file" "pem" {
  filename        = "/tmp/${local.name}.pem"
  content         = tls_private_key.ssh.private_key_pem
  file_permission = "0600"
}

resource "aws_key_pair" "admin" {
  key_name   = local.name
  public_key = tls_private_key.ssh.public_key_openssh
}

data "aws_subnet" "public1" {
  for_each = toset(data.aws_subnets.default.ids)
  id       = each.value
}

// data "cloudinit_config" "userData" {
//   part {
//     content      = file("install_rke2.sh")
//     content_type = "text/x-shellscript"
//   }
// }

resource "aws_security_group" "k3s" {
  name        = "${local.name}-sg"
  description = "Allow all inbound traffic"
  vpc_id      = data.aws_vpc.default.id
}

resource "aws_security_group_rule" "ssh_access" {
  from_port         = 22
  to_port           = 22
  protocol          = "tcp"
  security_group_id = aws_security_group.k3s.id
  type              = "ingress"
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "same_security_group" {
  from_port                = 0
  to_port                  = 0
  protocol                 = "-1"
  security_group_id        = aws_security_group.k3s.id
  type                     = "ingress"
  source_security_group_id = aws_security_group.k3s.id
}

resource "aws_security_group_rule" "nodeport_public" {
  from_port         = 30000
  to_port           = 32767
  protocol          = "tcp"
  security_group_id = aws_security_group.k3s.id
  type              = "ingress"
  cidr_blocks       = ["0.0.0.0/0"]
}


resource "aws_security_group_rule" "all_outbound" {
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  security_group_id = aws_security_group.k3s.id
  type              = "egress"
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_instance" "k3s" {
  ami                         = data.aws_ami.ubuntu.id
  associate_public_ip_address = true
  instance_type               = var.instance_type
  key_name                    = aws_key_pair.admin.key_name
  subnet_id                   = tolist(data.aws_subnets.default.ids)[0]
//   user_data                   = data.cloudinit_config.userData.rendered

  tags = {
    Schedule = var.schedule
    Name     = local.name
  }

  vpc_security_group_ids = [aws_security_group.k3s.id]

  root_block_device {
    volume_size           = 100
    volume_type           = "gp2"
    delete_on_termination = true
  }

  lifecycle {
    ignore_changes = [ebs_block_device]
  }

  provisioner "file" {
    source      = "install_enbuild_haul.sh"
    destination = "/tmp/install_enbuild_haul.sh"

    connection {
      type        = "ssh"
      user        = "ubuntu"
      private_key = file("/tmp/${local.name}.pem")
      host        = self.public_ip
    }
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /tmp/install_enbuild_haul.sh",
      "sudo /tmp/install_enbuild_haul.sh"
    ]

    connection {
      type        = "ssh"
      user        = "ubuntu"
      private_key = file("/tmp/${local.name}.pem")
      host        = self.public_ip
    }
  }
}

output "instance-ip" {
  value = aws_instance.k3s.public_ip
}

output "instance-ssh-key" {
  value = "/tmp/${local.name}.pem"
}

output "instance-login-user" {
  value = "ubuntu"
}

output "node_joining_token" {
  value = "/var/lib/rancher/k3s/server/node-token"
}


output "kubeconfig_file_for_remote_access" {
  value = "/tmp/kube_config.yaml"
}

