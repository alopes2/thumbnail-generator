resource "aws_sqs_queue" "queue" {
  name = "image-events"
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
  policy_id = "${aws_sqs_queue.queue.arn}/SQSDefaultPolicy"

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
      aws_sqs_queue.queue.arn,
    ]

    condition {
      test     = "ArnEquals"
      variable = "aws:SourceArn"

      values = aws_sns_topic.topic.arn
    }
  }
}
