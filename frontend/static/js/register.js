    // Создание снежной пыли 
    function createSnowDust() {
        const container = document.getElementById('snowDust');
        const count = 80;
        
        for (let i = 0; i < count; i++) {
            const particle = document.createElement('div');
            particle.classList.add('snow-particle');
            
            const size = Math.random() * 4 + 1;
            const posX = Math.random() * 100;
            const duration = Math.random() * 30 + 15;
            const delay = Math.random() * 10;
            const drift = (Math.random() - 0.5) * 30;
            
            particle.style.width = `${size}px`;
            particle.style.height = `${size}px`;
            particle.style.left = `${posX}%`;
            particle.style.animationDuration = `${duration}s`;
            particle.style.animationDelay = `${delay}s`;
            particle.style.setProperty('--offset', drift);
            
            container.appendChild(particle);
        }
    }

    // Плавное появление гор 
    function initMountains() {
        document.querySelectorAll('.peak').forEach((peak, index) => {
            peak.style.opacity = '0';
            peak.style.transform = `translateY(40px) ${peak.style.transform || ''}`;
            
            setTimeout(() => {
                peak.style.transition = 'opacity 0.8s ease-out, transform 1s ease-out';
                peak.style.opacity = '0.85';
                peak.style.transform = `translateY(0) ${peak.style.transform.replace('translateY(40px)', '')}`;
            }, 300 + index * 200);
        });
    }

    // отправка формы в микросервис
   async function registerUser(formData) {
    try {
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData)
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.message || `Ошибка ${response.status}: ${response.statusText}`);
        }

        return await response.json();
    } catch (error) {
        throw error;
    }
}

// Обработка отправки формы
document.getElementById('registerForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    // Собираем данные формы (БЕЗ name)
    const formData = {
        email: document.getElementById('email').value.trim(),
        password: document.getElementById('password').value,
        username: document.getElementById('username').value.trim()
    };

    // Простая валидация на клиенте (БЕЗ name)
    if (!formData.email || !formData.password || !formData.username) {
        alert('Пожалуйста, заполните все поля');
        return;
    }
    
    if (formData.password.length < 4) { // или 8, как у вас в требованиях
        alert('Пароль должен быть не менее 4 символов');
        return;
    }

    // Визуальная обратная связь
    const btn = this.querySelector('.submit-btn');
    const originalText = btn.innerHTML;
    btn.innerHTML = '<span class="snowflake">❄</span> Регистрация...';
    btn.disabled = true;

    try {
        // Отправляем данные в микросервис
        const result = await registerUser(formData);
        
        // Успешная регистрация
        btn.innerHTML = '<span class="snowflake">✓</span> Успешно!';
        setTimeout(() => {
            alert('Аккаунт создан! Теперь вы можете войти.');
            window.location.href = '/login.html'; 
        }, 1500);

    } catch (error) {
        // Обработка ошибок
        console.error('Ошибка регистрации:', error);
        btn.innerHTML = '<span class="snowflake">⚠</span> Ошибка!';
        setTimeout(() => {
            btn.innerHTML = originalText;
            btn.disabled = false;
            alert(`Ошибка регистрации: ${error.message}`);
        }, 2000);
    }
});

    // Инициализация
    document.addEventListener('DOMContentLoaded', () => {
        createSnowDust();
        initMountains();
    });