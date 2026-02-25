document.addEventListener('DOMContentLoaded', async () => {
    const token = localStorage.getItem('auth_token');
    
    if (!token) {
        window.location.href = '/login.html';
        return;
    }
    
    // Всегда загружаем свежие данные с сервера
    await loadProfileData();
    await loadPosts();
});    // Загрузка профиля
    async function loadProfileData() {
        const token = localStorage.getItem('auth_token');
        
        try {
            const response = await fetch('/dashboard', {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (response.status === 401) {
                window.location.href = '/login.html';
                return;
            }

            const data = await response.json();
            
            document.getElementById('profileName').textContent = data.username || 'Гость';
            document.getElementById('profileStatus').textContent = data.email || '';
            document.getElementById('profileAvatar').textContent = 
                (data.username || 'A').charAt(0).toUpperCase();
            
            if (data.location) document.getElementById('location').textContent = data.location;
            if (data.birthday) document.getElementById('birthday').textContent = data.birthday;
            if (data.work) document.getElementById('work').textContent = data.work;

        } catch (error) {
            console.error('Ошибка сети:', error);
        }
    }

    // Создание поста
    document.getElementById('createPost').addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const textarea = document.getElementById('newPostText');
        const text = textarea.value.trim();
        if (!text) return;

        const token = localStorage.getItem('auth_token');
        
        try {
            const response = await fetch('/api/posts', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ text: text })
            });

            if (response.status === 401) {
                window.location.href = '/login.html';
                return;
            }

            if (!response.ok) {
                throw new Error(`Ошибка ${response.status}`);
            }

            textarea.value = '';
            await loadPosts();            
        } catch (err) {
            console.error('Ошибка:', err);
            alert(`Не удалось опубликовать: ${err.message}`);
        }
    });



async function loadPosts() {
    const token = localStorage.getItem('auth_token');
    const container = document.getElementById('postsContainer');
    
    try {
        const response = await fetch('/api/posts', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        const rawText = await response.text();

        
        // Пытаемся распарсить
        const posts = JSON.parse(rawText);
        
        if (response.status === 401) {
            window.location.href = '/login.html';
            return;
        }

        if (!response.ok) {
            throw new Error(`Ошибка ${response.status}`);
        }
        
        container.innerHTML = '';

        if (posts.length === 0) {
            container.innerHTML = `
                <div style="text-align: center; color: var(--mountain); padding: 40px;">
                    <i class="fas fa-inbox" style="font-size: 48px; margin-bottom: 15px; opacity: 0.5;"></i>
                    <div>Пока нет постов. Будьте первым! ❄️</div>
                </div>
            `;
            return;
        }

        posts.forEach(post => {
            const postHTML = createPostCard(post);
            container.insertAdjacentHTML('beforeend', postHTML);
        });

    } catch (error) {
        console.error(' Ошибка загрузки постов:', error);
        container.innerHTML = `
            <div style="text-align: center; color: var(--mountain); padding: 40px;">
                <i class="fas fa-exclamation-triangle" style="font-size: 24px; margin-bottom: 10px;"></i>
                <div>Не удалось загрузить посты</div>
            </div>
        `;
    }
}

// Функция создания HTML карточки поста
function createPostCard(post) {
    const timeAgo = getTimeAgo(new Date(post.created_at));
    const avatarLetter = (post.username || 'U').charAt(0).toUpperCase();
    
    return `
        <div class="post-card">
            <div class="post-header">
                <div class="post-avatar">${avatarLetter}</div>
                <div class="post-user">
                    <div class="post-username">${escapeHtml(post.username || 'Аноним')}</div>
                    <div class="post-time">${timeAgo}</div>
                </div>
            </div>
            <div class="post-content">${escapeHtml(post.text)}</div>
            <div class="post-footer">
                <div class="post-actions">
                    <div class="post-action like" onclick="toggleLike(${post.id})">
                        <i class="far fa-heart"></i>
                        <span>${post.likes || 0}</span>
                    </div>
                    <div class="post-action">
                        <i class="far fa-comment"></i>
                        <span>${post.comments || 0}</span>
                    </div>
                </div>
            </div>
        </div>
    `;
}

function getTimeAgo(date) {
    const seconds = Math.floor((new Date() - date) / 1000);
    
    if (seconds < 60) return 'Только что';
    if (seconds < 3600) return `${Math.floor(seconds / 60)} мин назад`;
    if (seconds < 86400) return `${Math.floor(seconds / 3600)} час(а) назад`;
    return `${Math.floor(seconds / 86400)} дн назад`;
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

window.addEventListener('DOMContentLoaded', () => {
    loadProfileData();
    loadPosts(); 
});
    window.addEventListener('DOMContentLoaded', loadProfileData);
