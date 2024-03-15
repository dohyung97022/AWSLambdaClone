from flask import Flask
from handler import handler

app = Flask(__name__)

@app.route('/endpoint/<id>')
def listen(id):
    return handler()

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=443)
