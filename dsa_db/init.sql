-- データベースの作成
CREATE DATABASE IF NOT EXISTS dsa CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- データベースを使用
USE dsa;

-- Lectureテーブル(第1回授業, 第2回授業,...)の作成
CREATE TABLE IF NOT EXISTS Lecture (
    id INT AUTO_INCREMENT PRIMARY KEY, -- 授業エントリのID
    title VARCHAR(255) NOT NULL, -- 授業のタイトル名 e.g., 課題1, 課題2, ...
    start_date DATETIME NOT NULL, -- 課題ページの公開日
    end_date DATETIME NOT NULL -- 課題ページの公開終了日
);

-- Problemテーブル(課題1-1,1-2,2-1,...)の作成
CREATE TABLE IF NOT EXISTS Problem (
    lecture_id INT NOT NULL, -- Lecture.idからの外部キー
    assignment_id INT NOT NULL, -- 何番目の課題か, e.g., 1, 2, ...
    for_evaluation BOOLEAN NOT NULL, -- 課題採点用かどうか, True/False
    title VARCHAR(255) NOT NULL, -- 課題名 e.g., 基本課題1
    description_path VARCHAR(255) NOT NULL, -- 課題の説明文のファイルパス
    timeMS INT NOT NULL, -- ジャッジの制限時間[ms] e.g., 1000
    memoryMB INT NOT NULL, -- ジャッジの制限メモリ[MB] e.g., 1024
    PRIMARY KEY (lecture_id, assignment_id, for_evaluation),
    FOREIGN KEY (lecture_id) REFERENCES Lecture(id)
);

-- Executablesテーブル(実行ファイル名のリスト)の作成
CREATE TABLE IF NOT EXISTS Executables (
    id INT AUTO_INCREMENT PRIMARY KEY,
    lecture_id INT,
    assignment_id INT,
    for_evaluation BOOLEAN,
    name VARCHAR(255) NOT NULL, -- 実行ファイル名
    FOREIGN KEY (lecture_id, assignment_id, for_evaluation) REFERENCES Problem(lecture_id, assignment_id, for_evaluation)
);

-- ArrangedFilesテーブル(あらかじめこちらで用意したファイルリスト)の作成
CREATE TABLE IF NOT EXISTS ArrangedFiles (
    str_id VARCHAR(255) PRIMARY KEY, -- 文字列ID(ユニーク)
    lecture_id INT, -- 何回目の授業で出される課題か, e.g., 1, 2, ...
    assignment_id INT, -- 何番目の課題か, e.g., 1, 2, ...
    for_evaluation BOOLEAN, -- 課題採点用かどうか, True/False
    path VARCHAR(255) NOT NULL, -- ソースコードのパス(Makefileも全部含める)
    FOREIGN KEY (lecture_id, assignment_id, for_evaluation) REFERENCES Problem(lecture_id, assignment_id, for_evaluation)
);

-- RequiredFilesテーブル(ユーザに提出を求めれているファイルのリスト)の作成
CREATE TABLE IF NOT EXISTS RequiredFiles (
    id INT AUTO_INCREMENT PRIMARY KEY, -- ソースコードのID(auto increment)
    lecture_id INT, -- 何回目の授業で出される課題か, e.g., 1, 2, ...
    assignment_id INT, -- 何番目の課題か, e.g., 1, 2, ...
    for_evaluation BOOLEAN, -- 課題採点用かどうか, True/False
    name VARCHAR(255) NOT NULL, -- 提出が求められるファイルの名前
    FOREIGN KEY (lecture_id, assignment_id, for_evaluation) REFERENCES Problem(lecture_id, assignment_id, for_evaluation)
);


-- EvaluationItemsテーブル(課題に含まれる評価項目 e.g., func1, func2,...)の作成
CREATE TABLE EvaluationItems (
    str_id VARCHAR(255) PRIMARY KEY, -- 文字列ID(ユニーク)
    lecture_id INT, -- 何回目の授業で出される課題か, e.g., 1, 2, ...
    assignment_id INT, -- 何番目の課題か, e.g., 1, 2, ...
    for_evaluation BOOLEAN, -- 課題採点用かどうか, True/False
    title VARCHAR(255) NOT NULL, -- e.g., func1
    description TEXT, -- 説明
    score INT NOT NULL, -- 評価点
    type ENUM('Built', 'Judge') NOT NULL, -- 採点するタイミング
    arranged_file_id VARCHAR(255), -- 紐づいているソースコードのID, NULLABLE
    message_on_fail VARCHAR(255), -- 失敗した場合のメッセージ(一行、10文字程度)
    FOREIGN KEY (lecture_id, assignment_id, for_evaluation) REFERENCES Problem(lecture_id, assignment_id, for_evaluation),
    FOREIGN KEY (arranged_file_id) REFERENCES ArrangedFiles(str_id)
);


