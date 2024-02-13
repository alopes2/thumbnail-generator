resource "aws_s3_bucket" "my-app-images" {
  bucket = "my-super-app-images"
}

resource "aws_s3_bucket_notification" "images_put_notification" {
  bucket = aws_s3_bucket.my-app-images.id

  topic {
    topic_arn = aws_sns_topic.topic.arn
    events    = ["s3:ObjectCreated:*"]
  }
}
