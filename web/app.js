(() => {
  const $ = (s) => document.querySelector(s);
  const $$ = (s) => Array.from(document.querySelectorAll(s));

  const state = {
    baseURL: localStorage.getItem("baseURL") || "http://127.0.0.1:8888", // 可配置(默认8888)
    wsURL: "",
    token: localStorage.getItem("token") || "",
    ws: null,
  };

  // 简易 API 客户端
  async function request(path, { method = "GET", headers = {}, body, withAuth = true } = {}) {
    const url = `${state.baseURL}/${path.replace(/^\//, "")}`;
    const h = { ...headers };
    if (withAuth && state.token) h["Authorization"] = `Bearer ${state.token}`;
    // 仅当 body 为字符串时再显式设置 json，否则让浏览器自动决定（避免触发预检）
    if (typeof body === "string" && !h["Content-Type"]) {
      h["Content-Type"] = "application/json";
    }

    const res = await fetch(url, { method, headers: h, body });
    if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
    return res.json().catch(() => ({}));
  }
  // 后端地址保存
  const baseUrlInput = document.querySelector('#base-url');
  const saveBaseBtn = document.querySelector('#save-base-url');
  if (baseUrlInput) baseUrlInput.value = state.baseURL;
  if (saveBaseBtn) {
    saveBaseBtn.onclick = () => {
      const v = baseUrlInput.value.trim();
      if (!v) return alert('请输入有效地址');
      state.baseURL = v;
      localStorage.setItem('baseURL', v);
      alert('已保存');
    };
  }

  function show(pageId) {
    $$(".page").forEach((el) => el.classList.add("hidden"));
    $(`#${pageId}`).classList.remove("hidden");
  }

  // 导航
  $("#nav-login").onclick = () => show("page-login");
  $("#nav-todo").onclick = () => {
    show("page-todo");
    loadTodos();
  };
  $("#nav-approval").onclick = () => {
    show("page-approval");
    loadApprovals();
  };
  $("#nav-chat").onclick = () => show("page-chat");
  $("#nav-upload").onclick = () => show("page-upload");

  // 登录
  // 管理员登录：不保存 token，不带 Authorization
  $("#form-login-admin").addEventListener("submit", async (e) => {
    e.preventDefault();
    const params = new URLSearchParams(new FormData(e.target));
    try {
      const url = `${state.baseURL}/v1/user/login`;
      const res = await fetch(url, { method: "POST", body: params });
      if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
      const json = await res.json();
      // 保存 token（管理员登录也需要）
      state.token = json?.data?.token || json?.token || json?.AccessToken || "";
      if (!state.token) throw new Error("未返回 token");
      localStorage.setItem("token", state.token);
      $("#login-admin-result").textContent = JSON.stringify(json, null, 2);
      alert("管理员登录成功（已保存 token）");
    } catch (err) {
      $("#login-admin-result").textContent = String(err);
      alert("管理员登录失败: " + err.message);
    }
  });

  // 用户登录：保存 token，后续请求携带 Authorization
  $("#form-login-user").addEventListener("submit", async (e) => {
    e.preventDefault();
    // 使用 x-www-form-urlencoded，避免预检 OPTIONS
    const params = new URLSearchParams(new FormData(e.target));
    try {
      const url = `${state.baseURL}/v1/user/login`;
      const res = await fetch(url, { method: "POST", body: params });
      if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
      const json = await res.json();
      state.token = json?.data?.token || json?.token || json?.AccessToken || "";
      if (!state.token) throw new Error("未返回 token");
      localStorage.setItem("token", state.token);
      $("#login-user-result").textContent = JSON.stringify(json, null, 2);
      alert("用户登录成功（已保存 token）");
    } catch (err) {
      $("#login-user-result").textContent = String(err);
      alert("用户登录失败: " + err.message);
    }
  });

  // 待办
  async function loadTodos() {
    const userId = $("#todo-userId").value.trim();
    const query = new URLSearchParams();
    if (userId) query.set("userId", userId);
    const res = await request(`v1/todo/list?${query.toString()}`);
    const list = res?.data || res?.List || [];
    const tbody = $("#todo-table tbody");
    tbody.innerHTML = "";
    (list || []).forEach((item) => {
      const tr = document.createElement("tr");
      tr.innerHTML = `<td>${item.id || item.ID || ""}</td>
        <td>${item.title || ""}</td>
        <td>${item.status ?? item.todoStatus ?? ""}</td>
        <td>
          <button data-id="${item.id || item.ID || ""}" class="btn-finish">完成</button>
          <button data-id="${item.id || item.ID || ""}" class="btn-delete">删除</button>
        </td>`;
      tbody.appendChild(tr);
    });
  }
  $("#todo-refresh").onclick = () => loadTodos().catch(console.error);

  $("#form-todo-create").addEventListener("submit", async (e) => {
    e.preventDefault();
    const form = new FormData(e.target);
    const executeIds = String(form.get("executeIds") || "")
      .split(",")
      .map((s) => s.trim())
      .filter(Boolean);
    const payload = {
      title: form.get("title"),
      deadlineAt: Number(form.get("deadlineAt")) || 0,
      executeIds,
      desc: form.get("desc") || "",
    };
    await request("v1/todo", { method: "POST", body: JSON.stringify(payload) });
    await loadTodos();
    e.target.reset();
  });

  document.addEventListener("click", async (e) => {
    const finishBtn = e.target.closest(".btn-finish");
    const deleteBtn = e.target.closest(".btn-delete");
    if (finishBtn) {
      const todoId = finishBtn.getAttribute("data-id");
      await request("v1/todo/finish", {
        method: "POST",
        body: JSON.stringify({ todoId }),
      });
      await loadTodos();
    }
    if (deleteBtn) {
      const todoId = deleteBtn.getAttribute("data-id");
      await request(`v1/todo/${todoId}`, { method: "DELETE" });
      await loadTodos();
    }
  });

  // 审批
  async function loadApprovals() {
    const userId = $("#approval-userId").value.trim();
    const type = $("#approval-type").value;
    const q = new URLSearchParams();
    if (userId) q.set("userId", userId);
    if (type) q.set("type", type);
    const res = await request(`v1/approval/list?${q.toString()}`);
    const list = res?.data || [];
    const tbody = $("#approval-table tbody");
    tbody.innerHTML = "";
    (list || []).forEach((item) => {
      const tr = document.createElement("tr");
      tr.innerHTML = `<td>${item.id || ""}</td>
        <td>${item.title || ""}</td>
        <td>${item.status ?? ""}</td>
        <td><button data-id="${item.id || ""}" class="btn-approve-pass">通过</button>
            <button data-id="${item.id || ""}" class="btn-approve-reject">拒绝</button></td>`;
      tbody.appendChild(tr);
    });
  }
  $("#approval-refresh").onclick = () => loadApprovals().catch(console.error);

  document.addEventListener("click", async (e) => {
    const passBtn = e.target.closest(".btn-approve-pass");
    const rejectBtn = e.target.closest(".btn-approve-reject");
    if (passBtn || rejectBtn) {
      const approvalId = (passBtn || rejectBtn).getAttribute("data-id");
      const payload = { approvalId, status: passBtn ? 1 : 2, reason: passBtn ? "" : "不同意" };
      await request("v1/approval/dispose", { method: "PUT", body: JSON.stringify(payload) });
      await loadApprovals();
    }
  });

  // 聊天（HTTP）
  $("#form-chat").addEventListener("submit", async (e) => {
    e.preventDefault();
    const text = new FormData(e.target).get("prompts");
    const res = await request("v1/chat", { method: "POST", body: JSON.stringify({ prompts: text }) });
    const log = $("#chat-log");
    log.textContent += `\n> ${text}`;
    log.textContent += `\n< ${JSON.stringify(res?.data ?? res, null, 2)}\n`;
    log.scrollTop = log.scrollHeight;
    e.target.reset();
  });

  // 文件上传
  async function uploadFile(file, withChat = false) {
    const url = `${state.baseURL}/v1/upload/file`;
    const fd = new FormData();
    fd.append("file", file);
    if (withChat) fd.append("chat", "1");
    const headers = {};
    if (state.token) headers["Authorization"] = `Bearer ${state.token}`;
    const res = await fetch(url, { method: "POST", headers, body: fd });
    if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
    return res.json();
  }
  $("#btn-upload").onclick = async () => {
    const f = $("#file-input").files[0];
    if (!f) return alert("请选择文件");
    const res = await uploadFile(f, true);
    alert("上传成功: " + JSON.stringify(res));
  };
  $("#upload-file-btn").onclick = async () => {
    const f = $("#upload-file").files[0];
    if (!f) return alert("请选择文件");
    const res = await uploadFile(f, false);
    $("#upload-result").textContent = JSON.stringify(res, null, 2);
  };

  // WebSocket
  function connectWS() {
    const defaultWs = `ws://127.0.0.1:9000/ws`;
    const url = $("#ws-url").value.trim() || defaultWs;
    state.wsURL = url;
    const proto = state.token ? [state.token] : undefined;
    const ws = (state.ws = new WebSocket(url, proto));
    ws.onopen = () => {
      appendWSLog("[WS] 已连接");
    };
    ws.onmessage = (ev) => {
      appendWSLog(`[WS] 收到: ${ev.data}`);
    };
    ws.onclose = () => {
      appendWSLog("[WS] 已关闭");
      if (state.ws === ws) state.ws = null;
    };
    ws.onerror = (e) => {
      appendWSLog("[WS] 错误");
      console.error(e);
    };
  }
  function appendWSLog(text) {
    const log = $("#chat-log");
    log.textContent += `\n${text}`;
    log.scrollTop = log.scrollHeight;
  }
  $("#ws-connect").onclick = () => {
    if (state.ws) return alert("已连接");
    connectWS();
  };
  $("#ws-disconnect").onclick = () => {
    if (state.ws) state.ws.close();
  };

  // 默认展示
  show("page-login");
})();


