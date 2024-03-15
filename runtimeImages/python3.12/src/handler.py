from flask import request, jsonify

def handler():
    return jsonify(message="Hello world!", params=request.args), 200
