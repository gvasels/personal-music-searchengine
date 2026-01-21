# Upload Components - CLAUDE.md

## Overview

File upload components with drag-and-drop support for audio file uploads.

## Files

| File | Description |
|------|-------------|
| `UploadDropzone.tsx` | Drag-and-drop file upload area with progress |
| `index.ts` | Barrel export for upload components |
| `__tests__/UploadDropzone.test.tsx` | Unit tests (5 tests, 86.07% coverage) |

## Key Functions

### UploadDropzone.tsx
```typescript
interface UploadDropzoneProps {
  onFilesSelected: (files: File[]) => void;
  onError?: (error: string) => void;
  progress?: number;
  multiple?: boolean;
  maxSize?: number;        // Default: 500MB
  disabled?: boolean;
}

export function UploadDropzone(props: UploadDropzoneProps): JSX.Element
```
Renders drag-and-drop upload area with:
- Visual dropzone with instructions
- "Browse Files" button for manual selection
- Progress bar when uploading
- File type and size validation

**Supported Formats:**
- MP3 (audio/mpeg)
- FLAC (audio/flac)
- WAV (audio/wav)
- OGG (audio/ogg)
- M4A (audio/mp4)
- AAC (audio/aac)

**Validation:**
- File type must be audio/*
- File size must be under maxSize (default 500MB)
- Errors reported via `onError` callback

**Visual States:**
- Default: Dashed border
- Drag active: Primary color border with highlight
- Disabled: Reduced opacity, no interaction

## Dependencies

| Package | Usage |
|---------|-------|
| `react-dropzone` | `useDropzone` hook |

## Usage

```typescript
import { UploadDropzone } from '@/components/upload';

function UploadPage() {
  const [progress, setProgress] = useState(0);

  const handleFiles = async (files: File[]) => {
    // Upload logic with progress updates
  };

  return (
    <UploadDropzone
      onFilesSelected={handleFiles}
      onError={(msg) => toast.error(msg)}
      progress={progress}
    />
  );
}
```