-- TestCasesテーブル(実行するテストのリスト)の作成
CREATE TABLE IF NOT EXISTS TestCases (
    id INT AUTO_INCREMENT PRIMARY KEY, -- テストケースのID(auto increment)
    eval_id VARCHAR(255) NOT NULL, -- 対応する評価項目のID
    description TEXT, -- 簡単な1行の説明
    command VARCHAR(255) NOT NULL, -- e.g., "./run.sh", "ls", ...
    argument_path VARCHAR(255), -- スクリプトもしくは実行バイナリに渡す引数が記されたファイルのパス
    stdin_path VARCHAR(255), -- 標準入力のパス, path/to/stdin.txt
    stdout_path VARCHAR(255), -- 想定される標準出力のパス, path/to/stdout.txt
    stderr_path VARCHAR(255), -- 想定される標準エラー出力のパス, path/to/stderr.txt
    exit_code INT NOT NULL DEFAULT 0, -- 想定される戻り値
    FOREIGN KEY (eval_id) REFERENCES EvaluationItems(str_id)
);


-- Users テーブル
CREATE TABLE Users (
    user_id VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    role ENUM('admin', 'manager', 'student') NOT NULL,
    disabled BOOLEAN DEFAULT false NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    active_start_date DATETIME NULL,
    active_end_date DATETIME NULL
);


-- LoginHistory テーブル
CREATE TABLE IF NOT EXISTS LoginHistory (
    user_id VARCHAR(255) NOT NULL,
    login_at DATETIME NOT NULL,
    logout_at DATETIME NOT NULL, -- ログアウト予定の時刻(リフレッシュトークンにより更新される予定あり)
    refresh_count INT DEFAULT 0,  -- リフレッシュした回数、回数制限つける
    PRIMARY KEY (user_id, login_at),
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE
);


-- BatchSubmissionテーブルの作成
CREATE TABLE IF NOT EXISTS BatchSubmission (
    id INT AUTO_INCREMENT PRIMARY KEY, -- バッチ採点のID(auto increment)
    ts DATETIME DEFAULT CURRENT_TIMESTAMP, -- バッチ採点のリクエスト時刻
    user_id VARCHAR(255), -- リクエストした管理者のID
    FOREIGN KEY (user_id) REFERENCES Users(user_id)
);


-- Submissionテーブルの作成
CREATE TABLE IF NOT EXISTS Submission (
    id INT AUTO_INCREMENT PRIMARY KEY, -- 提出されたジャッジリクエストのID(auto increment)
    ts DATETIME DEFAULT CURRENT_TIMESTAMP, -- リクエストされた時刻
    batch_id INT, -- ジャッジリクエストが属しているバッチリクエストのID, 学生のフォーマットチェック提出ならNULL
    user_id VARCHAR(255) NOT NULL, -- 採点対象のユーザのID
    lecture_id INT NOT NULL, -- 何回目の授業で出される課題か, e.g., 1, 2, ...
    assignment_id INT NOT NULL, -- 何番目の課題か, e.g., 1, 2, ...
    for_evaluation BOOLEAN NOT NULL, -- 課題採点用かどうか, True/False
    progress ENUM('pending', 'queued', 'running', 'done') DEFAULT 'pending', -- リクエストの処理状況, pending/queued/running/done
    total_task INT NOT NULL DEFAULT 0, -- 実行しなければならないTestCaseの数
    completed_task INT NOT NULL DEFAULT 0, -- 現在実行完了しているTestCaseの数
    FOREIGN KEY (batch_id) REFERENCES BatchSubmission(id),
    FOREIGN KEY (user_id) REFERENCES Users(user_id),
    FOREIGN KEY (lecture_id, assignment_id, for_evaluation) REFERENCES Problem(lecture_id, assignment_id, for_evaluation)
);


-- UploadedFilesテーブルの作成
CREATE TABLE IF NOT EXISTS UploadedFiles (
    id INT AUTO_INCREMENT PRIMARY KEY, -- アップロードされたファイルのID(auto increment)
    ts DATETIME DEFAULT CURRENT_TIMESTAMP, -- アップロードされた時刻
    submission_id INT, -- そのファイルが必要なジャッジリクエストのID
    path VARCHAR(255) NOT NULL, -- アップロードされたファイルのパス
    FOREIGN KEY (submission_id) REFERENCES Submission(id)
);


