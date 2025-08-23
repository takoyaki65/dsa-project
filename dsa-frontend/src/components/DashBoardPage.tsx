import { useState } from "react";

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

const DashBoardPage: React.FC = () => {
  // Static data
  const lectureData: Lecture[] = [
    {
      id: 1,
      title: "課題1",
      start_date: "2025-10-01T10:00:00+09:00",
      deadline: "2025-12-01T10:00:00+09:00",
      problems: [
        {
          lecture_id: 1,
          problem_id: 1,
          title: "基本課題 1"
        },
        {
          lecture_id: 1,
          problem_id: 2,
          title: "基本課題 2"
        },
        {
          lecture_id: 1,
          problem_id: 3,
          title: "応用課題"
        }
      ]
    },
    {
      id: 2,
      title: "課題2",
      start_date: "2025-10-01T10:00:00+09:00",
      deadline: "2025-12-01T10:00:00+09:00",
      problems: [
        {
          lecture_id: 2,
          problem_id: 1,
          title: "基本課題・連結リスト"
        },
        {
          lecture_id: 2,
          problem_id: 2,
          title: "発展課題・双連結リスト"
        }
      ]
    },
    {
      id: 3,
      title: "課題3",
      start_date: "2025-10-01T10:00:00+09:00",
      deadline: "2025-12-01T10:00:00+09:00",
      problems: []
    },
    {
      id: 4,
      title: "課題4",
      start_date: "2025-10-01T10:00:00+09:00",
      deadline: "2025-12-01T10:00:00+09:00",
      problems: []
    },
    {
      id: 5,
      title: "課題5",
      start_date: "2025-10-01T10:00:00+09:00",
      deadline: "2025-12-01T10:00:00+09:00",
      problems: []
    },
    {
      id: 6,
      title: "課題6",
      start_date: "2025-10-01T10:00:00+09:00",
      deadline: "2025-12-01T10:00:00+09:00",
      problems: []
    },
    {
      id: 7,
      title: "課題7",
      start_date: "2025-10-01T10:00:00+09:00",
      deadline: "2025-12-01T10:00:00+09:00",
      problems: []
    },
    {
      id: 8,
      title: "課題8",
      start_date: "2025-10-01T10:00:00+09:00",
      deadline: "2025-12-01T10:00:00+09:00",
      problems: []
    }
  ];

  const [selectedLecture, setSelectedLecture] = useState<number | null>(null);

  // Filter lectures based on selection
  const displayedLectures = selectedLecture === null
    ? lectureData
    : lectureData.filter(lecture => lecture.id === selectedLecture);

  const handleDSAClick = () => {
    // Navigate to main page (implementation to be added)
    console.log("Navigate to main page");
  };

  const handleLogout = () => {
    // Logout process (implementation to be added)
    console.log("Logout");
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Top Navigation Bar */}
      <div className="bg-blue-500 text-white px-6 py-4 flex justify-between items-center">
        <button
          key="dsa-button"
          onClick={handleDSAClick}
          className="text-2xl font-bold hover:opacity-80 transition-opacity"
        >
          DSA
        </button>
        <button
          key="logout-button"
          onClick={handleLogout}
          className="hover:bg-blue-600 px-4 py-2  rounded transition-colors"
        >
          Logout
        </button>
      </div>

      {/* Main Content */}
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
    </div>
  )
}

export default DashBoardPage;