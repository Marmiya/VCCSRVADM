import json
import time
from http.server import BaseHTTPRequestHandler, HTTPServer
import threading

json_data = {
    "message": "Hello, World!",
    "timestamp": time.time()
}

# 更新JSON数据的函数
def update_json_data():
    global json_data
    while True:
        json_data = {
            "message": "Hello, World!",
            "timestamp": time.time()
        }
        time.sleep(30)

class RequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == '/json':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps(json_data).encode())

def run_server():
    server_address = ('', 9527)
    httpd = HTTPServer(server_address, RequestHandler)
    print('Starting server on port 9527...')
    httpd.serve_forever()

# 启动定时器线程
threading.Thread(target=update_json_data, daemon=True).start()

# 启动HTTP服务器
run_server()
