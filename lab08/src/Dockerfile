FROM alpine:3.18 AS builder

# Установка Python, pip и необходимых библиотек
RUN apk add --no-cache python3 py3-pip py3-virtualenv go

# Настройка рабочего каталога
WORKDIR /app

# Копирование файлов
COPY users_gen.py /app/users_gen.py
COPY gen_scheduler.go /app/gen_scheduler.go

# Создание виртуального окружения и установка зависимостей
RUN python3 -m venv /app/venv && \
    /app/venv/bin/pip install --no-cache-dir faker

# Сборка Go-приложения
RUN go build -o /app/scheduler /app/gen_scheduler.go

# Финальный образ
FROM alpine:3.18 AS runtime

# Установка Python
RUN apk add --no-cache python3 py3-pip

# Настройка рабочего каталога
WORKDIR /app

# Копирование артефактов
COPY --from=builder /app/scheduler /app/scheduler
COPY --from=builder /app/users_gen.py /app/users_gen.py
COPY --from=builder /app/venv /app/venv

# Настройка окружения PATH
ENV PATH="/app/venv/bin:$PATH"

# Команда запуска
CMD ["/app/scheduler"]
