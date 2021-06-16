from textblob import TextBlob

def handle_event(context, data):
    blob = TextBlob(data.decode())
    return {
        "sentiment": blob.sentiment.polarity,
        "subjectivity": blob.sentiment.subjectivity
    }

def init_context(context):
    context.logger.info("init context from custom image")