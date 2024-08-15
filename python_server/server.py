from flask import Flask, request, jsonify, url_for, send_from_directory, make_response
import os

app = Flask(__name__)

UPLOAD_FOLDER = '../uploads'
app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER

@app.route('/upload', methods=['POST'])
def upload_file():
    # 检查目录是否存在，如果不存在则创建它
    if not os.path.exists(app.config['UPLOAD_FOLDER']):
        os.makedirs(app.config['UPLOAD_FOLDER'])
    
    file = request.files['file']
    if file:
        # 使用 os.path.join 来构建文件路径
        file_path = os.path.join(app.config['UPLOAD_FOLDER'], file.filename)
        print("Saving file to:", file_path)
        file.save(file_path)
        # 返回完整的文件路径
        file_url = url_for('uploaded_file', filename=file.filename, _external=True)
        return jsonify({"url": file_url}), 200
    
    return jsonify({"error": "No file uploaded"}), 400

@app.route('/uploads/<filename>')
def uploaded_file(filename):
    # 获取文件的扩展名，设置合适的 Content-Type
    file_ext = os.path.splitext(filename)[1].lower()
    mime_type = 'application/octet-stream'
    if file_ext == '.mp3':
        mime_type = 'audio/mpeg'
    elif file_ext in ['.m4a']:
        mime_type = 'audio/mp4'

    response = make_response(send_from_directory(app.config['UPLOAD_FOLDER'], filename))
    response.headers['Content-Type'] = mime_type
    return response

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5001)