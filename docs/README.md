# オンラインジャッジWebアプリ 仕様書

## 1. 概要
プログラミング演習課題のオンラインジャッジWebアプリケーション

### 1.1 プロジェクト名
`dsa-project`

### 1.2 目的
プログラミング演習課題で提出されたプログラムコードのコンパイル・実行・テストケースの確認を
自動で行い、結果をWebアプリケーション上で確認する。

### 1.3 対象ユーザー
- システム管理者 (Admin)
  * 運用管理者の作成・削除
  * 1人のみ、初期設定時に作成される
- 運用管理者 (Manager)
  * 課題の作成・削除
  * 学生ユーザーの作成・削除
  * 提出された課題の詳細なジャッジ
    * 全てのテストケースを実行し、結果を表示
- 学生 (Student)
  * 提出ファイルのバリデーション
    * コンパイルが通るか、実行ができるか検証し、結果を表示

## 2. システム要件

### 2.1 機能要件
- ユーザー認証・認可システム
- コード提出システム
  - 提出されたコードをsandbox上でコンパイル・実行し、結果を表示
  - バリデーションリクエスト
    * 提出されたコードをsandbox上で限られたタスクを実行し、自動採点ができるかチェック
  - 採点リクエスト
    * 提出されたコードをsandbox上で全てのタスクを実行し、結果を表示
    * 運用管理者、システム管理者のみ
- 結果表示システム
  - コンパイル・実行・テストケースの確認結果を表示
- 管理者機能
  - 管理者ユーザーの作成・削除
  - 課題の作成・削除
  - 学生ユーザーの作成・削除
  - 学生が提出したファイルを一つにまとめたzipファイルをアップロードし、まとめてコンパイル・実行・テストケースの確認を行う
    - (高難易度) フォーマットが微妙に異なることでチェックができない提出に対して、その場で修正して再チェックすることができる

### 2.2 非機能要件
- セキュリティ
  - ログイン認証時に、ロール毎に異なる権限を設定
  - 時間が経過すると自動でログアウト
  - パスワードはハッシュ化して保存
  - sandbox上での任意のコード実行時のセキュリティ
    - CPUコア数、メモリ使用量の制限
    - 実行時間制限
    - フォルダ・ファイルの読み込み・書き込み制限
    - ネットワークアクセス制限
  - 監視・ログ収集
    - ダッシュボードでシステム負荷を監視
    - WARN以上のログをメールで通知
    - 高負荷時にメール・Slackメッセージで通知
- パフォーマンス
  - 同時アクセス対応
    - 必要に応じて、バックエンドサーバーのプロセス数を増やしてスケーリングさせる
  - sandbox環境のパフォーマンス
    - コンパイル・実行・テストケースの確認結果を高速に行える
  - データベースのパフォーマンス
    - 不必要なデータは定期的に削除する
      - ログイン履歴履歴は一週間単位で削除する
      - 学生のバリデーション結果は一週間単位で削除する
- 可用性
  - 24時間稼働
- 可搬性
  - 簡単にデプロイできる
    - ハイパラメータを設定する箇所が少ない、または一か所にまとまっている。
    - コマンド一つでデプロイできる。

## 3. システム構成

### 3.1 アーキテクチャ
- フロントエンド: React (Vite) + TypeScript
- バックエンド: Python (FastAPI)
- データベース: MySQL
- ジャッジサーバー: Python + Docker

## 4. データベース設計

### 4.1 主要テーブル
強調したカラムは必須項目。
- **Users**: ユーザ情報管理
  - **id**: ユーザーID (文字列)
  - **name**: ユーザー名 (文字列)
  - **hashed_password**: パスワードのハッシュ値 (文字列)
  - **role**: ユーザーの権限 (enum: 'admin', 'manager', 'student')
  - **disabled**: ユーザーが無効化されているかどうか (boolean)
  - email: メールアドレス (文字列)
- **LoginHistory**: ログイン履歴管理、及び強制ログアウト機能のためのテーブル
  - **id**: ログイン履歴ID (auto increment)
  - **user_id**: ユーザーID (文字列)
  - **login_at**: ログイン時刻 (datetime, 1s精度)
  - **logout_at**: ログアウト予定時刻 (datetime, 1s精度)
  - **is_revoked**: ログアウトされたかどうか (boolean)
- **Lecture**: 授業情報管理
  - **id**: 授業 (整数)
  - **title**: 授業タイトル (文字列)
  - **start_date**: 公開開始日時 (datetime, 1s精度)
  - **end_date**: 公開終了日時 (datetime, 1s精度)
  - **deadline**: 課題提出締め切り日時 (datetime, 1s精度)
