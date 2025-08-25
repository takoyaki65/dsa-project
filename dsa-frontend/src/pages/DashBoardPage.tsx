import { useState } from "react";
import { useAuthQuery } from "../auth/hooks";
import NavigationBar from "../components/NavigationBar";

interface Problem {
  lecture_id: number;
  problem_id: number;
  title: string;
}

interface Lecture {
  id: number;
  title: string;
  start_date: string;
  deadline: string;
  problems: Problem[];
}

// url: /dashboard
const DashBoardPage: React.FC = () => {
  const [selectedLecture, setSelectedLecture] = useState<number | null>(null);

  const lectureDataQuery = useAuthQuery<Lecture[]>({
    queryKey: ['lectures'],
    endpoint: '/problem/fetch/list',
    options: {
      queryOptions: {
        retry: 2,
      }
    }
  });

  const isPending = lectureDataQuery.isPending;
  const lectureData = lectureDataQuery.data;
  const error = lectureDataQuery.error;

  const mainContent = () => {
    if (isPending) {
      return (
        <div className="min-h-screen flex items-center justify-center">
          <div className="text-gray-500">Loading...</div>
        </div>
      );
    }
    if (error) {
      return (
        <div className="min-h-screen flex items-center justify-center">
          <div className="text-red-500">Error loading lectures: {error.message}</div>
        </div>
      );
    }

    if (!lectureData) {
      return (
        <div className="min-h-screen flex items-center justify-center">
          <div className="text-gray-500">No lectures available</div>
        </div>
      );
    }

    // Filter lectures based on selection
    const displayedLectures = selectedLecture === null
      ? lectureData
      : lectureData.filter(lecture => lecture.id === selectedLecture);

    return (
      <div className="container mx-auto px-6 py-8">
        {/* Title */}
        <h1 className="text-2xl font-bold mb-6">Problem List</h1>

        {/* Tab Navigation */}
        <div className="flex gap-1 mb-8 border-b border-gray-300 overflow-x-auto">
          <button
            key="all-button"
            onClick={() => setSelectedLecture(null)}
            className={`px-4 py-2 whitespace-nowrap  transition-colors ${selectedLecture === null
              ? 'border-b-2 border-blue-500 text-blue-600 font-medium'
              : 'text-gray-600 hover:text-gray-800'
              }`}
          >
            All
          </button>
          {lectureData.map((lecture) => (
            <button
              key={lecture.id}
              onClick={() => setSelectedLecture(lecture.id)}
              className={`px-4 py-2 whitespace-nowrap transition-colors ${selectedLecture === lecture.id
                ? 'border-b-2 border-blue-500 text-blue-600 font-medium'
                : 'text-gray-600 hover:text-gray-800'
                }`}
            >
              {lecture.title}
            </button>
          ))}
        </div>

        {/* Lecture List */}
        <div className="space-y-6">
          {displayedLectures.map((lecture) => (
            <div key={lecture.id} className="bg-white rounded-lg shadow-sm border border-gray-200">
              <div className="px-6 py-4 border-b border-gray-200 bg-gray-50">
                <h2 className="text-lg font-semibold">{lecture.title}</h2>
              </div>
              <div className="p-6">
                {lecture.problems.length > 0 ? (
                  <div className="space-y-3">
                    {lecture.problems.map((problem) => (
                      <div
                        key={`${problem.lecture_id}-${problem.problem_id}`}
                        className="p-4 border border-gray-200 rounded hover:bg-gray-50 transition-colors cursor-pointer"
                      >
                        <span className="text-gray-700">{problem.title}</span>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-gray-400 text-center py-8">
                    No problems available
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      </div>
    )

  }

  return (
    <div className="min-h-screen bg-gray-50">
      <NavigationBar />
      {mainContent()}
    </div>
  )
}

export default DashBoardPage;