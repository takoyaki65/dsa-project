import NavigationBar from "../components/NavigationBar";
import Markdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { oneDark } from 'react-syntax-highlighter/dist/esm/styles/prism';

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
  const exampleData: ProblemDetail = {
    lecture_id: 1,
    problem_id: 1,
    title: "基本課題",
    description: "# 基本課題\n[課題リンク](https://www.coins.tsukuba.ac.jp/~amagasa/lecture/dsa-jikken/report1/#_4)\n\n教科書リスト1-4（p. 7）の「ユークリッドの互除法」に基づいたプログラム`gcd_euclid.c`および`main_euclid.c`を作成しなさい。\n\n## ファイル gcd_euclid.c\n```c\n#include <stdio.h>\n#include <stdlib.h>\n\n// Find the greatest common divisor of the two integers, n and m.\nint gcd_euclid(int n, int m) {\n\n    // 関数を完成させよ\n\n    return n;\n}\n```\n\n## ファイル main_euclid.c\n```c\n#include <stdio.h>\n#include <stdlib.h>\n\n// Find the greatest common divisor of the two integers, n and m.\nint gcd_euclid(int, int);\n\nint main(int argc, char *argv[]) {\n  // Process arguments.\n  if (argc != 3) {\n    fprintf(stderr, \"Usage: %s <number1> <number2>\\n\", argv[0]);\n    return EXIT_FAILURE;\n  }\n  int n = atoi(argv[1]);\n  int m = atoi(argv[2]);\n\n  // Compute and output the greatest common divisor.\n  int gcd = gcd_euclid(n, m);\n  printf(\"The GCD of %d and %d is %d.\\n\", n, m, gcd);\n\n  return EXIT_SUCCESS;\n}\n```\n\n# 提出方法\n`Makefile`, `gcd_euclid.c`, `main_euclid.c`の3点を提出せよ。\n* `Makefile` : 以下の内容が含まれたビルドスクリプト\n```Makefile\ngcd_euclid: gcd_euclid.o main_euclid.o\n```\n* `gcd_euclid.c` : 二つの整数から最大公約数を計算する関数`gcd_euclid`が定義されているCプログラム\n* `main_euclid.c` : `main`関数が定義されているCプログラム\n",
    time_ms: 1000,
    memory_mb: 512,
    required_files: [
      "gcd_euclid.c",
      "main_euclid.c",
      "Makefile"
    ]
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <NavigationBar />

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header Section */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 mb-6">
          <div className="px-6 py-4 border-b border-gray-200">
            <div className="flex items-center justify-between">

              {/* Meta information of problem */}
              <div className="flex items-center space-x-3">
                <span className="text-sm text-gray-500">
                  Lecture {exampleData.lecture_id}
                </span>
                <span className="text-sm text-gray-300">•</span>
                <span className="text-sm text-gray-500">
                  Problem {exampleData.problem_id}
                </span>
              </div>

              {/* Execution Time / Memory Limit */}
              <div className="flex items-center space-x-4 text-sm">
                <div className="flex items-center space-x-2">
                  <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <span className="text-gray-600">{exampleData.time_ms} ms</span>
                </div>
                <div className="flex items-center space-x-2">
                  <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2z" />
                  </svg>
                  <span className="text-gray-600">{exampleData.memory_mb} MB</span>
                </div>
              </div>

            </div>

            {/* Title */}
            <h1 className="text-2xl font-bold text-gray-900 mt-4">
              {exampleData.title}
            </h1>
          </div>
        </div>


        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* Main Content - Problem Description */}
          <div className="lg:col-span-3">
            <div className="bg-white rounded-lg shadow-sm border border-gray-200">
              <div className="px-6 py-4 border-b border-gray-200">
                <h2 className="text-lg font-semibold text-gray-900">問題文</h2>
              </div>

              <div className="p-6">
                {/* Markdown Content with Typography */}
                <div className="prose prose-sm max-w-none
                                prose-headings:text-gray-900
                                prose-h1:text-2xl prose-h1:font-bold prose-h1:mb-4 prose-h1:mt-6
                                prose-h2:text-xl prose-h2:font-semibold prose-h2:mb-3 prose-h2:mt-5
                                prose-h3:text-lg prose-h3:font-semibold prose-h3:mb-2 prose-h3:mt-4
                                prose-p:text-gray-700 prose-p:leading-7 prose-p:mb-4
                                prose-strong:text-gray-900
                                prose-pre:bg-[rgb(40,44,52)] prose-pre:overflow-x-auto prose-pre:rounded-lg prose-pre:p-4
                                prose-ul:list-disc prose-ul:pl-6 prose-ul:my-4
                                prose-ol:list-decimal prose-ol:pl-6 prose-ol:my-4
                                prose-li:text-gray-700 prose-li:my-1">
                  <Markdown
                    remarkPlugins={[remarkGfm]}
                    components={{
                      // Customization of code block
                      pre: ({ children }) => (
                        <div className="relative group">
                          <pre className="overflow-x-auto">{children}</pre>
                          <button
                            className="absolute top-2 right-2 bg-gray-700 hover:bg-gray-600 text-white px-2 py-1 rounded text-xs"
                            onClick={(e) => {
                              const code = e.currentTarget.parentElement?.querySelector('code')?.textContent;
                              if (code) navigator.clipboard.writeText(code);
                            }}
                          >Copy</button>
                        </div>
                      ),
                      // Customization of inline code
                      code: ({ children, className, node, ...rest }) => {
                        const match = /language-(\w+)/.exec(className || '');
                        const isCodeBlock = !!match;
                        const language = isCodeBlock ? match[1] : '';
                        if (isCodeBlock) {
                          // return <code className="px-0 py-0.5 rounded text-sm" {...rest}>{children}</code>
                          return (
                            <SyntaxHighlighter
                              // refを除くrest
                              className="text-base"
                              PreTag="pre"
                              children={String(children).replace(/\n$/, '')}
                              language={language}
                              style={oneDark}
                            />
                          )
                        }
                        return (
                          <code className="bg-gray-100 text-pink-600 px-1 py-0.5 rounded text-sm before:content-[''] after:content-['']" {...rest}>
                            {children}
                          </code>
                        );
                      },
                      // Customization of link
                      a: ({ children, href }) => (
                        <a
                          href={href}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-blue-600 no-underline hover:underline inline-flex items-center gap-1"
                        >
                          {children}
                          <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                          </svg>
                        </a>
                      ),
                    }}
                  >
                    {exampleData.description}
                  </Markdown>
                </div>
              </div>
            </div>
          </div>

          {/* Sidebar - Required Files & Submit */}
          <div className="lg:col-span-1 space-y-6">
            {/* Required Files */}
            <div className="bg-white rounded-lg shadow-sm border border-gray-200">
              <div className="px-4 py-3 border-b border-gray-200">
                <h3 className="text-sm font-semibold text-gray-900">Required Files</h3>
              </div>
              <div className="p-4">
                <ul className="space-y-2">
                  {exampleData.required_files.map((file, index) => (
                    <li key={index} className="flex items-center space-x-2">
                      <svg className="w-4 h-4 text-gray-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                      </svg>
                      <span>{file}</span>
                    </li>
                  ))}
                </ul>
              </div>
            </div>

          </div>
        </div>
      </div>


    </div>
  )
}

export default ProblemStatementPage;
