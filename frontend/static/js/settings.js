    //  Снежинки
        function createSnow() {
            const c = document.getElementById('snow');
            for (let i = 0; i < 45; i++) {
                const s = document.createElement('div');
                s.className = 'snowflake';
                const sz = Math.random() * 3 + 1.5;
                s.style.width = s.style.height = sz + 'px';
                s.style.left = Math.random() * 100 + '%';
                s.style.animationDuration = (Math.random() * 10 + 7) + 's';
                s.style.animationDelay = (Math.random() * 4) + 's';
                c.appendChild(s);
            }
        }

        //  Показать/скрыть пароль
        function togglePass(id, btn) {
            const inp = document.getElementById(id);
            const ic = btn.querySelector('i');
            if (inp.type === 'password') {
                inp.type = 'text';
                ic.className = 'fas fa-eye-slash';
            } else {
                inp.type = 'password';
                ic.className = 'fas fa-eye';
            }
        }

        // валидация поля
        function validate(input, errId, fn) {
            const err = document.getElementById(errId);
            const val = input.value.trim();
            const ok = fn(val);
            if (!ok && val !== '') {
                input.classList.add('invalid');
                if (err) err.classList.add('show');
                return false;
            }
            input.classList.remove('invalid');
            if (err) err.classList.remove('show');
            return true;
        }

        //  Универсальная отправка отдельного блока
        async function submitBlock(formEl, btnEl, statusEl, getDataFn, validateFn) {
            statusEl.className = 'block-status';
            statusEl.style.display = 'none';

            // Валидация
            if (validateFn && !validateFn()) return;

            // UI: загрузка
            btnEl.classList.add('loading');
            btnEl.disabled = true;
            const originalText = btnEl.querySelector('.btn-text').textContent;

            try {
                const data = getDataFn();
                const endpoint = formEl.dataset.endpoint;
                const url = `http://localhost:8080${endpoint}`;

                console.log(`[${endpoint}] Отправка:`, data);

                const res = await fetch(url, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(data)
                });

                if (!res.ok) {
                    const txt = await res.text();
                    let emsg = `Ошибка ${res.status}`;
                    try { const j = JSON.parse(txt); emsg = j.message || emsg; } catch(_) {}
                    throw new Error(emsg);
                }

                await res.json();
                showStatus(statusEl, 'success', ' Сохранено!');

                // Очистка паролей после успеха
                if (formEl.id === 'formPass') {
                    document.getElementById('inputOldPass').value = '';
                    document.getElementById('inputNewPass').value = '';
                }

            } catch (err) {
                console.error(` [${formEl.dataset.endpoint}]`, err);
                showStatus(statusEl, 'error', + (err.message || 'Не удалось сохранить'));
            } finally {
                btnEl.classList.remove('loading');
                btnEl.disabled = false;
            }
        }

        //  Показать статус в блоке
        function showStatus(el, type, text) {
            el.textContent = text;
            el.className = `block-status ${type}`;
            el.style.display = 'block';
            setTimeout(() => { el.style.display = 'none'; }, 3500);
        }

        // Инициализация обработчиков для каждого блока
        window.addEventListener('load', () => {
            createSnow();

        document.getElementById('formName').addEventListener('submit', (e) => {
            e.preventDefault();
            const btn = document.getElementById('btnName');
            const status = document.getElementById('statusName');
            const input = document.getElementById('inputName');
            
            submitBlock(
                e.currentTarget, btn, status,
                () => ({ newName: input.value.trim() }), 
                () => validate(input, 'errName', v => v.length >= 2)
            );
        });
        document.getElementById('inputName').addEventListener('blur', function() {
            validate(this, 'errName', v => v.length >= 2);
        });

        document.getElementById('formWork').addEventListener('submit', (e) => {
            e.preventDefault();
            const btn = document.getElementById('btnWork');
            const status = document.getElementById('statusWork');
            const input = document.getElementById('inputWork');
            
            submitBlock(
                e.currentTarget, btn, status,
                () => ({ work_location: input.value.trim() || null }),
                null 
            );
        });

        document.getElementById('formPass').addEventListener('submit', (e) => {
            e.preventDefault();
            const btn = document.getElementById('btnPass');
            const status = document.getElementById('statusPass');
            const oldPass = document.getElementById('inputOldPass');
            const newPass = document.getElementById('inputNewPass');
            
            submitBlock(
                e.currentTarget, btn, status,
                () => ({
                    OldPass: oldPass.value || null,
                    NewPass: newPass.value || null
                }),
                () => {
                    if (!newPass.value) return true; 
                    if (newPass.value.length < 8) {
                        newPass.classList.add('invalid');
                        document.getElementById('errPass').classList.add('show');
                        return false;
                    }
                    newPass.classList.remove('invalid');
                    document.getElementById('errPass').classList.remove('show');
                    return true;
                }
            );
        });
        document.getElementById('inputNewPass').addEventListener('input', function() {
            if (this.value && this.value.length < 8) {
                this.classList.add('invalid');
                document.getElementById('errPass').classList.add('show');
            } else {
                this.classList.remove('invalid');
                document.getElementById('errPass').classList.remove('show');
            }
        });

        document.getElementById('formCity').addEventListener('submit', (e) => {
            e.preventDefault();
            const btn = document.getElementById('btnCity');
            const status = document.getElementById('statusCity');
            const input = document.getElementById('inputCity');
            
            submitBlock(
                e.currentTarget, btn, status,
                () => ({ city: input.value.trim() }),
                () => validate(input, 'errCity', v => v.length >= 2)
            );
        });
        document.getElementById('inputCity').addEventListener('blur', function() {
            validate(this, 'errCity', v => v.length >= 2);
        });

});
        document.querySelectorAll('.form-control').forEach(inp => {
            inp.addEventListener('keypress', (e) => {
                if (e.key === 'Enter') {
                    e.preventDefault();
                    inp.closest('form').requestSubmit();
                }
            });
        });