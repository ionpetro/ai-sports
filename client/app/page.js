'use client';
import { useState } from 'react';

export default function Home() {
  const [image, setImage] = useState(null);
  const [analysis, setAnalysis] = useState(null);
  const [loading, setLoading] = useState(false);
  const [context, setContext] = useState('');

  const handleImageUpload = async (e) => {
    const file = e.target.files[0];
    if (file && file.type.startsWith('image/')) {
      setImage(file);
      setAnalysis(null);
    } else {
      alert('Please upload a valid image file');
    }
  };

  const analyzeImage = async () => {
    if (!image) return;

    setLoading(true);
    try {
      const formData = new FormData();
      formData.append('image', image);
      formData.append('context-ai', context);

      const response = await fetch('http://localhost:8080/', {
        method: 'POST',
        body: formData,
      });

      const data = await response.json();
      console.log('Response data:', data);
      setAnalysis(data.api_response);
    } catch (error) {
      console.error('Error analyzing image:', error);
      alert('Error analyzing image. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-background p-8">
      <div className="max-w-2xl mx-auto space-y-8">
        <h1 className="text-4xl font-bold text-foreground text-center">
          Sports Image Analysis
        </h1>
        
        <div className="space-y-4">
          <div>
            <input
              type="file"
              accept="image/*"
              onChange={handleImageUpload}
              className="hidden"
              id="image-upload"
            />
            <label
              htmlFor="image-upload"
              className="cursor-pointer inline-block px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 transition-colors"
            >
              Select Image
            </label>
            {image && (
              <p className="mt-2 text-sm text-gray-600">
                Selected: {image.name}
              </p>
            )}
          </div>

          <div>
            <textarea
              value={context}
              onChange={(e) => setContext(e.target.value)}
              placeholder="Enter analysis context or instructions..."
              className="w-full p-2 border rounded-md dark:bg-gray-800 dark:border-gray-700"
              rows={4}
            />
          </div>

          <button
            onClick={analyzeImage}
            disabled={!image || loading}
            className={`w-full py-2 px-4 rounded-md text-white transition-colors ${
              !image || loading
                ? 'bg-gray-400 cursor-not-allowed'
                : 'bg-green-500 hover:bg-green-600'
            }`}
          >
            {loading ? 'Analyzing...' : 'Analyze Image'}
          </button>
        </div>

        {analysis && (
          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-lg">
            <h2 className="text-2xl font-semibold mb-4">Analysis Results</h2>
            <pre className="whitespace-pre-wrap">{analysis}</pre>
          </div>
        )}
      </div>
    </div>
  );
}
