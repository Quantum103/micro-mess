const API = "http://localhost:8080";

let ws = null;
let currentChatUser = null;

const token = localStorage.getItem("token");
const currentUserId = localStorage.getItem("user_id");

if (!token) {
    alert("Не авторизован");
    window.location.href = "/login.html";
}

/* ================= WS ================= */

function connectWS() {
    ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

    ws.onopen = () => {
        console.log("WS connected");
    };

    ws.onmessage = (event) => {
        const msg = JSON.parse(event.data);

        if (!currentChatUser) return;

        if (msg.from == currentChatUser || msg.to == currentChatUser) {
            renderMessage(msg);
        }
    };

    ws.onclose = () => {
        console.log("WS closed → reconnect");
        setTimeout(connectWS, 2000);
    };
}

/* ================= USERS ================= */

async function loadUsers() {
    const res = await fetch(API + "/api/users", {
        headers: {
            "Authorization": "Bearer " + token
        }
    });

    if (!res.ok) {
        console.error(await res.text());
        return;
    }

    const users = await res.json();

    if (!Array.isArray(users)) return;

    const container = document.getElementById("users");
    container.innerHTML = "";

    users.forEach(u => {
        const div = document.createElement("div");
        div.className = "user";
        div.innerText = u.username || ("User " + u.id);

        div.onclick = () => selectUser(u, div);

        container.appendChild(div);
    });
}

/* ================= SELECT ================= */

async function selectUser(user, element) {
    currentChatUser = String(user.id);

    document.getElementById("header").innerText =
        "Чат с: " + (user.username || user.id);

    document.querySelectorAll(".user")
        .forEach(u => u.classList.remove("active"));

    element.classList.add("active");

        const res = await fetch(
            API + `/api/history?user_id=${currentUserId}&partner_id=${currentChatUser}`,
            {
                headers: {
                    "Authorization": "Bearer " + token
                }
            }
        );

    const messages = await res.json();

    const box = document.getElementById("messages");
    box.innerHTML = "";

    if (Array.isArray(messages)) {
        messages.forEach(renderMessage);
    }
}

/* ================= RENDER ================= */

function renderMessage(msg) {
    const div = document.createElement("div");
    div.classList.add("message");

    const isMine = msg.from_me || msg.from === undefined;

    if (isMine) div.classList.add("me");
    else div.classList.add("other");

    div.innerText = msg.content;

    const box = document.getElementById("messages");
    box.appendChild(div);
    box.scrollTop = box.scrollHeight;
}

/* ================= SEND ================= */

function sendMessage() {
    if (!ws || ws.readyState !== WebSocket.OPEN) return;

    const input = document.getElementById("messageInput");
    const text = input.value.trim();

    if (!text || !currentChatUser) return;

    ws.send(JSON.stringify({
        to: currentChatUser,
        content: text
    }));

    input.value = "";
}

/* ================= INIT ================= */

window.onload = () => {
    connectWS();
    loadUsers();
};