-- SubmissionSummary(一つの提出における、全体の採点結果)
CREATE TABLE IF NOT EXISTS SubmissionSummary (
    submission_id INT PRIMARY KEY, -- 対象のSubmissionリクエストのID
    batch_id INT, -- Submissionリクエストに紐づいたBatchリクエストのID
    user_id VARCHAR(255), -- 採点対象のユーザのID
    lecture_id INT NOT NULL, -- 何回目の授業で出される課題か, e.g., 1, 2, ...
    assignment_id INT NOT NULL, -- 何番目の課題か, e.g., 1, 2, ...
    for_evaluation BOOLEAN NOT NULL, -- 課題採点用かどうか, True/False
    /* Aggregation attributes over SubmissionSummary */
    result ENUM('AC', 'WA', 'TLE', 'MLE', 'RE', 'CE', 'OLE', 'IE', 'FN') NOT NULL, -- Submissionリクエスト全体の実行結果, FN(File Not Found)
    message VARCHAR(255), -- メッセージ(5文字～10文字程度)
    detail VARCHAR(255), -- 詳細(ファイルが足りない場合: "main.c func.c....", 実行ファイルが足りない場合: "main, func,...")
    score INT NOT NULL, -- 集計スコア (該当Submissionリクエストの全scoreの合計)
    timeMS INT DEFAULT 0, -- 実行時間[ms]
    memoryKB INT DEFAULT 0, -- 消費メモリ[KB]
    FOREIGN KEY (submission_id) REFERENCES Submission(id),
    FOREIGN KEY (batch_id) REFERENCES BatchSubmission(id),
    FOREIGN KEY (user_id) REFERENCES Users(user_id),
    FOREIGN KEY (lecture_id, assignment_id, for_evaluation) REFERENCES Problem(lecture_id, assignment_id, for_evaluation)
);


-- EvaluationSummary(一つの提出における、各評価項目の採点結果)
CREATE TABLE IF NOT EXISTS EvaluationSummary (
    id INT AUTO_INCREMENT PRIMARY KEY,
    parent_id INT NOT NULL, -- 対象のSubmissionリクエストのID
    batch_id INT, -- Submissionリクエストに紐づいたBatchリクエストのID
    user_id VARCHAR(255), -- 採点対象のユーザのID
    lecture_id INT NOT NULL, -- 何回目の授業で出される課題か, e.g., 1, 2, ...
    assignment_id INT NOT NULL, -- 何番目の課題か, e.g., 1, 2, ...
    for_evaluation BOOLEAN NOT NULL, -- 課題採点用かどうか, True/False
    eval_id VARCHAR(255) NOT NULL, -- 評価項目の文字列ID
    arranged_file_id VARCHAR(255), -- 紐づいているソースコードのID
    /* Aggregation attribltes over JudgeResult */
    result ENUM('AC', 'WA', 'TLE', 'MLE', 'RE', 'CE', 'OLE', 'IE') NOT NULL, -- 評価項目に含まれる全TestCaseの実行結果
    message VARCHAR(255), -- メッセージ(5文字～10文字程度)
    detail VARCHAR(255), -- 詳細 (ファイルが足りない場合: "main.c func.c....", 実行ファイルが足りない場合: "main, func,...")
    score INT NOT NULL, -- 集計結果 (ACの場合、EvaluationItems.scoreの値、それ以外は0点)
    timeMS INT DEFAULT 0, -- 実行時間[ms]
    memoryKB INT DEFAULT 0, -- 消費メモリ[KB]
    -- 以下、外部キー関係ではないけどEvaluationItemsやArrangedFilesから取ってくる値
    eval_title VARCHAR(255) NOT NULL, -- EvaluationItems.title
    eval_description TEXT, -- EvaluationItems.description
    eval_type ENUM('Built', 'Judge') NOT NULL, -- EvaluationItems.type
    arranged_file_path VARCHAR(255), -- ArrangedFiles.path
    -- 外部キー関係
    FOREIGN KEY (parent_id) REFERENCES SubmissionSummary(submission_id),
    FOREIGN KEY (batch_id) REFERENCES BatchSubmission(id),
    FOREIGN KEY (user_id) REFERENCES Users(user_id),
    FOREIGN KEY (lecture_id, assignment_id, for_evaluation) REFERENCES Problem(lecture_id, assignment_id, for_evaluation),
    FOREIGN KEY (eval_id) REFERENCES EvaluationItems(str_id),
    FOREIGN KEY (arranged_file_id) REFERENCES ArrangedFiles(str_id)
);


