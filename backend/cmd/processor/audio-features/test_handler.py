"""
TDD Red Phase -- Failing tests for HLS segment sampling functionality.

These tests target three new functions/behaviors:
  1. list_hls_segments(user_id, track_id) -> dict grouped by quality
  2. download_hls_sample(segments_by_quality) -> (wav_path, quality)
  3. handler() integration: prefer HLS segments, fall back to original,
     and include 'source' field in the result.

All tests MUST FAIL because the implementation does not exist yet.

NOTE: boto3, librosa, numpy, and soundfile are Lambda-only dependencies
not installed locally. We mock them at the sys.modules level before
importing handler.py so that the *real* failure reason is the missing
functions, not a missing package.
"""

import os
import sys
import unittest
from unittest.mock import MagicMock, patch

# ---------------------------------------------------------------------------
# 1. Set environment variables the handler expects
# ---------------------------------------------------------------------------
os.environ['MEDIA_BUCKET'] = 'test-media-bucket'
os.environ['NUMBA_CACHE_DIR'] = '/tmp'

# ---------------------------------------------------------------------------
# 2. Pre-populate sys.modules with mocks for heavy/unavailable packages
#    so that `import boto3` etc. inside handler.py resolve without error.
# ---------------------------------------------------------------------------
mock_boto3 = MagicMock()
mock_librosa = MagicMock()
mock_numpy = MagicMock()
mock_soundfile = MagicMock()

sys.modules.setdefault('boto3', mock_boto3)
sys.modules.setdefault('librosa', mock_librosa)
sys.modules.setdefault('librosa.beat', MagicMock())
sys.modules.setdefault('librosa.feature', MagicMock())
sys.modules.setdefault('numpy', mock_numpy)
sys.modules.setdefault('soundfile', mock_soundfile)

# ---------------------------------------------------------------------------
# 3. Now we can safely import handler (existing code will load fine).
#    Importing list_hls_segments / download_hls_sample will fail with
#    ImportError because they don't exist yet -- that is the Red signal.
# ---------------------------------------------------------------------------
import handler as handler_module  # noqa: E402  (must come after sys.modules patching)


class TestListHlsSegmentsGroupsByQuality(unittest.TestCase):
    """list_hls_segments should group S3 objects by quality level and exclude non-.ts files."""

    def test_groups_by_quality_and_filters_non_ts(self):
        # This import MUST fail -- list_hls_segments does not exist yet.
        from handler import list_hls_segments

        mock_s3 = MagicMock()
        mock_paginator = MagicMock()
        mock_paginator.paginate.return_value = [
            {
                'Contents': [
                    {'Key': 'hls/user1/track1/96k_00001.ts'},
                    {'Key': 'hls/user1/track1/96k_00002.ts'},
                    {'Key': 'hls/user1/track1/192k_00001.ts'},
                    {'Key': 'hls/user1/track1/320k_00001.ts'},
                    {'Key': 'hls/user1/track1/master.m3u8'},
                ]
            }
        ]
        mock_s3.get_paginator.return_value = mock_paginator

        with patch.object(handler_module, 's3', mock_s3):
            result = list_hls_segments('user1', 'track1')

        # Should have exactly the three quality buckets
        self.assertIsInstance(result, dict)
        self.assertEqual(set(result.keys()), {'96k', '192k', '320k'})

        # Correct counts per quality
        self.assertEqual(len(result['96k']), 2)
        self.assertEqual(len(result['192k']), 1)
        self.assertEqual(len(result['320k']), 1)

        # master.m3u8 must NOT appear anywhere
        all_keys = [k for keys in result.values() for k in keys]
        for key in all_keys:
            self.assertFalse(key.endswith('.m3u8'), f'Non-TS file leaked through: {key}')


class TestListHlsSegmentsEmpty(unittest.TestCase):
    """list_hls_segments should return an empty dict when no objects exist at the prefix."""

    def test_returns_empty_dict_when_no_objects(self):
        from handler import list_hls_segments

        mock_s3 = MagicMock()
        mock_paginator = MagicMock()
        mock_paginator.paginate.return_value = [{}]
        mock_s3.get_paginator.return_value = mock_paginator

        with patch.object(handler_module, 's3', mock_s3):
            result = list_hls_segments('user1', 'track1')

        self.assertIsInstance(result, dict)
        self.assertEqual(len(result), 0)


