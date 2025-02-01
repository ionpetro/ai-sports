import { NextResponse } from 'next/server';

export async function POST(request) {
  try {
    const formData = await request.formData();
    const video = formData.get('video');

    if (!video) {
      return NextResponse.json(
        { error: 'No video file provided' },
        { status: 400 }
      );
    }

    // Create a new FormData object to send to the Go server
    const goServerFormData = new FormData();
    goServerFormData.append('video', video);

    // Forward the request to the Go server
    const response = await fetch('http://localhost:8080/analyze', {
      method: 'POST',
      body: goServerFormData,
    });

    const data = await response.json();
    return NextResponse.json(data);

  } catch (error) {
    console.error('Error processing video:', error);
    return NextResponse.json(
      { error: 'Error processing video' },
      { status: 500 }
    );
  }
} 