-- JudgeResultテーブルの作成
CREATE TABLE IF NOT EXISTS JudgeResult (
    id INT AUTO_INCREMENT PRIMARY KEY, -- ジャッジ結果のID(auto increment)
    ts DATETIME DEFAULT CURRENT_TIMESTAMP, -- ジャッジ結果が出た時刻
    parent_id INT NOT NULL, -- 親のEvaluationSummaryのID
    submission_id INT NOT NULL, -- ジャッジ結果に紐づいているジャッジリクエストのID
    testcase_id INT NOT NULL, -- ジャッジ結果に紐づいているテストケースのID
    result ENUM('AC', 'WA', 'TLE', 'MLE', 'RE', 'CE', 'OLE', 'IE') NOT NULL, -- 実行結果のステータス、 AC/WA/TLE/MLE/CE/RE/OLE/IE, 参考: https://atcoder.jp/contests/abc367/glossary
    timeMS INT NOT NULL, -- 実行時間[ms]
    memoryKB INT NOT NULL, -- 消費メモリ[KB]
    exit_code INT NOT NULL, -- 戻り値
    stdout TEXT NOT NULL, -- 標準出力
    stderr TEXT NOT NULL, -- 標準エラー出力
    -- 以降、外部キー関係ではないけどTestCasesから取ってくる値
    description TEXT, -- TestCases.description
    command TEXT NOT NULL, -- 実行したコマンド e.g., "./a.out 1 2 -loption..."
    stdin TEXT, -- 標準入力(実体)
    expected_stdout TEXT, -- 期待される標準出力
    expected_stderr TEXT, -- 期待される標準エラー出力
    expected_exit_code INT NOT NULL DEFAULT 0, -- 期待される戻り値
    FOREIGN KEY (parent_id) REFERENCES EvaluationSummary(id),
    FOREIGN KEY (submission_id) REFERENCES Submission(id),
    FOREIGN KEY (testcase_id) REFERENCES TestCases(id)
);


-- EvaluationResult(学生の提出に対する、各授業の採点結果)
CREATE TABLE IF NOT EXISTS EvaluationResult (
    id INT AUTO_INCREMENT PRIMARY KEY, -- 採点結果のID(auto increment)
    ts DATETIME DEFAULT CURRENT_TIMESTAMP, -- 採点結果が出た時刻
    user_id VARCHAR(255), -- 採点対象のユーザのID
    lecture_id INT NOT NULL, -- 何回目の授業で出される課題か, e.g., 1, 2, ...
    score INT, -- 集計スコア (該当Submissionリクエストの全scoreの合計)
    report_path VARCHAR(255), -- 採点結果のレポートのパス
    comment TEXT, -- コメント
    FOREIGN KEY (user_id) REFERENCES Users(user_id),
    FOREIGN KEY (lecture_id) REFERENCES Lecture(id)
);


-- 課題1のデータを挿入
INSERT INTO Lecture (title, start_date, end_date) VALUES
('課題1', '2023-10-01 00:00:00', '2025-12-31 23:59:59');

INSERT INTO Problem (lecture_id, assignment_id, for_evaluation, title, description_path, timeMS, memoryMB) VALUES
(1, 1, false, '基本課題', 'ex1-1/description.md', 1000, 1024),
(1, 2, false, '発展課題', 'ex1-2/description.md', 1000, 1024);

INSERT INTO Executables
(lecture_id, assignment_id, for_evaluation, name) VALUES
(1         , 1            , false         , 'gcd_euclid'),
(1         , 2            , false         , 'gcd_recursive');

INSERT INTO ArrangedFiles
(str_id         , lecture_id, assignment_id, for_evaluation, path) VALUES
('1-1-make'     , 1         , 1            , false         , 'ex1-1/Makefile'),
('1-1-testlink' , 1         , 1            , false         , 'ex1-1/test_link.c'),
('1-2-make'     , 1         , 2            , false         , 'ex1-2/Makefile'),
('1-2-testlink' , 1         , 2            , false         , 'ex1-2/test_link.c');

INSERT INTO RequiredFiles (lecture_id, assignment_id, for_evaluation, name) VALUES
(1, 1, false, 'gcd_euclid.c'),
(1, 1, false, 'main_euclid.c'),
(1, 1, false, 'Makefile'),
(1, 2, false, 'gcd_recursive.c'),
(1, 2, false, 'main_recursive.c'),
(1, 2, false, 'Makefile');

