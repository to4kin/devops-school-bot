terraform {
  required_version = ">= 0.12.2"

  backend "s3" {
    region         = "eu-central-1"
    bucket         = "tf-school-devops-school-bot"
    key            = "terraform.tfstate"
    dynamodb_table = "tf-school-devops-school-bot-lock"
    profile        = ""
    role_arn       = ""
    encrypt        = "true"
  }
}
