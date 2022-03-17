from flask import request, Flask, jsonify, abort

app = Flask(__name__)
app.config['JSON_AS_ASCII'] = False
recognize_info = {}

@app.route('/', methods=['POST', 'GET'])
def post_Data():
    recognize_info["instance_name"] = request.form['instance_name']
    recognize_info["disk_size"] = request.form['disk_size']
    return jsonify(recognize_info), {"create_success": 200}

@app.route('/get', methods=['GET'])
def get_Data():
    if request.args.get("id") == "weiyi_demo_id":
        return jsonify(recognize_info), 200
    abort(400)

@app.route('/update', methods=['PUT'])
def update_Data():
    if request.args.get("id") != "weiyi_demo_id":
        abort(400)
    if request.form['instance_name']:
        recognize_info["instance_name"] = request.form['instance_name']
    if request.form['disk_size']:
        recognize_info["disk_size"] = request.form['disk_size']
    return jsonify(recognize_info), 200

@app.route('/delete', methods=['DELETE'])
def delete_Data():
    if request.args.get("id") == "weiyi_demo_id":
        recognize_info = {}
        return jsonify(recognize_info), 200
    abort(400)

if __name__ == '__main__':
    app.run(debug=False, host='0.0.0.0', port=8888)
    

# 安装依赖 pip3 install flask
# 启动服务： python3 mock.py