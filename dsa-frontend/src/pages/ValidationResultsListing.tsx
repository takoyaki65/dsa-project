import type React from "react";
import NavigationBar from "../components/NavigationBar";
import { formatTimestamp } from "../util/timestamp";
import ResultBadge from "../components/ResultBadge";
import { useAuthQuery } from "../auth/hooks";
import { useEffect, useState } from "react";

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
  const [anchor, setAnchor] = useState<number>(15000000);
  const [direction, setDirection] = useState<"next" | "prev">("next");
  const [currentData, setCurrentData] = useState<APIResponse | null>(null);

  // API call using useAuthQuery
  const { data, isLoading, error } = useAuthQuery<APIResponse>({
    queryKey: ['validationResults', `${anchor}`, `${direction}`],
    endpoint: `/problem/result/validation/list?anchor=${anchor}&direction=${direction}`,
    options: {
      queryOptions: {
        enabled: true,
        staleTime: 0, // always fetch fresh data
      }
    }
  })

  // Update currentData when new data is fetched
  useEffect(() => {
    if (data && data.results && data.results.length > 0) {
      setCurrentData(data);
    }
  }, [data]);

  const getProblemTitle = (lectureId: number, problemId: number): string => {
    if (!currentData) return "Unknown";
    const lecture = currentData.lecture_info.find(l => l.lecture_id === lectureId);
    if (!lecture) return "Unknown";
    const problem = lecture.problems.find(p => p.problem_id === problemId);
    if (!problem) return "Unknown";
    return `${lecture.title} · ${problem.title}`;
  }

  const handleNextPage = async () => {
    if (!currentData || !currentData.results || currentData.results.length === 0) return;

    // Get the smallest ID (oldest entry) from current data
    const minId = Math.min(...currentData.results.map(r => r.id));

    // Update pagination parameters
    setAnchor(minId);
    setDirection("next");

    // Refetch will be triggered automatically by useAuthQuery due to queryKey change
  };

  const handlePrevPage = async () => {
    if (!currentData || !currentData.results || currentData.results.length === 0) return;

    // Get the largest ID (newest entry) from current data
    const maxId = Math.max(...currentData.results.map(r => r.id));

    // Update pagination parameters
    setAnchor(maxId);
    setDirection("prev");

    // Refetch will be triggered automatically by useAuthQuery due to queryKey change
  };

  if (isLoading && !currentData) {
    return (
      <div className="min-h-screen bg-gray-50">
        <NavigationBar />
        <div className="container mx-auto px-8 py-6">
          <h1 className="text-3xl font-semibold mb-6">Validation Results</h1>
          <div className="flex justify-center items-center h-64">
            <div className="text-gray-600">Loading...</div>
          </div>
        </div>
      </div>
    )
  }

  if (error && !currentData) {
    return (
      <div className="min-h-screen bg-gray-50">
        <NavigationBar />
        <div className="container mx-auto px-8 py-6">
          <h1 className="text-3xl font-semibold mb-6">Validation Results</h1>
          <div className="flex justify-center items-center h-64">
            <div className="text-red-600">Error loading data.</div>
          </div>
        </div>
      </div>
    )
  }

  if (!currentData || !currentData.results || currentData.results.length === 0) {
    return (
      <div className="min-h-screen bg-gray-50">
        <NavigationBar />
        <div className="container mx-auto px-8 py-6">
          <h1 className="text-3xl font-semibold mb-6">Validation Results</h1>
          <div className="flex justify-center items-center h-64">
            <div className="text-gray-600">No validation results available.</div>
          </div>
        </div>
      </div>
    )
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
          <button
            onClick={handlePrevPage}
            className="text-gray-500 hover:text-gray-700"
            disabled={isLoading}>
            &lt; Prev
          </button>
          <button
            onClick={handleNextPage}
            className="text-blue-500 hover:text-blue-700"
            disabled={isLoading}>
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
              {currentData.results.map((result, index) => (
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
          <button
            onClick={handlePrevPage}
            className="text-gray-500 hover:text-gray-700"
            disabled={isLoading}>
            &lt; Prev
          </button>
          <button
            onClick={handleNextPage}
            className="text-blue-500 hover:text-blue-700"
            disabled={isLoading}>
            Next &gt;
          </button>
        </div>
      </div>
    </div>
  )
}

export default ValidationResultsListing;