    // === ГЕНЕРАЦИЯ СНЕЖИНОК ===
    function createSnowflakes() {
        const snowContainer = document.getElementById('snow');
        const count = 50;
        for (let i = 0; i < count; i++) {
            const flake = document.createElement('div');
            flake.classList.add('snowflake');
            const size = Math.random() * 4 + 2;
            flake.style.width = `${size}px`;
            flake.style.height = `${size}px`;
            flake.style.left = `${Math.random() * 100}%`;
            flake.style.animationDuration = `${Math.random() * 10 + 5}s`;
            flake.style.animationDelay = `${Math.random() * 5}s`;
            snowContainer.appendChild(flake);
        }
    }

    // === TOGGLE PASSWORD ===
    function togglePass(btn) {
        const input = btn.parentElement.querySelector('input');
        const icon = btn.querySelector('i');
        if (input.type === 'password') {
            input.type = 'text';
            icon.className = 'fas fa-eye-slash';
        } else {
            input.type = 'password';
            icon.className = 'fas fa-eye';
        }
    }

    
    // === ОБРАБОТКА ФОРМ ===

        document.querySelectorAll('form').forEach(form => {
        form.addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const btn = this.querySelector('.save-btn');
            const originalText = btn.innerText;
            
            btn.innerText = 'СОХРАНЕНИЕ...';
            btn.disabled = true;

            const formData = new FormData(this);
            const data = Object.fromEntries(formData.entries());
            const formId = this.id;
            let apiUrl = '';


            // Маршрутизация
            switch(formId) {
                case 'geoForm':
                    apiUrl = '/change/city';
                    break;
                case 'nameForm':
                    apiUrl = '/change/username';
                    break;
                case 'passForm':
                    apiUrl = '/change/Pass';
                    break;
                case 'workForm':
                    apiUrl = '/change/job';
                    break;
            }

            try {
                const token = localStorage.getItem('auth_token');
                
                const headers = {
                    'Content-Type': 'application/json',
                };
                
                if (token) {
                    headers['Authorization'] = `Bearer ${token}`;
                }

                const response = await fetch(apiUrl, {
                    method: 'POST',
                    headers: headers,
                    body: JSON.stringify(data)
                });

                if (!response.ok) {
                    const errorText = await response.text();
                    if (response.status === 401 || response.status === 403) {
                        throw new Error('Сессия истекла. Требуется вход.');
                    }
                    throw new Error(errorText || `Ошибка ${response.status}`);
                }

                btn.innerText = 'ГОТОВО!';
                btn.style.background = 'linear-gradient(135deg, #4CAF50, #45a049)';
               


                    
                setTimeout(() => {
                    btn.innerText = originalText;
                    btn.disabled = false;
                    btn.style.background = '';
                }, 2000);

            } catch (error) {
                console.error('Ошибка:', error);
                btn.innerText = 'ОШИБКА';
                btn.style.background = '#d9534f';
                
                if (error.message.includes('Сессия')) {
                    localStorage.removeItem('auth_token');
                    setTimeout(() => window.location.href = '/index.html', 1500);
                } else {
                    alert(`⚠️ ${error.message}`);
                    setTimeout(() => {
                        btn.innerText = originalText;
                        btn.disabled = false;
                        btn.style.background = '';
                    }, 2000);
                }
            }
        });
    });
    // Запуск
    window.onload = createSnowflakes;
