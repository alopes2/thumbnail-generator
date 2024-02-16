resource "aws_sns_topic" "topic" {
  name   = "image-events"
  policy = data.aws_iam_policy_document.sns-topic-policy.json
}

resource "aws_sns_topic_subscription" "topic_subscription" {
  topic_arn            = aws_sns_topic.topic.arn
  protocol             = "lambda"
  endpoint             = aws_lambda_function.lambda.arn
  raw_message_delivery = true
}

resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowExecutionFromSNS"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda.arn
  principal     = "sns.amazonaws.com"
  source_arn    = aws_sns_topic.topic.arn
}

data "aws_iam_policy_document" "sns-topic-policy" {
  policy_id = "arn:aws:sns:${var.region}:${var.account_id}:image-events/SNSS3NotificationPolicy"

  statement {
    sid    = "s3-allow-send-messages"
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["s3.amazonaws.com"]
    }

    actions = [
      "SNS:Publish",
    ]

    resources = [
      "arn:aws:sns:${var.region}:${var.account_id}:image-events",
    ]

    condition {
      test     = "ArnEquals"
      variable = "aws:SourceArn"

      values = [
        aws_s3_bucket.my-app-images.arn
      ]
    }
  }
}
