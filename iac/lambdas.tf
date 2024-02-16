resource "aws_iam_role" "iam_for_lambda" {
  name               = "thumbnail-generator-lambda-role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
  inline_policy {
    name   = "DefaultPolicy"
    policy = data.aws_iam_policy_document.lambda_role_policies.json
  }
}

resource "aws_lambda_function" "lambda" {
  filename      = data.archive_file.lambda.output_path
  function_name = "thumbnail-generator"
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "main"
  runtime       = "go1.x"
}

data "archive_file" "lambda" {
  type        = "zip"
  source_file = "./lambda_init_code/main"
  output_path = "thumbnail_generator_lambda_function_payload.zip"
}

data "aws_iam_policy_document" "assume_role" {

  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

data "aws_iam_policy_document" "lambda_role_policies" {
  statement {
    effect = "Allow"

    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]

    resources = ["arn:aws:logs:*:*:*"]
  }

  statement {
    effect = "Allow"

    actions = [
      "s3:GetObject",
    ]

    resources = [
      format("%s/%s*", aws_s3_bucket.my-app-images.arn, aws_s3_object.images_folder.key)
    ]
  }

  statement {
    effect = "Allow"

    actions = [
      "s3:PutObject",
    ]

    resources = [
      format("%s/%s*", aws_s3_bucket.my-app-images.arn, aws_s3_object.thumbnails_folder.key)
    ]
  }
}
