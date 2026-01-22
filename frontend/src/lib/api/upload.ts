/**
 * Upload API Module - Wave 4
 */
import { apiClient } from './client';

export interface PresignedUploadRequest {
  filename: string;
  contentType: string;
  fileSize: number;
}

export interface PresignedUploadResponse {
  uploadId: string;
  uploadUrl: string;
  expiresAt: string;
  maxFileSize: number;
}

export interface UploadConfirmResponse {
  status: 'processing' | 'completed' | 'failed';
  trackId: string | null;
}

export interface UploadStatusResponse {
  uploadId: string;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  trackId: string | null;
  error: string | null;
}

export async function getPresignedUploadUrl(data: PresignedUploadRequest): Promise<PresignedUploadResponse> {
  const response = await apiClient.post<PresignedUploadResponse>('/upload/presigned', data);
  return response.data;
}

export async function confirmUpload(uploadId: string): Promise<UploadConfirmResponse> {
  const response = await apiClient.post<UploadConfirmResponse>('/upload/confirm', { uploadId });
  return response.data;
}

export async function getUploadStatus(uploadId: string): Promise<UploadStatusResponse> {
  const response = await apiClient.get<UploadStatusResponse>(`/uploads/${uploadId}`);
  return response.data;
}
