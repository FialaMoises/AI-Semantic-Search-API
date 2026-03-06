from flask import Flask, request, jsonify
from sentence_transformers import SentenceTransformer
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = Flask(__name__)

MODEL_NAME = "sentence-transformers/all-MiniLM-L6-v2"
logger.info(f"Loading model: {MODEL_NAME}")
model = SentenceTransformer(MODEL_NAME)
logger.info("Model loaded successfully")

@app.route('/health', methods=['GET'])
def health():
    """Health check endpoint"""
    return jsonify({
        "status": "healthy",
        "model": MODEL_NAME,
        "dimension": model.get_sentence_embedding_dimension()
    })

@app.route('/embed', methods=['POST'])
def embed():
    """Generate embeddings for input texts"""
    try:
        data = request.get_json()

        if not data or 'texts' not in data:
            return jsonify({"error": "Missing 'texts' field in request body"}), 400

        texts = data['texts']

        if not isinstance(texts, list):
            return jsonify({"error": "'texts' must be a list"}), 400

        if len(texts) == 0:
            return jsonify({"error": "'texts' list cannot be empty"}), 400

        logger.info(f"Generating embeddings for {len(texts)} texts")
        embeddings = model.encode(texts, convert_to_numpy=True)

        embeddings_list = embeddings.tolist()

        return jsonify({
            "embeddings": embeddings_list,
            "dimension": model.get_sentence_embedding_dimension()
        })

    except Exception as e:
        logger.error(f"Error generating embeddings: {str(e)}")
        return jsonify({"error": str(e)}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8001, debug=False)
