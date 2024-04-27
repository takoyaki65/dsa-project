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
('課題4', 'report4', '2023-11-01 00:00:00', '2025-12-31 23:59:59'),
('課題5', 'report5', '2023-11-01 00:00:00', '2023-12-31 23:59:59'), -- もう終了しているケース
('課題6', 'report6', '2023-11-01 00:00:00', '2023-12-31 23:59:59'), -- もう終了しているケース
('課題7', 'report7', '2023-11-01 00:00:00', '2023-12-31 23:59:59'), -- もう終了しているケース
('課題8', 'report8', '2023-11-01 00:00:00', '2025-12-31 23:59:59');

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
(1, 1, '基本課題1', 'sub1', 'gcd_euclid: gcd_euclid.o main_gcd_euclid.o', 'gcd_euclid.c', 'main_gcd_euclid.c', 'test1.txt'),
(1, 2, '発展課題1', 'sub2', 'gcd_recursive: gcd_recursive.o main_recursive.o', 'gcd_recursive.c', 'main_recursive.c', 'test1.txt'),
(2, 1, '基本課題1', 'sub3', 'linked_list: linked_list.o main_linked_list.o', 'linked_list.c', 'main_linked_list.c', 'test1.txt'),
(2, 2, '基本課題2', 'sub1', 'queue: queue.o main_queue.o', 'queue.c', 'main_queue.c', 'test1.txt'),
(2, 3, '発展課題1', 'sub2', 'doublylinked_list: doublylinked_list.o main_doublylinked_list.o', 'doublylinked_list.c', 'main_doublylinked_list.c', 'test1.txt'),
(3, 1, '基本課題1', 'sub1', 'open_addressing: open_addressing.o main_open_addressing.o', 'open_addressing.c', 'main_open_addressing.c', 'test1.txt'),
(4, 1, '基本課題1', 'sub1',"binarytree: binarytree.o main_binarytree.o", 'binarytree.c', 'main_binarytree.c', 'test1.txt'),
(4, 2, '基本課題2', 'sub2', 'binarysearchtree: binarysearchtree.o main_binarysearchtree.o', 'binarysearchtree.c', 'main_binarysearchtree.c', 'test1.txt'),
(4, 3, '発展課題1', 'sub3', 'binarytree_mirror: binarytree_mirror.o main_binarytree_mirror.o', 'binarytree_mirror.c', 'main_binarytree_mirror.c', 'test1.txt'),
(4, 4, '発展課題2', 'sub4', 'bst_advanced: bst_advanced.o main_bst_advanced.o', 'bst_advanced.c', 'main_bst_advanced.c', 'test1.txt'),
(8, 1, "基本課題1", "sub1", "knapsack: knapsack.o main_knapsack.o", "knapsack.c", "main_knapsack.c", "test1.txt"),
(8, 2, "基本課題2", "sub2", "knapsackDP: knapsackDP.o main_knapsackDP.o", "knapsackDP.c", "main_knapsackDP.c", "test1.txt"),
(8, 3, "発展課題1", "sub3", "knapsack: knapsackDP2.o main_knapsackDP2.o", "knapsackDP2.c", "main_knapsackDP2.c", "test1.txt"),
(8, 4, "発展課題2", "sub4", "subsetsum: subsetsum.o main_subsetsum.o", "sumsetsum.c", "main_subsetsum.c", "test1.txt");

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

CREATE TABLE IF NOT EXISTS function_test (
    id INT NOT NULL,
    sub_id INT NOT NULL,
    func_id INT NOT NULL,
    func_name VARCHAR(255) NOT NULL,
    exec_command VARCHAR(255) ,
    PRIMARY KEY (id, sub_id, func_id),
    CONSTRAINT fk_sub_assignments_id FOREIGN KEY (id, sub_id) 
        REFERENCES sub_assignments(id, sub_id)
);

INSERT INTO function_test (id, sub_id, func_id, func_name, exec_command) VALUES
(1, 1, 1, 'gcd_euclid', './{unique_id}/bin/report1/sub1/gcd_euclid'),
(1, 1, 2, 'gcd_recursive', './{unique_id}/bin/report1/sub2/gcd_recursive'),
(2, 1, 1, 'insert_cell', './{unique_id}/bin/report2/sub1/incert_cell'),
(2, 1, 2, 'insert_cell_top', './{unique_id}/bin/report2/sub1/incert_cell_top'),
(2, 1, 3, 'delete_cell', './{unique_id}/bin/report2/sub1/delete_cell'),
(2, 1, 4, 'delete_cell_top', './{unique_id}/bin/report2/sub1/delete_cell_top'),
(2, 1, 5, 'display', './{unique_id}/bin/report2/sub1/display'),
(4, 1, 1, 'create_tree', './{unique_id}/bin/report4/sub1/create_tree'),
(4, 1, 2, 'preorder', './{unique_id}/bin/report4/sub1/preorder'),
(4, 1, 3, 'inorder', './{unique_id}/bin/report4/sub1/inorder'),
(4, 1, 4, 'postorder', './{unique_id}/bin/report4/sub1/postorder'),
(4, 1, 5, 'display', './{unique_id}/bin/report4/sub1/display'),
(4, 1, 6, 'breadth_first_search', './{unique_id}/bin/report4/sub1/breadth_first_search'),
(4, 1, 7, 'height', './{unique_id}/bin/report4/sub1/height'),
(4, 1, 8, 'delete_tree', './{unique_id}/bin/report4/sub1/delete_tree'),
(4, 2, 1, 'min_bst', './{unique_id}/bin/report4/sub2/min_bst'),
(4, 2, 2, 'search_bst', './{unique_id}/bin/report4/sub2/search_bst'),
(4, 2, 3, 'insert_bst', './{unique_id}/bin/report4/sub2/insert_bst'),
(4, 2, 4, 'delete_bst', './{unique_id}/bin/report4/sub2/delete_bst'),
(4, 3, 1, 'create_mirror', './{unique_id}/bin/report4/sub3/create_mirror'),
(4, 3, 2, 'are_mirrors', './{unique_id}/bin/report4/sub3/are_mirrors'),
(4, 4, 1, 'insert_bst', './{unique_id}/bin/report4/sub4/insert_bst'),
(4, 4, 2, 'delete_bst', './{unique_id}/bin/report4/sub4/delete_bst'),
(8, 1, 1, 'knapsack', './{unique_id}/bin/report8/sub1/knapsack'),
(8, 2, 1, 'knapsackDP', './{unique_id}/bin/report8/sub2/knapsackDP'),
(8, 3, 1, 'knapsackDP2', './{unique_id}/bin/report8/sub3/knapsackDP2'),
(8, 4, 1, 'subsetsum', './{unique_id}/bin/report8/sub4/subsetsum');
