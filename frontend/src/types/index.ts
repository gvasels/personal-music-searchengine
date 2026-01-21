export interface Track {
  id: string;
  title: string;
  artist: string;
  album: string;
  albumId?: string;
  duration: number;
  trackNumber?: number;
  year?: number;
  genre?: string;
  format: string;
  fileSize: number;
  s3Key: string;
  artworkS3Key?: string;
  tags: string[];
  createdAt: string;
  updatedAt: string;
}

export interface Album {
  id: string;
  name: string;
  artist: string;
  year?: number;
  trackCount: number;
  coverArt?: string;
  createdAt: string;
}

export interface Artist {
  id: string;
  name: string;
  trackCount: number;
  albumCount: number;
}

export interface Playlist {
  id: string;
  name: string;
  description?: string;
  trackIds: string[];
  trackCount: number;
  coverArt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Tag {
  name: string;
  trackCount: number;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  limit: number;
  offset: number;
}

export type RepeatMode = 'off' | 'all' | 'one';
export type Theme = 'light' | 'dark';