INSERT INTO EvaluationItems
(str_id         , lecture_id, assignment_id, for_evaluation, title         , description          , score, type        , arranged_file_id, message_on_fail      ) VALUES
('1-1-build'    , 1         , 1            , false         , 'compile'     , 'コンパイル'           , 0    , 'Built'     , '1-1-make'       , 'コンパイルに失敗しました'),
('1-1-check'    , 1         , 1            , false         , 'check'       , 'gcd_euclidの定義'    , 0    , 'Built'     , '1-1-make'       , 'gcd_euclidが定義されていません'),
('1-1-small'    , 1         , 1            , false         , 'smallNumber' , '小さい数同士のGCD'     , 0    , 'Judge'     , NULL             , '小さい数同士のGCDを求められていません'),
('1-1-invalid1' , 1         , 1            , false         , 'invalidArg'  , '引数が2つでない場合'    , 0    , 'Judge'     , NULL             , '引数が2つでない場合のエラー出力ができていません'),
('1-2-build'    , 1         , 2            , false         , 'compile'     , 'コンパイル'           , 0    , 'Built'     , '1-2-make'       , 'コンパイルに失敗しました'),
('1-2-check'    , 1         , 2            , false         , 'check'       , 'gcd_recursiveの定義' , 0    , 'Built'     , '1-2-make'       , 'gcd_recursiveが定義されていません'),
('1-2-small'    , 1         , 2            , false         , 'smallNumber' , '小さい数同士のGCD'     , 0    , 'Judge'     , NULL             , '小さい数同士のGCDを求められていません'),
('1-2-invalid1' , 1         , 2            , false         , 'invalidArg'  , '引数が2つでない場合'    , 0    , 'Judge'     , NULL             , '引数が2つでない場合のエラー出力ができていません');

INSERT INTO TestCases 
(eval_id             , description , command                      , argument_path                    , stdin_path, stdout_path                     , stderr_path                     , exit_code) VALUES
( '1-1-build'        , ''          , 'make gcd_euclid'            , NULL                             , NULL      , NULL                            , NULL                            , 0),
( '1-1-check'        , ''          , 'make test_link'             , NULL                             , NULL      , NULL                            , NULL                            , 0),
( '1-1-small'        , ''          , './gcd_euclid'               , 'ex1-1/testcases/easy1.arg'      , NULL      , 'ex1-1/testcases/easy1.out'     , 'ex1-1/testcases/easy1.err'     , 0),
( '1-1-small'        , ''          , './gcd_euclid'               , 'ex1-1/testcases/easy2.arg'      , NULL      , 'ex1-1/testcases/easy2.out'     , 'ex1-1/testcases/easy2.err'     , 0),
( '1-1-small'        , ''          , './gcd_euclid'               , 'ex1-1/testcases/easy3.arg'      , NULL      , 'ex1-1/testcases/easy3.out'     , 'ex1-1/testcases/easy3.err'     , 0),
( '1-1-small'        , ''          , './gcd_euclid'               , 'ex1-1/testcases/easy4.arg'      , NULL      , 'ex1-1/testcases/easy4.out'     , 'ex1-1/testcases/easy4.err'     , 0),
( '1-1-invalid1'     , ''          , './gcd_euclid'               , 'ex1-1/testcases/exception1.arg' , NULL      , 'ex1-1/testcases/exception1.out', 'ex1-1/testcases/exception1.err', 1),
( '1-2-build'        , ''          , 'make gcd_recursive'         , NULL                             , NULL      , NULL                            , NULL                            , 0),
( '1-2-check'        , ''          , 'make test_link'             , NULL                             , NULL      , NULL                            , NULL                            , 0),
( '1-2-small'        , ''          , './gcd_recursive'            , 'ex1-1/testcases/easy1.arg'      , NULL      , 'ex1-1/testcases/easy1.out'     , 'ex1-1/testcases/easy1.err'     , 0),
( '1-2-small'        , ''          , './gcd_recursive'            , 'ex1-1/testcases/easy2.arg'      , NULL      , 'ex1-1/testcases/easy2.out'     , 'ex1-1/testcases/easy2.err'     , 0),
( '1-2-small'        , ''          , './gcd_recursive'            , 'ex1-1/testcases/easy3.arg'      , NULL      , 'ex1-1/testcases/easy3.out'     , 'ex1-1/testcases/easy3.err'     , 0),
( '1-2-small'        , ''          , './gcd_recursive'            , 'ex1-1/testcases/easy4.arg'      , NULL      , 'ex1-1/testcases/easy4.out'     , 'ex1-1/testcases/easy4.err'     , 0),
( '1-2-invalid1'     , ''          , './gcd_recursive'            , 'ex1-2/testcases/exception.arg' , NULL      , 'ex1-2/testcases/exception.out', 'ex1-2/testcases/exception.err', 1);
