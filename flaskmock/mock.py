from flask import request, Flask, jsonify

app = Flask(__name__)
app.config['JSON_AS_ASCII'] = False


@app.route('/', methods=['POST'])
def post_Data():
    instance_name = request.form['instance_name']
    disk_size = request.form['disk_size']
    recognize_info = {'instance_name': instance_name, 'disk_size': disk_size}
    return jsonify(recognize_info), 200


if __name__ == '__main__':
    app.run(debug=False, host='0.0.0.0', port=8888)

# 安装依赖 pip3 install flask
# 启动服务： python3 mock.py