前端使用说明（默认端口已适配后端配置）

- 直接打开 index.html 即可使用（推荐同源部署）。默认 API 为 http://127.0.0.1:8888。
- 顶部可设置后端地址，保存到 localStorage。
- 登录成功后会把 JWT 保存到 localStorage.token。
- WebSocket：在“聊天”页填 ws://127.0.0.1:9000/ws（可改），使用 JWT 作为子协议。

接口映射
- POST /v1/user/login
- GET  /v1/todo/list?userId=
- POST /v1/todo
- POST /v1/todo/finish
- DELETE /v1/todo/:id
- GET  /v1/approval/list?userId=&type=
- PUT  /v1/approval/dispose
- POST /v1/chat
- POST /v1/upload/file

同源部署建议（Nginx）
- location /v1/ { proxy_pass http://127.0.0.1:8080; }
- location /ws  { proxy_pass http://127.0.0.1:8081; proxy_set_header Upgrade $http_upgrade; proxy_set_header Connection "upgrade"; }

