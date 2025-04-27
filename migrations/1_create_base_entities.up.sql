-- Включаем расширение для UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Справочник категорий товаров
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- уникальный ID категории
    name TEXT NOT NULL UNIQUE, -- название категории
    created_at TIMESTAMP NOT NULL DEFAULT now() -- дата создания категории
);

-- Справочник магазинов
CREATE TABLE IF NOT EXISTS shops (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- уникальный ID магазина
    name TEXT NOT NULL UNIQUE, -- название магазина
    address TEXT, -- адрес магазина
    created_at TIMESTAMP NOT NULL DEFAULT now() -- дата создания магазина
);

-- Таблица покупок
CREATE TABLE IF NOT EXISTS purchases (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- уникальный ID покупки
    user_id UUID NOT NULL, -- ID пользователя
    name TEXT NOT NULL, -- название товара
    category_id UUID, -- ID категории товара
    shop_id UUID, -- ID магазина
    estimated_price NUMERIC(10,2), -- предполагаемая цена
    actual_price NUMERIC(10,2), -- реальная цена
    created_at TIMESTAMP NOT NULL DEFAULT now(), -- дата создания покупки
    deadline TIMESTAMP, -- срок до которого нужно купить
    is_active BOOLEAN NOT NULL DEFAULT true, -- флаг актуальности
    is_purchased BOOLEAN NOT NULL DEFAULT false, -- флаг покупки
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(id) ON DELETE SET NULL
);

-- Индексы для оптимизации поиска по покупкам
CREATE INDEX idx_purchases_user_id ON purchases(user_id); -- поиск покупок по пользователю
CREATE INDEX idx_purchases_is_active ON purchases(is_active); -- активные покупки
CREATE INDEX idx_purchases_category_id ON purchases(category_id); -- поиск по категории
CREATE INDEX idx_purchases_shop_id ON purchases(shop_id); -- поиск по магазину

-- Бюджеты пользователей
CREATE TABLE IF NOT EXISTS budgets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- уникальный ID бюджета
    user_id UUID NOT NULL, -- ID пользователя
    month INT NOT NULL, -- месяц
    year INT NOT NULL, -- год
    amount NUMERIC(10,2) NOT NULL, -- сумма бюджета
    created_at TIMESTAMP NOT NULL DEFAULT now(), -- дата создания бюджета
    UNIQUE(user_id, month, year) -- один бюджет на месяц для одного пользователя
);

-- Напоминания о покупках
CREATE TABLE IF NOT EXISTS reminders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- уникальный ID напоминания
    purchase_id UUID NOT NULL, -- ID покупки
    remind_at TIMESTAMP NOT NULL, -- время напоминания
    message TEXT, -- сообщение напоминания
    created_at TIMESTAMP NOT NULL DEFAULT now(), -- дата создания напоминания
    CONSTRAINT fk_reminder_purchase FOREIGN KEY (purchase_id) REFERENCES purchases(id) ON DELETE CASCADE
);

-- История действий с покупками
CREATE TABLE IF NOT EXISTS purchase_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- уникальный ID истории
    purchase_id UUID NOT NULL, -- ID покупки
    action TEXT NOT NULL, -- действие (created, updated, purchased, deleted)
    description TEXT, -- описание действия
    created_at TIMESTAMP NOT NULL DEFAULT now(), -- время действия
    CONSTRAINT fk_history_purchase FOREIGN KEY (purchase_id) REFERENCES purchases(id) ON DELETE CASCADE
);

-- Справочник меток
CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- уникальный ID метки
    name TEXT NOT NULL UNIQUE, -- название метки
    created_at TIMESTAMP NOT NULL DEFAULT now() -- дата создания метки
);

-- Связь многие-ко-многим покупки и метки
CREATE TABLE IF NOT EXISTS purchase_tags (
    purchase_id UUID NOT NULL, -- ID покупки
    tag_id UUID NOT NULL, -- ID метки
    PRIMARY KEY (purchase_id, tag_id), -- составной ключ
    CONSTRAINT fk_purchase_tag_purchase FOREIGN KEY (purchase_id) REFERENCES purchases(id) ON DELETE CASCADE,
    CONSTRAINT fk_purchase_tag_tag FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);
