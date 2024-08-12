from flask import Flask, request, jsonify

app = Flask(__name__)

@app.route('/upload', methods=['POST'])
def upload_file():
    file = request.files['file']
    if file:
        file.save(f'./uploads/{file.filename}')
        return jsonify({"url": f"/uploads/{file.filename}"}), 200
    return jsonify({"error": "No file uploaded"}), 400

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)