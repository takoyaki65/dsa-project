import type React from "react";
import NavigationBar from "../components/NavigationBar";
import { formatTimestamp } from "../util/timestamp";
import ResultBadge from "../components/ResultBadge";

const resultIDtoString = {
  0: "AC",
  1: "WA",
  2: "TLE",
  3: "MLE",
  4: "RE",
  5: "CE",
  6: "OLE",
  7: "IE",
  8: "FN",
  9: "Judging",
  10: "WJ",
}

const resultIDtoExplanation = {
  0: "Accepted",
  1: "Wrong Answer",
  2: "Time Limit Exceeded",
  3: "Memory Limit Exceeded",
  4: "Runtime Error",
  5: "Compilation Error",
  6: "Output Limit Exceeded",
  7: "Internal Error",
  8: "File Not Found",
  9: "Judging",
  10: "Waiting for Judging",
}

interface ValidationResult {
  id: number;
  ts: number;
  user_id: string;
  user_name: string;
  lecture_id: number;
  problem_id: number;
  result_id: number;
  time_ms: number;
  memory_kb: number;
}

interface ProblemInfo {
  lecture_id: number;
  problem_id: number;
  registered_at: number;
  title: string;
}

interface LectureInfo {
  lecture_id: number;
  title: string;
  start_date: number;
  deadline: number;
  problems: ProblemInfo[];
}

interface APIResponse {
  results: ValidationResult[];
  lecture_info: LectureInfo[];
}

const ValidationResultsListing: React.FC = () => {
  const sampleData: APIResponse = {
    results: [
      {
        id: 30,
        ts: 1758441667,
        user_id: "admin",
        user_name: "admin",
        lecture_id: 1,
        problem_id: 2,
        result_id: 5,
        time_ms: 163,
        memory_kb: 7680
      },
      {
        id: 29,
        ts: 1758441667,
        user_id: "admin",
        user_name: "admin",
        lecture_id: 1,
        problem_id: 1,
        result_id: 0,
        time_ms: 163,
        memory_kb: 7680
      },
    ],
    lecture_info: [
      {
        lecture_id: 1,
        title: "課題1",
        start_date: 1759280400,
        deadline: 1764550800,
        problems: [
          {
            lecture_id: 1,
            problem_id: 1,
            registered_at: 1758441326,
            title: "基本課題"
          },
          {
            lecture_id: 1,
            problem_id: 2,
            registered_at: 1758441345,
            title: "発展課題"
          }
        ]
      }
    ]
  };

  const getProblemTitle = (lectureId: number, problemId: number): string => {
    const lecture = sampleData.lecture_info.find(l => l.lecture_id === lectureId);
    if (!lecture) return "Unknown";
    const problem = lecture.problems.find(p => p.problem_id === problemId);
    if (!problem) return "Unknown";
    return `${lecture.title} · ${problem.title}`;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <NavigationBar />

      <div className="container mx-auto px-8 py-6">
        <h1 className="text-3xl font-semibold mb-6">Validation Results</h1>

        {/* Warning Message */}
        <div className="bg-white border-2 border-red-500 rounded-lg p-4 mb-6">
          <p className="text-sm">注） ここで行った提出で、課題の評価はされません。</p>
          <p className="text-sm">注） 問題無く採点可能であることを確認した後、manabaで提出してください。</p>
          <p className="text-sm">注） 提出してから一週間程度で結果は削除されます。</p>
        </div>

        {/* Pagination - Top */}
        <div className="flex justify-center mb-4 gap-4">
          <button className="text-gray-500 hover:text-gray-700">
            &lt; Prev
          </button>
          <button className="text-blue-500 hover:text-blue-700">
            Next &gt;
          </button>
        </div>

        {/* Table */}
        <div className="bg-white shadow-sm border border-gray-300">
          <table className="w-full">
            <thead>
              <tr className="border-b-2 border-gray-800">
                <th className="border-r border-gray-400 px-4 py-3 text-left font-semibold">提出日時</th>
                <th className="border-r border-gray-400 px-4 py-3 text-left font-semibold">問題</th>
                <th className="border-r border-gray-400 px-4 py-3 text-left font-semibold">ユーザ</th>
                <th className="border-r border-gray-400 px-4 py-3 text-left font-semibold">結果</th>
                <th className="border-r border-gray-400 px-4 py-3 text-left font-semibold">実行時間</th>
                <th className="border-r border-gray-400 px-4 py-3 text-left font-semibold">メモリ</th>
                <th className="px-4 py-3 text-center font-semibold"></th>
              </tr>
            </thead>
            <tbody>
              {sampleData.results.map((result, index) => (
                <tr key={`result-${result.id}`} className={`border-b border-gray-400 ${index % 2 === 0 ? "bg-white" : "bg-gray-50"}`}>
                  <td className="border-r border-gray-400 px-4 py-3 text-sm">
                    {formatTimestamp(result.ts)}
                  </td>
                  <td className="border-r border-gray-400 px-4 py-3 text-sm">
                    {getProblemTitle(result.lecture_id, result.problem_id)}
                  </td>
                  <td className="border-r border-gray-400 px-4 py-3 text-sm">
                    {result.user_name} ({result.user_id})
                  </td>
                  <td className="border-r border-gray-400 px-4 py-3 text-sm">
                    <div className="flex justify-center">
                      <ResultBadge resultID={result.result_id} />
                    </div>
                  </td>
                  <td className="border-r border-gray-400 px-4 py-3 text-sm">
                    {result.time_ms} ms
                  </td>
                  <td className="border-r border-gray-400 px-4 py-3 text-sm">
                    {result.memory_kb} KiB
                  </td>
                  <td className="px-4 py-3 text-center">
                    <button className="text-blue-500 hover:text-blue-700 underline text-sm">
                      詳細
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Pagination - Top */}
        <div className="flex justify-center mt-4 gap-4">
          <button className="text-gray-500 hover:text-gray-700">
            &lt; Prev
          </button>
          <button className="text-blue-500 hover:text-blue-700">
            Next &gt;
          </button>
        </div>
      </div>
    </div>
  )
}

export default ValidationResultsListing;