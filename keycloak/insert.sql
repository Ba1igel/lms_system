-- =============================================================================
-- LT10: Keycloak — прямые SQL-вставки в PostgreSQL (образовательный файл)
-- =============================================================================
-- ВАЖНО: Этот файл показывает, как Keycloak хранит данные внутри PostgreSQL.
-- В реальных проектах используйте realm-export.json (--import-realm),
-- а не прямые SQL-вставки. Прямой SQL нужен только для понимания структуры
-- или нестандартных миграций.
-- =============================================================================

-- =============================================================================
-- СТРУКТУРА КЛЮЧЕВЫХ ТАБЛИЦ KEYCLOAK
-- =============================================================================

-- 1. REALM — конфигурация realm (создаётся Keycloak автоматически)
--    id        : UUID realm
--    name      : 'lms-realm'

-- 2. USER_ENTITY — пользователи
--    id        : UUID (генерируем сами через gen_random_uuid())
--    realm_id  : ссылка на realm
--    username  : логин
--    email     : email
--    enabled   : true/false

-- 3. KEYCLOAK_ROLE — роли
--    id        : UUID
--    name      : 'admin' / 'teacher' / 'student'
--    realm_id  : ссылка на realm
--    client_role: false (realm role), true (client role)

-- 4. USER_ROLE_MAPPING — связь юзер ↔ роль
--    user_id   : UUID из user_entity
--    role_id   : UUID из keycloak_role

-- 5. CREDENTIAL — пароли (хранятся захэшированными!)
--    id            : UUID
--    type          : 'password'
--    user_id       : UUID из user_entity
--    secret_data   : JSON с хэшем и солью
--    credential_data: JSON с алгоритмом и итерациями

-- =============================================================================
-- КАК KEYCLOAK ХЭШИРУЕТ ПАРОЛИ
-- =============================================================================
-- Алгоритм: PBKDF2-SHA256
-- Итерации: 27500
-- Хранится в credential.secret_data как JSON:
--   {"value":"<base64-hash>","salt":"<base64-salt>"}
-- Credential_data:
--   {"hashIterations":27500,"algorithm":"pbkdf2-sha256","additionalParameters":{}}
--
-- => Нельзя сгенерировать валидный хэш в чистом SQL!
--    Используйте realm-export.json (пароль в открытом виде, Keycloak сам хэширует при импорте)
--    ИЛИ Keycloak Admin REST API для установки паролей.
-- =============================================================================

-- =============================================================================
-- ПРИМЕР: Посмотреть данные после запуска Keycloak с realm-export.json
-- =============================================================================

-- Посмотреть все realm:
SELECT id, name, enabled FROM realm;

-- Посмотреть всех пользователей в нашем realm:
SELECT
    ue.id,
    ue.username,
    ue.email,
    ue.enabled,
    ue.created_timestamp
FROM user_entity ue
JOIN realm r ON ue.realm_id = r.id
WHERE r.name = 'lms-realm';

-- Посмотреть все роли realm:
SELECT
    kr.id,
    kr.name,
    kr.client_role
FROM keycloak_role kr
JOIN realm r ON kr.realm_id = r.id
WHERE r.name = 'lms-realm'
  AND kr.client_role = false;

-- Посмотреть какие роли у каждого пользователя:
SELECT
    ue.username,
    ue.email,
    kr.name AS role
FROM user_role_mapping urm
JOIN user_entity ue ON urm.user_id = ue.id
JOIN keycloak_role kr ON urm.role_id = kr.id
JOIN realm r ON ue.realm_id = r.id
WHERE r.name = 'lms-realm';

-- =============================================================================
-- ПРЯМАЯ ВСТАВКА РОЛИ (без пароля — для служебных записей)
-- =============================================================================
-- Этот блок показывает структуру. Выполняйте ПОСЛЕ запуска Keycloak,
-- когда realm 'lms-realm' уже создан через realm-export.json.

-- Вставить дополнительную роль 'moderator':
DO $$
DECLARE
    v_realm_id VARCHAR(36);
BEGIN
    SELECT id INTO v_realm_id FROM realm WHERE name = 'lms-realm';

    -- Вставляем роль только если её ещё нет
    IF NOT EXISTS (
        SELECT 1 FROM keycloak_role
        WHERE name = 'moderator' AND realm_id = v_realm_id AND client_role = false
    ) THEN
        INSERT INTO keycloak_role (id, name, realm_id, client_role, client_realm_constraint)
        VALUES (
            gen_random_uuid()::varchar,
            'moderator',
            v_realm_id,
            false,
            v_realm_id
        );
        RAISE NOTICE 'Role moderator created';
    ELSE
        RAISE NOTICE 'Role moderator already exists';
    END IF;
END $$;

-- =============================================================================
-- НАЗНАЧИТЬ РОЛЬ СУЩЕСТВУЮЩЕМУ ПОЛЬЗОВАТЕЛЮ
-- =============================================================================

DO $$
DECLARE
    v_realm_id VARCHAR(36);
    v_user_id  VARCHAR(36);
    v_role_id  VARCHAR(36);
BEGIN
    SELECT id INTO v_realm_id FROM realm WHERE name = 'lms-realm';

    -- Найти пользователя
    SELECT id INTO v_user_id
    FROM user_entity
    WHERE username = 'student_alice' AND realm_id = v_realm_id;

    -- Найти роль
    SELECT id INTO v_role_id
    FROM keycloak_role
    WHERE name = 'moderator' AND realm_id = v_realm_id AND client_role = false;

    -- Назначить роль (если ещё не назначена)
    IF v_user_id IS NOT NULL AND v_role_id IS NOT NULL THEN
        INSERT INTO user_role_mapping (user_id, role_id)
        VALUES (v_user_id, v_role_id)
        ON CONFLICT DO NOTHING;
        RAISE NOTICE 'Role assigned successfully';
    ELSE
        RAISE NOTICE 'User or role not found';
    END IF;
END $$;

-- =============================================================================
-- ОТКЛЮЧИТЬ ПОЛЬЗОВАТЕЛЯ
-- =============================================================================

UPDATE user_entity
SET enabled = false
WHERE username = 'student_alice'
  AND realm_id = (SELECT id FROM realm WHERE name = 'lms-realm');

-- =============================================================================
-- ИТОГ: ЧТО ИСПОЛЬЗОВАТЬ В ПРОЕКТЕ
-- =============================================================================
-- ✅ realm-export.json + --import-realm  → первичная загрузка (роли + юзеры + клиенты)
-- ✅ Keycloak Admin REST API             → динамическое управление (создать/удалить юзера)
-- ✅ Прямой SQL (этот файл)             → SELECT для мониторинга + точечные UPDATE/INSERT
-- ❌ Прямой SQL для паролей             → НЕЛЬЗЯ, пароли нужно хэшировать через Keycloak
-- =============================================================================
