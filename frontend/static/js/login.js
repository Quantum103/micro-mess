   // ========================================
    // СОЗДАНИЕ СНЕЖИНОК
    // ========================================
    function createSnowflakes() {
        const snowContainer = document.getElementById('snow');
        const snowflakeCount = 50;
        
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

    // ========================================
    // ПОКАЗАТЬ/СКРЫТЬ ПАРОЛЬ
    // ========================================
    function togglePassword() {
        const input = document.getElementById('password');
        const button = event.currentTarget;
        const icon = button.querySelector('i');
        
        if (input.type === 'password') {
            input.type = 'text';
            icon.className = 'fas fa-eye-slash';
        } else {
            input.type = 'password';
            icon.className = 'fas fa-eye';
        }
    }

    // ========================================
    // ЧЕКБОКС "ЗАПОМНИТЬ МЕНЯ"
    // ========================================
    function toggleRemember() {
        const checkbox = document.getElementById('rememberCheckbox');
        checkbox.classList.toggle('checked');
    }

    // ========================================
    // ЗАБЫЛИ ПАРОЛЬ
    // ========================================
    function forgotPassword(e) {
        e.preventDefault();
        alert('Функция восстановления пароля будет доступна позже.');
    }


    // ========================================
    // ОБРАБОТКА ОТПРАВКИ ФОРМЫ 
    // ========================================
    document.getElementById('loginForm').addEventListener('submit', async function(e) {
        e.preventDefault();
        
        console.log(' Форма отправлена');
        
        // Собираем данные формы
        const usernameInput = document.getElementById('username');
        const passwordInput = document.getElementById('password');
        
        const username = usernameInput.value.trim();
        const password = passwordInput.value;

        
        if (!username) {
            alert(' Пожалуйста, введите логин или email');
            return;
        }
        
        if (!password) {
            alert(' Пожалуйста, введите пароль');
            return;
        }
        
        if (password.length < 4) {
            alert(' Пароль должен быть не менее 4 символов');
            return;
        }

        // Визуальная обратная связь
        const btn = document.getElementById('loginButton');
        const originalText = btn.innerHTML;
        btn.innerHTML = '❄ Выполняется вход...';
        btn.disabled = true;

        try {
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            const isEmail = emailRegex.test(username);
            
            //  данные для отправки
            const formData = {
                identifier: username,
                password: password
            };


            // Отправляем запрос
            const response = await fetch('http://localhost:8080/api/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData)
            });

<<<<<<< HEAD
=======
            console.log(' Ответ от сервера - статус:', response.status);
            console.log('📥Ответ от сервера - headers:', Object.fromEntries(response.headers.entries()));

>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56
            if (!response.ok) {
                const errorText = await response.text();
                console.error('Текст ошибки от сервера:', errorText);
                
                let errorMessage = `Ошибка ${response.status}`;
                try {
                    const errorData = JSON.parse(errorText);
                    errorMessage = errorData.message || errorMessage;
                } catch (e) {
<<<<<<< HEAD
=======
                    // Если не JSON, используем текст как есть
>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56
                    errorMessage = errorText || errorMessage;
                }
                
                throw new Error(errorMessage);
            }

            const result = await response.json();
            console.log('Успешный ответ:', result);
            
            //  УСПЕШНАЯ АВТОРИЗАЦИЯ
            if (result.token) {
                console.log(' Токен получен, сохраняем в localStorage');
                localStorage.setItem('auth_token', result.token);
<<<<<<< HEAD
                if (result.token && result.token !== "undefined") {
                    localStorage.setItem("token", result.token);
                } else {
                    console.error("Сервер вернул неправильный token:", result.token);
                }                
                localStorage.setItem('user_id', result.user_id);                
                console.log('Сохранён user_id:', result.user_id);   // ← Для отладки
=======
>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56
                
                btn.innerHTML = 'Успешно!';
                
                setTimeout(() => {
                    console.log('Редирект на dashboard');
                    window.location.href = '/dashboard.html';
                }, 1000);
            } else {
                throw new Error('Сервер не вернул токен');
            }

        } catch (error) {
            //  ОШИБКА
            console.error(' Ошибка входа:', error);
            console.error(' Стек ошибки:', error.stack);
            alert(` Ошибка входа: ${error.message}`);
            
            // Восстанавливаем кнопку
            setTimeout(() => {
                btn.innerHTML = originalText;
                btn.disabled = false;
            }, 2000);
        }
    });

    // ========================================
    // ИНИЦИАЛИЗАЦИЯ
    // ========================================
    window.onload = () => {
<<<<<<< HEAD
        console.log('Страница загружена');
=======
        console.log('🏠 Страница загружена');
>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56
        createSnowflakes();
        
        // Enter для отправки формы
        document.getElementById('password').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                console.log('⏎ Enter нажат, отправляем форму');
                document.getElementById('loginForm').requestSubmit();
            }
        });
    };