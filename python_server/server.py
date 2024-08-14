from flask import Flask, request, jsonify
import os

app = Flask(__name__)

UPLOAD_FOLDER = '../uploads'

@app.route('/upload', methods=['POST'])
def upload_file():
    # 检查目录是否存在，如果不存在则创建它
    if not os.path.exists(UPLOAD_FOLDER):
        os.makedirs(UPLOAD_FOLDER)
    
    file = request.files['file']
    if file:
        # 使用 os.path.join 来构建文件路径
        file_path = os.path.join(UPLOAD_FOLDER, file.filename)
        print("Saving file to:", file_path)
        file.save(file_path)
        return jsonify({"url": f"/uploads/{file.filename}"}), 200
    return jsonify({"error": "No file uploaded"}), 400

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5001)