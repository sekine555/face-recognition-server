-- mst_user
-- 一般ユーザー
INSERT INTO mst_user(password, email, username, photo, s3_key, created_at)
    VALUES ('xxxx', 'test1@test.co.jp', 'テスト太郎1', 'xxx.jpg', 'sasakinozomi-smile.jpg', CURRENT_TIMESTAMP);
-- 管理者
INSERT INTO mst_user(password, email, username, photo, s3_key, is_admin, created_at)
    VALUES ('xxxx', 'admin@test.co.jp', '管理者太郎1', 'xxx.jpg', 'sasakinozomi-smile.jpg', true, CURRENT_TIMESTAMP);