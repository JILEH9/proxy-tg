# proxy-tg

> HTTP-прокси на Go для перенаправления Telegram вебхуков на любой внутренний сервис (n8n и др.)

Принимает входящий запрос от Telegram и проксирует его на адрес из параметра `target_url`, сохраняя метод, заголовки и тело. Работает за Traefik с автоматическим HTTPS через Let's Encrypt.

---

## Требования

1. Сервер в Европе (VPS)
2. Домен — желательно зарегистрированный не в России, например через [Cloudflare](https://cloudflare.com)
3. A-запись домена или субдомена, направленная на IP сервера
4. Docker >= 20.10

---

## Установка

```bash
git clone https://github.com/JILEH9/proxy-tg.git
cd proxy-tg
cp .env.example .env
nano .env                  # укажите DOMAIN_NAME и SSL_EMAIL
docker compose up -d
```

---

## Настройка вебхука Telegram

Выполните POST-запрос, чтобы направить вебхуки Telegram через прокси:

```
POST https://api.telegram.org/bot<TOKEN>/setWebhook
Content-Type: application/json

{
  "url": "https://your-domain.com?target_url=https://n8n.your-domain.com"
}
```

Все входящие апдейты от Telegram будут перенаправлены на адрес в `target_url`.

---

## Переменные окружения

| Переменная    | Описание                               |
|---------------|----------------------------------------|
| `DOMAIN_NAME` | Домен, на котором работает прокси      |
| `SSL_EMAIL`   | Email для уведомлений Let's Encrypt    |

---

## Стек

- **Go 1.23**
- **Traefik v3** — reverse proxy + автоматический TLS
- **Docker Compose**

---

## Контакты

- Telegram: [@Miirrox](https://t.me/Miirrox)
- GitHub: [JILEH9](https://github.com/JILEH9)

---

## Лицензия

[MIT](./LICENSE)
