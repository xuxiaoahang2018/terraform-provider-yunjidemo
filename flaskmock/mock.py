import json
import uuid
from flask import request, Flask, jsonify, abort

app = Flask(__name__)
app.config['JSON_AS_ASCII'] = False
recognize_info = {}

@app.route('/create', methods=['POST', 'GET'])
def post_Data():
    # 获取传入的参数
    get_Data=request.get_data()
    # 传入的参数为bytes类型，需要转化成json
    get_Data=json.loads(get_Data)
    recognize_info["instance_name"] = get_Data.get('instance_name')
    recognize_info["disk_size"] = get_Data.get('disk_size')
    recognize_info["networks"] = get_Data.get('networks')
    recognize_info["memory"] = get_Data.get('memory')
    recognize_info["config_json"] = get_Data.get('config_json')
    recognize_info["set_demo"] = get_Data.get('set_demo')
    recognize_info["uuid"] = uuid.uuid1()
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
    get_Data=request.get_data()
    # 传入的参数为bytes类型，需要转化成json
    get_Data=json.loads(get_Data)
    if get_Data.get('instance_name'):
        recognize_info["instance_name"] = get_Data.get('instance_name')
    if get_Data.get('disk_size'):
        recognize_info["disk_size"] = get_Data.get('disk_size')
    if get_Data.get('networks'):
        recognize_info["networks"] = get_Data.get('networks')
    if get_Data.get('memory'):
        recognize_info["memory"] = get_Data.get('memory')
    if get_Data.get('config_json'):
        recognize_info["config_json"] = get_Data.get('config_json')
    if get_Data.get("set_demo"):
        recognize_info["set_demo"] = get_Data.get('set_demo')
    return jsonify(recognize_info), 200

@app.route('/delete', methods=['DELETE'])
def delete_Data():
    if request.args.get("id") == "weiyi_demo_id":
        recognize_info = {}
        return jsonify(recognize_info), 200
    abort(400)

@app.route("/data_source", methods=["GET"])
def data_source():
    if request.args.get("name") == "ecs":
        rep = {"name": "ecs",
               "id" : "ecs_id",
               }
        return jsonify(rep), 200
    return jsonify({}), 200



if __name__ == '__main__':
    app.run(debug=False, host='0.0.0.0', port=8888)
    

# 安装依赖 pip3 install flask
# 启动服务： python3 mock.py