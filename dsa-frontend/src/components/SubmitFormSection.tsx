import type React from "react";
import { useState } from "react";

interface UploadedFile {
  file: File;
  id: string;
}

interface SubmitFormSectionProps {
  onSubmit: (files: File[]) => Promise<void>;
  isLoading?: boolean;
}

const SubmitFormSection: React.FC<SubmitFormSectionProps> = ({ onSubmit, isLoading = false }) => {
  const [uploadedFiles, setUploadedFiles] = useState<UploadedFile[]>([]);
  const [isDragging, setIsDragging] = useState<boolean>(false);
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);

  const handleFileSelect = (files: FileList | null, inputElement?: HTMLInputElement) => {
    if (!files) return;

    const newFiles: UploadedFile[] = Array.from(files).map(file => ({
      file,
      id: `${file.name}-${Date.now()}-${Math.random()}`
    }));

    setUploadedFiles(prev => [...prev, ...newFiles]);

    // Reset input values to allow re-selecting the same file
    if (inputElement) {
      inputElement.value = '';
    }
  };

  // Handle drag events
  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);

    const files = e.dataTransfer.files;
    handleFileSelect(files);
  };

  // Remove a file
  const removeFile = (fileId: string) => {
    setUploadedFiles(prev => prev.filter(f => f.id !== fileId));
  };

  // Clear all files
  const clearAll = () => {
    setUploadedFiles([]);
  };

  // Handle form submission
  const handleSubmit = async () => {
    if (uploadedFiles.length === 0) return;

    setIsSubmitting(true);

    try {
      // Extract File objects from UploadedFile array
      const files = uploadedFiles.map(uf => uf.file);
      await onSubmit(files);

      // Clear files after successful submission
      clearAll();
    } catch (error) {
      console.error("Submission page error:", error);
    } finally {
      setIsSubmitting(false);
    };
  }

  // Calculate total size
  const totalSize = uploadedFiles.reduce((acc, uf) => acc + uf.file.size, 0);
  const totalSizeMB = (totalSize / (1024 * 1024)).toFixed(2);

  return (
    <>
      {/* Submit File Section */}
      <div className="bg-white rounded-lg shadow-sm border border-gray-200">
        <div className="px-4 py-3 border-b border-gray-200">
          <h3 className="text-sm font-semibold text-gray-900">Submit file</h3>
        </div>

        <div className="p-4 space-y-4">

          {/* Upload Area */}
          <div
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
            className={`relative border-2 border-dashed rounded-lg transition-all duration-200 ${isDragging
              ? 'border-blue-400 bg-blue-50'
              : 'border-gray-300 hover:border-gray-400 bg-white'
              }`}
          >
            <input
              type="file"
              id="file-upload"
              multiple
              onChange={(e) => handleFileSelect(e.target.files, e.target)}
              className="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10"
            //accept=".c,.h,.cpp,.hpp,.py,.java,.rs,.go,.js,.ts,.makefile,Makefile"
            />

            <div className="p-8 text-center">
              <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
              </svg>

              <p className="mt-3 text-sm font-medium text-gray-900">
                Click to upload or drag and drop
              </p>
              <p className="mt-1 text-xs text-gray-500">
                Multiple files supported
              </p>
            </div>
          </div>
        </div>

        {/* Uploaded Files List */}
        {uploadedFiles.length > 0 && (
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium text-gray-700">
                {uploadedFiles.length} file{uploadedFiles.length > 1 ? 's' : ''} • {totalSizeMB} MB
              </span>
              <button
                onClick={clearAll}
                className="text-sm text-red-500 hover:text-red-700 transition-colors"
              >
                Clear all
              </button>
            </div>

            <div className="max-h-48 overflow-y-auto space-y-1.5 pl-1 pr-1 group">
              {uploadedFiles.map((uploadedFile) => (
                <div
                  key={uploadedFile.id}
                  className="flex items-center justify-between p-2 bg-gray-50 rounded-lg border border-gray-200  group-hover:bg-gray-100 transition-colors"
                >
                  <div className="flex items-center space-x-2 flex-1 min-w-0">
                    <svg className="w-4 h-4 text-gray-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    <span className="text-sm text-gray-700 truncate">
                      {uploadedFile.file.name}
                    </span>
                    <span className="text-sm text-gray-400">
                      ({(uploadedFile.file.size / 1024).toFixed(1)} KB)
                    </span>
                  </div>

                  <button
                    onClick={() => removeFile(uploadedFile.id)}
                    className="ml-2"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>

                </div>
              ))}
            </div>
          </div>
        )}

        {/* Submit Button */}
        <button
          onClick={handleSubmit}
          disabled={uploadedFiles.length === 0 || isSubmitting || isLoading}
          className={`w-full py-2.5 px-4 rounded-lg font-medium transition-all duration-200 flex items-center justify-center space-x-2 ${uploadedFiles.length > 0 && !isSubmitting
            ? 'bg-blue-600 hover:bg-blue-700 text-white shadow-sm hover:shadow-md transform hover:-translate-y-0.5'
            : 'bg-gray-100 text-gray-400 cursor-not-allowed'
            }`}
        >
          {isSubmitting ? (
            <>
              <svg className="animate-spin h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span>Submitting...</span>
            </>
          )
            : (
              <>
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
                </svg>
                <span>Submit</span>
              </>
            )}
        </button>
      </div>
    </>
  )
}

export default SubmitFormSection;