- **Problem**: 課題情報管理
  - **lecture_id**: 授業ID (整数)
  - **problem_id**: 課題ID (整数)
    - 小課題の番号を表す、e.g., "課題3-1"の"1", "課題3-2"の"2"
    - 授業IDと課題IDの組み合わせで一意
  - **title**: 課題タイトル (文字列)
  - **resource_location_id**: 課題リソースファイルへのパス (FileLocation.id)
- **Submission**: ジャッジリクエスト情報管理
  - **id**: ジャッジリクエストID (auto increment)
  - **ts**: ジャッジリクエスト時刻 (datetime, 1s精度)
  - **user_id**: 採点対象のユーザーID (文字列)
  - **submission_ts**: 提出時刻 (datetime, 1s精度)
    - 提出時刻は、実際に課題がManaba等の媒体で提出された際の時刻
    - 採点リクエスト時に提出時刻が指定される
    - 採点リクエストでも無い場合は、提出時刻はジャッジリクエスト時刻と同一となる
  - **request_user_id**: ジャッジリクエストしたユーザーのID (文字列)
    - 管理者が学生の提出ファイルをジャッジする場合、提出者と採点対象が一致しない場合がある
  - **eval**: 課題採点リクエストかどうか, True/False
  - **lecture_id**: 授業ID (整数)
  - **problem_id**: 課題ID (整数)
  - **upload_location_id**: 提出ファイルへのパス (FileLocation.id)
  - **total_task**: 実行しなければならないTestCaseの数 (整数, デフォルト0)
  - **completed_task**: 現在実行完了しているTestCaseの数 (整数, デフォルト0)
  - **result**: 採点結果 (enum: 'WJ', 'Judging', 'AC', 'WA', 'TLE', 'MLE', 'RE', 'CE', 'OLE', 'IE', 'FN')
    - 種類:
      - 'WJ': Wait for Judge
      - 'Judging': Under Judging
      - 'AC': Accepted, all tasks have passed
      - 'WA': Wrong Answer, some judge tasks have wrong answer
      - 'TLE': Time Limit Exceeded, execution time exceeds the limit in some tasks
      - 'MLE': Memory Limit Exceeded, memory usage exceeds the limit in some tasks
      - 'RE': Runtime Error, runtime error occurs in some tasks
      - 'CE': Compile Error, compile error occurs in some tasks
      - 'OLE': Output Limit Exceeded, output exceeds the limit in some tasks
      - 'IE': Internal Error, internal error occurs in some tasks
      - 'FN': File Not Found, all tasks have aborted because some required file not found
    - デフォルトは'WJ'
    - 全順序: 'WJ' < 'Judging' < 'AC' < 'WA' < 'TLE' < 'MLE' < 'RE' < 'CE' < 'OLE' < 'IE' < 'FN'
    - 各タスクのジャッジ結果の内、最大値がストアされる
  - **log_location_id**: ジャッジログへのパス (FileLocation.id)
    - 各テストケースの実行結果が記録される
      - 実行結果 (AC～IE)、実行時間、消費メモリ、実行コマンド、標準入力、標準出力、標準エラー出力
    - その他、メッセージログ等も記録される。
  - **timeMS**: 全タスクの最大実行時間[ms]
  - **memoryKB**: 全タスクの最大消費メモリ[KB]
- **FileLocation**: アップロードされたファイルの管理
  - **id**: アップロードファイルID (auto increment)
  - **path**: アップロードファイルへのパス (文字列)
  - **ts**: アップロード日時 (datetime, 1s精度)

* 注意: テーブルの数が膨大になり、ORM wrapperの実装やスキーマの管理が非常に煩雑になるため、課題リソースデータに関する詳細 (各テストケースの入出力データ、テストケースの実行時間制限、メモリ制限等) は、データベースに保存しない。代わりに該当フォルダへのパスをDBに保存し、設定データはJSONファイルで管理する。
* 実装の簡潔さのために、課題情報が更新された場合、古い課題情報及びその課題に対するジャッジ結果・アップロードファイルは全て削除される。

## 9. 付録
### 9.1 用語集
- Judge: 提出されたソースコードのに対していくつかのタスクを行い、想定された出力をするか確認すること。タスクには以下の種類がある。
  - Build Task: ソースコードをコンパイルするタスク
  - Run Task: コンパイルされたプログラムを実行し、与えられた入力に対して想定された出力をするか確認するタスク
- Lecture (授業): 複数の課題を含む、授業の単位。具体例としては、「ハッシュ」「木構造」「グラフ」「動的計画法」等がある。
- Problem (課題): 授業内の課題。具体例としては、「必須課題1」「必須課題2」「応用課題」等がある。