class TestDownloadHlsSampleSelectsRandom(unittest.TestCase):
    """download_hls_sample should pick a random quality, sample segments, and return a wav path."""

    def test_returns_file_path_and_quality(self):
        from handler import download_hls_sample

        segments = {
            '320k': [
                'hls/user1/track1/320k_00001.ts',
                'hls/user1/track1/320k_00002.ts',
                'hls/user1/track1/320k_00003.ts',
            ],
            '96k': [
                'hls/user1/track1/96k_00001.ts',
            ],
        }

        with patch.object(handler_module, 's3', MagicMock()), \
             patch('handler.random') as mock_random:
            mock_random.choice.return_value = '320k'
            mock_random.randint.return_value = 3
            mock_random.sample.return_value = [
                'hls/user1/track1/320k_00001.ts',
                'hls/user1/track1/320k_00002.ts',
                'hls/user1/track1/320k_00003.ts',
            ]

            result = download_hls_sample(segments)

        # Should return a tuple of (file_path, quality_string)
        self.assertIsInstance(result, tuple)
        self.assertEqual(len(result), 2)

        file_path, quality = result
        self.assertIsInstance(file_path, str)
        self.assertTrue(file_path.endswith('.wav'), f'Expected .wav path, got: {file_path}')
        self.assertEqual(quality, '320k')


class TestHandlerUsesHlsWhenAvailable(unittest.TestCase):
    """When HLS segments exist, handler should use them instead of downloading the original file."""

    def test_prefers_hls_over_original(self):
        mock_s3 = MagicMock()

        with patch.object(handler_module, 's3', mock_s3), \
             patch.object(handler_module, 'list_hls_segments') as mock_list, \
             patch.object(handler_module, 'download_hls_sample') as mock_download, \
             patch.object(handler_module, 'librosa') as mock_lib, \
             patch.object(handler_module, 'detect_key', return_value=('A minor', '8A')):

            mock_list.return_value = {
                '192k': ['hls/u1/t1/192k_00001.ts', 'hls/u1/t1/192k_00002.ts'],
            }
            mock_download.return_value = ('/tmp/sample.wav', '192k')
            mock_lib.load.return_value = (MagicMock(), 22050)
            mock_lib.beat.beat_track.return_value = ([128.0], None)

            event = {
                'trackId': 't1',
                'userId': 'u1',
                's3Key': 'media/u1/t1.mp3',
            }

            result = handler_module.handler(event, None)

        # Must report HLS source
        self.assertIn('source', result)
        self.assertEqual(result['source'], 'hls')

        # Must NOT have downloaded the original file
        mock_s3.download_file.assert_not_called()

        # Should have called the HLS path
        mock_list.assert_called_once_with('u1', 't1')
        mock_download.assert_called_once()


class TestHandlerFallsBackToOriginal(unittest.TestCase):
    """When no HLS segments exist, handler should fall back to downloading the original file."""

    def test_uses_original_when_no_hls(self):
        mock_s3 = MagicMock()

        with patch.object(handler_module, 's3', mock_s3), \
             patch.object(handler_module, 'list_hls_segments') as mock_list, \
             patch.object(handler_module, 'librosa') as mock_lib, \
             patch.object(handler_module, 'detect_key', return_value=('C major', '8B')), \
             patch.object(handler_module, 'os') as mock_os:

            mock_list.return_value = {}
            mock_lib.load.return_value = (MagicMock(), 22050)
            mock_lib.beat.beat_track.return_value = ([120.0], None)

            event = {
                'trackId': 't1',
                'userId': 'u1',
                's3Key': 'media/u1/t1.mp3',
            }

            result = handler_module.handler(event, None)

        # Must report original source
        self.assertIn('source', result)
        self.assertEqual(result['source'], 'original')

        # Must have downloaded the original file
        mock_s3.download_file.assert_called_once()
        call_args = mock_s3.download_file.call_args
        self.assertEqual(call_args[0][0], 'test-media-bucket')
        self.assertEqual(call_args[0][1], 'media/u1/t1.mp3')


class TestHandlerReturnsSourceField(unittest.TestCase):
    """The handler result must always contain a 'source' field."""

    def test_result_contains_source_key(self):
        mock_s3 = MagicMock()

        with patch.object(handler_module, 's3', mock_s3), \
             patch.object(handler_module, 'librosa') as mock_lib, \
             patch.object(handler_module, 'detect_key', return_value=('G major', '9B')), \
             patch.object(handler_module, 'os') as mock_os:

            mock_lib.load.return_value = (MagicMock(), 22050)
            mock_lib.beat.beat_track.return_value = ([140.0], None)

            event = {
                'trackId': 't1',
                'userId': 'u1',
                's3Key': 'media/u1/t1.mp3',
            }

            result = handler_module.handler(event, None)

        self.assertIn('source', result,
                       "handler result must include 'source' field indicating 'hls' or 'original'")
        self.assertIn(result['source'], ('hls', 'original'),
                       "source must be either 'hls' or 'original'")


if __name__ == '__main__':
    unittest.main()
