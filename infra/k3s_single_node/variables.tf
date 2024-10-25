variable "aws_region" {
  description = "AWS region for all resources."
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment Tag"
  type        = string
  default     = "prod"
}

variable "owner" {
  description = "Owner Tag"
  type        = string
  default     = "Enbuild"
}

variable "cluster_name" {
  description = "Cluster Name"
  type        = string
  default     = "Enbuild-in-k3s"
}

variable "instance_type" {
  description = "instance_type"
  type        = string
  default     = "t3.large"
}

variable "schedule" {
  description = "instance_schedule"
  type        = string
  default     = "on=(M,01);off=(F,23);tz=ist"
}   