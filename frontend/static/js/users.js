// === UI эффекты ===
function createSnowflakes() {
    const snowContainer = document.getElementById('snow');
    if (!snowContainer) return;

    const snowflakeCount = 40;

    for (let i = 0; i < snowflakeCount; i++) {
        const snowflake = document.createElement('div');
        snowflake.classList.add('snowflake');

        const size = Math.random() * 4 + 2;
        snowflake.style.width = `${size}px`;
        snowflake.style.height = `${size}px`;
        snowflake.style.left = `${Math.random() * 100}%`;
        snowflake.style.top = `${Math.random() * 100}%`;

        const duration = Math.random() * 12 + 6;
        snowflake.style.animationDuration = `${duration}s`;

        const delay = Math.random() * 5;
        snowflake.style.animationDelay = `${delay}s`;

        snowContainer.appendChild(snowflake);
    }
}

// === API helper ===
async function apiRequest(endpoint, options = {}) {
    const token = localStorage.getItem('auth_token');

    const response = await fetch(endpoint, {
        ...options,
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Authorization': `Bearer ${token}`,
            ...options.headers
        }
    });

    if (response.status === 401) {
        window.location.href = '/login';
        return;
    }

    if (!response.ok) {
        const err = await response.json().catch(() => ({}));
        throw new Error(err.error || `Ошибка ${response.status}`);
    }

    return response.json();
}

// === Рендер пользователя ===
function renderUser(user, showAction = true) {
    return `
        <li class="user-item">
            <div class="user-info">
                <div class="user-name">${escapeHtml(user.username)}</div>
                <div class="user-meta">
                    ${user.city ? '📍 ' + escapeHtml(user.city) : ''}
                </div>
            </div>
            <div>
                ${showAction 
                    ? renderAction(user) 
                    : `<span class="status accepted">Друг</span>`
                }
            </div>
        </li>
    `;
}

// === Действия ===
function renderAction(user) {
    switch (user.friend_status) {
        case 'accepted':
            return `<span class="status accepted">Друг</span>`;
        case 'pending':
            return `<span class="status pending">Заявка отправлена</span>`;
        case 'incoming':
            return `<button class="btn btn-primary" onclick="acceptFriend(${user.id})">Принять</button>`;
        default:
            return `<button class="btn btn-primary" onclick="addFriend(${user.id})">Добавить</button>`;
    }
}

// === Ошибки ===
function showError(msg) {
    const el = document.getElementById('error');
    if (!el) return;

    el.textContent = msg;
    el.style.display = 'block';

    setTimeout(() => el.style.display = 'none', 5000);
}

// === Безопасный HTML ===
function escapeHtml(text) {
    if (!text) return '';
    const map = { '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#039;' };
    return text.replace(/[&<>"']/g, m => map[m]);
}

// === Загрузка всех пользователей ===
async function loadAllUsers() {
    const content = document.getElementById('content');
    content.innerHTML = '<div class="loading">Загрузка...</div>';

    try {
        const data = await apiRequest('/api/people?limit=50');
        const users = data.users; // ✅ фикс

        if (!users || users.length === 0) {
            content.innerHTML = '<div class="empty">Пользователи не найдены</div>';
            return;
        }

        content.innerHTML = `
            <ul class="user-list">
                ${users.map(u => renderUser(u)).join('')}
            </ul>
        `;
    } catch (e) {
        showError(e.message);
        content.innerHTML = '<div class="empty">Ошибка загрузки</div>';
    }
}

// === Загрузка друзей ===
async function loadFriends() {
    const content = document.getElementById('content');
    content.innerHTML = '<div class="loading">Загрузка друзей...</div>';
    try {
        const data = await apiRequest('/api/friends'); // или /friends
        
        const friends = data.users || data.friends || data; 
        
        if (!Array.isArray(friends) || friends.length === 0) {
            content.innerHTML = '<div class="empty">У вас пока нет друзей</div>';
            return;
        }
        content.innerHTML = `<ul class="user-list">${friends.map(u => renderUser(u, false)).join('')}</ul>`;
    } catch (e) {
        showError(e.message);
        content.innerHTML = '<div class="empty">Ошибка загрузки</div>';
    }
}

// === Добавить в друзья ===
async function addFriend(friendId) {
    try {
        await apiRequest('/api/friends', {
            method: 'POST',
            body: JSON.stringify({ friend_id: friendId })
        });
        loadAllUsers();
    } catch (e) {
        showError(e.message);
    }
}

// === Принять заявку ===
async function acceptFriend(friendId) {
    try {
        await apiRequest('/api/friends/accept', {
            method: 'POST',
            body: JSON.stringify({ friend_id: friendId })
        });
        loadAllUsers();
    } catch (e) {
        showError(e.message);
    }
}

// === Tabs ===
function initTabs() {
    document.querySelectorAll('.tab').forEach(tab => {
        tab.addEventListener('click', (e) => {
            document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
            e.currentTarget.classList.add('active');

            const tabName = e.currentTarget.dataset.tab;
            tabName === 'all' ? loadAllUsers() : loadFriends();
        });
    });
}

// === INIT ===
document.addEventListener('DOMContentLoaded', () => {
    createSnowflakes();
    initTabs();
    loadAllUsers();
});
