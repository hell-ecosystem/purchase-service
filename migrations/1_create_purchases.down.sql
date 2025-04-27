-- Удаляем индексы сначала
DROP INDEX IF EXISTS idx_purchases_is_active;
DROP INDEX IF EXISTS idx_purchases_category;
DROP INDEX IF EXISTS idx_purchases_user_id;

-- Удаляем таблицу
DROP TABLE IF EXISTS purchases;

-- (Опционально) Отключение расширения uuid-ossp, если хочешь чистить БД до конца:
-- DROP EXTENSION IF EXISTS "uuid-ossp";
