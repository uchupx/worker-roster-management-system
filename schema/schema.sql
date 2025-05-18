-- Drop the table if it exists
DROP TABLE IF EXISTS users;
CREATE TABLE  users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    deleted_at DATETIME
);
INSERT INTO users (name, email, password) VALUES
('admin', 'admin@uchupx.tech', 'GBeIP-fo4jVVAtaKQq7HkEmXuBrBcDeIiD5IK8Fr-VCGPqEXomLasxmIu3aBlxIEv_6SN9v1slm-Gt7Q7h0f9S8A49EjhkELWJlY55bHhlTQMECq8Kj3p_N8satwFlZoTWhv7mj8X7fs_-bLat3keWHoxggq3uk365zZQAZ98DUp3SOWCfe6ruS0MIVwhr5l1z6SvxaUKSDTH-aC3Tnw3g5ZvNCv5ZkHeTx5cb6CKC8Kpe8ynDsYIPDW3PJJBBCl1VLhZ7IiGdpWPpknQSbV4yR5gPy4NaS7a9at9HPibcGej78hDTh2XppXDf3jvQaLfPxVDtGOp_OtKQ_ucb2lhA=='),
('John D', 'employee@uchupx.tech', 'GBeIP-fo4jVVAtaKQq7HkEmXuBrBcDeIiD5IK8Fr-VCGPqEXomLasxmIu3aBlxIEv_6SN9v1slm-Gt7Q7h0f9S8A49EjhkELWJlY55bHhlTQMECq8Kj3p_N8satwFlZoTWhv7mj8X7fs_-bLat3keWHoxggq3uk365zZQAZ98DUp3SOWCfe6ruS0MIVwhr5l1z6SvxaUKSDTH-aC3Tnw3g5ZvNCv5ZkHeTx5cb6CKC8Kpe8ynDsYIPDW3PJJBBCl1VLhZ7IiGdpWPpknQSbV4yR5gPy4NaS7a9at9HPibcGej78hDTh2XppXDf3jvQaLfPxVDtGOp_OtKQ_ucb2lhA=='),
('Luffy D', 'luffy@uchupx.tech', 'GBeIP-fo4jVVAtaKQq7HkEmXuBrBcDeIiD5IK8Fr-VCGPqEXomLasxmIu3aBlxIEv_6SN9v1slm-Gt7Q7h0f9S8A49EjhkELWJlY55bHhlTQMECq8Kj3p_N8satwFlZoTWhv7mj8X7fs_-bLat3keWHoxggq3uk365zZQAZ98DUp3SOWCfe6ruS0MIVwhr5l1z6SvxaUKSDTH-aC3Tnw3g5ZvNCv5ZkHeTx5cb6CKC8Kpe8ynDsYIPDW3PJJBBCl1VLhZ7IiGdpWPpknQSbV4yR5gPy4NaS7a9at9HPibcGej78hDTh2XppXDf3jvQaLfPxVDtGOp_OtKQ_ucb2lhA==');

DROP TABLE IF EXISTS roles;
CREATE TABLE roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME
);


INSERT INTO roles (name) VALUES
('admin'),
('employee');

DROP TABLE IF EXISTS user_roles;
CREATE TABLE user_roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

INSERT INTO user_roles (user_id, role_id) VALUES
(1, 1),
(1, 2),
(2, 2),
(3, 2);

DROP TABLE IF EXISTS shifts;
CREATE TABLE shifts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    start_time TEXT NOT NULL,
    end_time TEXT NOT NULL,
    shift_date DATE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status INTEGER DEFAULT 0, -- 0 is pending, 1 is approved, 2 is rejected
    updated_at DATETIME,
    deleted_at DATETIME,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO shifts (user_id, start_time, end_time, shift_date, status) VALUES
(2, '08:00:00', '16:00:00', '2025-05-21', 0),
(3, '08:00:00', '16:00:00', '2025-05-21', 0),
(2, '16:00:00', '00:00:00', '2025-05-22', 0),
(3, '16:00:00', '00:00:00', '2025-05-22', 0),
(1, '00:00:00', '08:00:00', '2023-10-03', 1);
