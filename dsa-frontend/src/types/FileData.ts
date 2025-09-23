interface CompressedFileData {
  filename: string;
  data: string; // base64 encoded string, compressed with gzip
  compression: string; // e.g., "gzip"
  original_size: number; // original size in bytes
}

interface FileData {
  filename: string;
  data: Blob; // decompressed file data as Blob
  original_size: number; // original size in bytes
}

async function decompressFileData(compressedFileData: CompressedFileData): Promise<FileData> {
  if (compressedFileData.compression !== 'gzip') {
    console.warn(`Unknown compression type: ${compressedFileData.compression}`);
    // decoding base64 to binary only,
    const binaryString = atob(compressedFileData.data);
    const bytes = new Uint8Array(binaryString.length);
    for (let i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i);
    }
    return {
      filename: compressedFileData.filename,
      data: new Blob([bytes]),
      original_size: bytes.length,
    };
  }

  try {
    /* Base64 decode */
    const binaryString = atob(compressedFileData.data);
    const bytes = new Uint8Array(binaryString.length);
    for (let i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i);
    }

    /* Decompress */
    const decompressedData = await decompressWithStream(bytes);

    /* Create Blob from decompressed data */
    const blob = createBlobFromUint8Array(decompressedData);

    return {
      filename: compressedFileData.filename,
      data: blob,
      original_size: decompressedData.length,
    };
  } catch (error) {
    console.error('Error decompressing file data:', error);
    throw error;
  }
}

function createBlobFromUint8Array(data: Uint8Array): Blob {
  const buffer = new ArrayBuffer(data.length);
  const view = new Uint8Array(buffer);
  view.set(data);
  return new Blob([buffer]);
}

async function decompressWithStream(compressedData: Uint8Array): Promise<Uint8Array> {
  // Check for DecompressionStream support
  if (!('DecompressionStream' in window)) {
    throw new Error('DecompressionStream is not supported in this browser.');
  }

  const ds = new DecompressionStream('gzip');

  const blob = createBlobFromUint8Array(compressedData);
  const compressedStream = blob.stream();

  const decompressedstream = compressedStream.pipeThrough(ds);
  const decompressedArrayBuffer = await new Response(decompressedstream).arrayBuffer();
  return new Uint8Array(decompressedArrayBuffer);
}

export type { CompressedFileData, FileData };
export { decompressFileData };
