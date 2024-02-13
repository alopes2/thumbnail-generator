resource "aws_sqs_queue" "queue" {
  name   = "image-events"
  policy = data.aws_iam_policy_document.sqs-queue-policy.json
}

resource "aws_sns_topic" "topic" {
  name = "image-events"
}

resource "aws_sns_topic_subscription" "topic_subscription" {
  topic_arn            = aws_sns_topic.topic.arn
  protocol             = "sqs"
  endpoint             = aws_sqs_queue.queue.arn
  raw_message_delivery = true
}

data "aws_iam_policy_document" "sqs-queue-policy" {
  policy_id = "arn:aws:sqs:${var.region}:${var.account_id}:movie-updates-queue/SQSDefaultPolicy"

  statement {
    sid    = "image-events-allow-send-messages"
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["*"]
    }

    actions = [
      "SQS:SendMessage",
    ]

    resources = [
      "arn:aws:sqs:${var.region}:${var.account_id}:movie-updates-queue",
    ]

    condition {
      test     = "ArnEquals"
      variable = "aws:SourceArn"

      values = [
        aws_sns_topic.topic.arn
      ]
    }
  }
}
