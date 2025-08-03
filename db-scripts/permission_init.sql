UPDATE users SET activated = true WHERE email = '';


-- すべてのユーサに「movie:read」権限を与える.
INSERT INTO users_permissions SELECT id, (SELECT id FROM permissions WHERE code = 'movies:read') FROM users;


-- 「memes@gmail」に["movies:write"]権限を与える
INSERT INTO users_permissions
VALUES (
		(SELECT id FROM users WHERE users.email='konas@memes.com'),
		(SELECT id FROM permissions WHERE permissions.code='movies:write')
		);


-- List all activated users and their permissions
SELECT email, array_agg(p.code) as permi
FROM permissions as p
INNER JOIN users_permissions as up
ON up.permission_id = p.id 
INNER JOIN users as u 
ON u.id = up.user_id
WHERE u.activated = true
GROUP BY email;


