-- Откатить связи и таблицы в правильном порядке

DROP TABLE IF EXISTS purchase_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS purchase_history;
DROP TABLE IF EXISTS reminders;
DROP TABLE IF EXISTS budgets;
DROP INDEX IF EXISTS idx_purchases_shop_id;
DROP INDEX IF EXISTS idx_purchases_category_id;
DROP INDEX IF EXISTS idx_purchases_is_active;
DROP INDEX IF EXISTS idx_purchases_user_id;
DROP TABLE IF EXISTS purchases;
DROP TABLE IF EXISTS shops;
DROP TABLE IF EXISTS categories;

-- Можно оставить расширение uuid-ossp, если оно используется другими таблицами
-- Иначе можно удалить его явно:
-- DROP EXTENSION IF EXISTS "uuid-ossp";
