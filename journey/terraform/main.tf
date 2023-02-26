terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region  = "ap-southeast-1"
  profile = "journey_user"
}

resource "aws_s3_bucket" "website_bucket" {
  bucket = "journey.bustanil.com"

  tags = {
    purpose = "miniproject"
  }
}

resource "aws_s3_bucket_website_configuration" "website_config" {
  bucket = aws_s3_bucket.website_bucket.bucket
  index_document {
    suffix = "index.html"
  }
}

resource "aws_s3_bucket_acl" "website_bucket_acl" {
  bucket = aws_s3_bucket.website_bucket.bucket
  acl    = "public-read"
}

resource "aws_s3_bucket_website_configuration" "journey_web_config" {
  bucket = aws_s3_bucket.website_bucket.id

  index_document {
    suffix = "index.html"
  }

}

data "aws_iam_policy_document" "allow_object_read" {
  statement {
    principals {
      type = "AWS"
      identifiers = ["*"]
    }

    actions   = ["s3:GetObject"]
    resources = ["${aws_s3_bucket.website_bucket.arn}/*"]
  }
}

resource "aws_route53_record" "journey" {
  zone_id = "Z08238352NOVWWT6JWSWJ"
  name    = "journey"
  type    = "A"

  alias {
    name                   =  aws_s3_bucket_website_configuration.website_config.website_domain
    zone_id                = aws_s3_bucket.website_bucket.hosted_zone_id
    evaluate_target_health = false
  }
}

