from tensorflow.keras.applications.resnet50 import ResNet50
from tensorflow.keras.preprocessing import image
from tensorflow.keras.applications.resnet50 import preprocess_input, decode_predictions
import numpy as np

# Try it with a config like:

# processor:
#     workers: 2
# source:
#     kind: http

# curl -X POST --data-binary "@Elephant.jpg" http://localhost:8080
def handle_event(context, data):
    img_path = f"/tmp/{context.worker_name}-file.jpg"

    with open(img_path, "wb") as f:
        f.write(data)

    img = image.load_img(img_path, target_size=(224, 224))
    x = image.img_to_array(img)
    x = np.expand_dims(x, axis=0)
    x = preprocess_input(x)

    preds = context.model.predict(x)
    data = decode_predictions(preds, top=3)[0]
    out = []
    for i in data:
        out.append({"label": str(i[1]), "prob": float(i[2])})
    context.logger.info(out)
    return out

def init_context(context):
    model = ResNet50(weights='imagenet')
    setattr(context, "model", model)