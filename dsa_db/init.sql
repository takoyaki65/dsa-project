-- データベースの作成
CREATE DATABASE IF NOT EXISTS dsa CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- データベースを使用
USE dsa;

-- テーブルの作成
CREATE TABLE IF NOT EXISTS assignments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    test_dir_name VARCHAR(255) NOT NULL,
    start_date DATETIME NOT NULL, 
    end_date DATETIME NOT NULL
);

-- 初期データの挿入
INSERT INTO assignments (title, test_dir_name, start_date, end_date) VALUES
('課題1', 'report1', '2023-10-01 00:00:00', '2025-12-31 23:59:59'),
('課題2', 'report2', '2023-10-08 00:00:00', '2025-12-31 23:59:59'),
('課題3', 'report3', '2024-10-30 00:00:00', '2025-12-31 23:59:59'), -- まだ開始していないケース  
('課題4', 'report4', '2023-11-01 00:00:00', '2023-12-31 23:59:59'); -- もう終了しているケース  

CREATE TABLE IF NOT EXISTS sub_assignments (
    id INT NOT NULL,
    sub_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    test_dir_name VARCHAR(255) NOT NULL,
    makefile VARCHAR(255) NOT NULL,
    required_file_name VARCHAR(255) NOT NULL,
    main_file_name VARCHAR(255) NOT NULL,
    test_case_name VARCHAR(255) NULL,
    PRIMARY KEY (id, sub_id),
    CONSTRAINT fk_assignments_id FOREIGN KEY (id) 
        REFERENCES assignments(id)
) ENGINE=InnoDB;

INSERT INTO sub_assignments (id, sub_id, title, test_dir_name, makefile, required_file_name, main_file_name, test_case_name) VALUES
(1, 1, '基本課題1', 'sub1', 'gcd_euclid: gcd_euclid.o main_euclid.o', 'gcd_euclid.c', 'main_euclid.c', 'case1.txt'),
(1, 2, '発展課題1', 'sub2', 'gcd_recursive: gcd_recursive.o main_recursive.o', 'gcd_recursive.c', 'main_recursive.c', 'case1.txt'),
(2, 1, '基本課題1', 'sub3', 'linked_list: linked_list.o main_linked_list.o', 'linked_list.c', 'main_linked_list.c', 'case1.txt'),
(2, 2, '基本課題2', 'sub1', 'queue: queue.o main_queue.o', 'queue.c', 'main_queue.c', 'case1.txt'),
(2, 3, '発展課題1', 'sub2', 'doublylinked_list: doublylinked_list.o main_doublylinked_list.o', 'doublylinked_list.c', 'main_doublylinked_list.c', 'case1.txt'),
(3, 1, '基本課題1', 'sub1', 'open_addressing: open_addressing.o main_open_addressing.o', 'open_addressing.c', 'main_open_addressing.c', 'case1.txt'),
(4, 1, '基本課題1', 'sub1',"binarytree: binarytree.o main_binarytree.o", 'binarytree.c', 'main_binarytree.c', 'case1.txt'),
(4, 2, '基本課題2', 'sub2', 'binarysearchtree: binarysearchtree.o main_binarysearchtree.o', 'binarysearchtree.c', 'main_binarysearchtree.c', 'case1.txt'),
(4, 3, '発展課題1', 'sub3', 'binary_tree_mirror: binary_tree_mirror.o main_binary_tree_mirror.o', 'binary_tree_mirror.c', 'main_binary_tree_mirror.c', 'case1.txt'),
(4, 4, '発展課題2', 'sub4', 'bst_advanced: bst_advanced.o main_bst_advanced.o', 'bst_advanced.c', 'main_bst_advanced.c', 'case1.txt');

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    hashed_password VARCHAR(255) NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    disabled BOOLEAN DEFAULT FALSE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    active_start_date DATETIME NULL,
    active_end_date DATETIME NULL
);

CREATE TABLE IF NOT EXISTS  auth_codes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(255) NOT NULL,
    expired_at DATETIME NOT NULL,
    is_expired BOOLEAN DEFAULT FALSE,
    user_id INT NULL,
    CONSTRAINT fk_users_id FOREIGN KEY (user_id) 
        REFERENCES users(id)
);
