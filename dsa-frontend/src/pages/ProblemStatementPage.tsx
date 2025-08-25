import NavigationBar from "../components/NavigationBar";

interface ProblemDetail {
  lecture_id: number;
  problem_id: number;
  title: string;
  description: string;
  time_ms: number;
  memory_mb: number;
  required_files: string[];
};

// url: /problem/:lectureid/:problemid
const ProblemStatementPage: React.FC = () => {
  /**
   * example: 
  {
    "lecture_id": 1,
    "problem_id": 1,
    "title": "基本課題",
    "description": "# 基本課題\n[課題リンク](https://www.coins.tsukuba.ac.jp/~amagasa/lecture/dsa-jikken/report1/#_4)\n\n教科書リスト1-4（p. 7）の「ユークリッドの互除法」に基づいたプログラム`gcd_euclid.c`および`main_euclid.c`を作成しなさい。\n\n## ファイル gcd_euclid.c\n```c\n#include <stdio.h>\n#include <stdlib.h>\n\n// Find the greatest common divisor of the two integers, n and m.\nint gcd_euclid(int n, int m) {\n\n    // 関数を完成させよ\n\n    return n;\n}\n```\n\n## ファイル main_euclid.c\n```c\n#include <stdio.h>\n#include <stdlib.h>\n\n// Find the greatest common divisor of the two integers, n and m.\nint gcd_euclid(int, int);\n\nint main(int argc, char *argv[]) {\n  // Process arguments.\n  if (argc != 3) {\n    fprintf(stderr, \"Usage: %s <number1> <number2>\\n\", argv[0]);\n    return EXIT_FAILURE;\n  }\n  int n = atoi(argv[1]);\n  int m = atoi(argv[2]);\n\n  // Compute and output the greatest common divisor.\n  int gcd = gcd_euclid(n, m);\n  printf(\"The GCD of %d and %d is %d.\\n\", n, m, gcd);\n\n  return EXIT_SUCCESS;\n}\n```\n\n# 提出方法\n`Makefile`, `gcd_euclid.c`, `main_euclid.c`の3点を提出せよ。\n* `Makefile` : 以下の内容が含まれたビルドスクリプト\n```Makefile\ngcd_euclid: gcd_euclid.o main_euclid.o\n```\n* `gcd_euclid.c` : 二つの整数から最大公約数を計算する関数`gcd_euclid`が定義されているCプログラム\n* `main_euclid.c` : `main`関数が定義されているCプログラム\n",
    "time_ms": 1000,
    "memory_mb": 512,
    "required_files": [
      "gcd_euclid.c",
      "main_euclid.c",
      "Makefile"
    ]
  }
  */

  return (
    <div>
      <NavigationBar />
      {/* Problem Statement... */}
    </div>
  )
}

export default ProblemStatementPage;
