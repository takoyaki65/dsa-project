-- データベースの作成
CREATE DATABASE IF NOT EXISTS dsa CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- データベースを使用
USE dsa;

-- テーブルの作成
CREATE TABLE IF NOT EXISTS assignments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    start_date DATETIME NOT NULL, 
    end_date DATETIME NOT NULL
);

-- 初期データの挿入
INSERT INTO assignments (title, start_date, end_date) VALUES
('課題1', '2023-10-01 00:00:00', '2025-12-31 23:59:59'),
('課題2', '2023-10-08 00:00:00', '2025-12-31 23:59:59'),
('課題3', '2024-10-30 00:00:00', '2025-12-31 23:59:59'), -- まだ開始していないケース  
('課題4', '2023-11-01 00:00:00', '2023-12-31 23:59:59'); -- もう終了しているケース  

CREATE TABLE IF NOT EXISTS sub_assignments (
    id INT NOT NULL,
    sub_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    makefile VARCHAR(255) NOT NULL,
    required_file_name VARCHAR(255) NOT NULL,
    test_file_name VARCHAR(255) NOT NULL,
    test_input_dir VARCHAR(255),
    test_output_dir VARCHAR(255),
    test_program_dir VARCHAR(255),
    test_case_name VARCHAR(255),
    test_program_name VARCHAR(255),
    PRIMARY KEY (id, sub_id),
    CONSTRAINT fk_assignments_id FOREIGN KEY (id) 
        REFERENCES assignments(id)
) ENGINE=InnoDB;

INSERT INTO sub_assignments (id, sub_id, title, makefile, required_file_name, test_file_name, test_input_dir, test_output_dir, test_program_dir, test_case_name, test_program_name) VALUES
(1, 1, '基本課題1', 'gcd_euclid: gcd_euclid.o main_euclid.o', 'gcd_euclid.c', 'main_euclid.c', '/app/dsa_test_case/report1/sub1/in', '/app/dsa_test_case/report1/sub1/out', '/app/dsa_test_program/report1/sub1', 'case1.txt', 'show/main_iter.c'),
(1, 2, '発展課題1', 'gcd_recursive: gcd_recursive.o main_recursive.o', 'gcd_recursive.c', 'main_recursive.c', NULL, NULL, NULL, NULL, NULL),
(2, 1, '基本課題1', 'linked_list: linked_list.o main_linked_list.o', 'linked_list.c', 'main_linked_list.c', NULL, NULL, NULL, NULL, NULL),
(2, 2, '基本課題2', 'queue: queue.o main_queue.o', 'queue.c', 'main_queue.c', NULL, NULL, NULL, NULL, NULL),
(2, 3, '発展課題1', 'doublylinked_list: doublylinked_list.o main_doublylinked_list.o', 'doublylinked_list.c', 'main_doublylinked_list.c', NULL, NULL, NULL, NULL, NULL),
(3, 1, '基本課題1', 'open_addressing: open_addressing.o main_open_addressing.o', 'open_addressing.c', 'main_open_addressing.c', NULL, NULL, NULL, NULL, NULL),
(4, 1, '基本課題1', "binarytree: binarytree.o main_binarytree.o", 'binarytree.c', 'main_binarytree.c', NULL, NULL, NULL, NULL, NULL),
(4, 2, '基本課題2', 'binarysearchtree: binarysearchtree.o main_binarysearchtree.o', 'binarysearchtree.c', 'main_binarysearchtree.c', NULL, NULL, NULL, NULL, NULL),
(4, 3, '発展課題1', 'binary_tree_mirror: binary_tree_mirror.o main_binary_tree_mirror.o', 'binary_tree_mirror.c', 'main_binary_tree_mirror.c', NULL, NULL, NULL, NULL, NULL),
(4, 4, '発展課題2', 'bst_advanced: bst_advanced.o main_bst_advanced.o', 'bst_advanced.c', 'main_bst_advanced.c', NULL, NULL, NULL, NULL, NULL);

