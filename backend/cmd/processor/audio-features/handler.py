"""
Audio Features Lambda - BPM and Key Detection using librosa
"""
import os
os.environ['NUMBA_CACHE_DIR'] = '/tmp'

import json
import boto3
import librosa
import numpy as np
from tempfile import NamedTemporaryFile

s3 = boto3.client('s3')
MEDIA_BUCKET = os.environ.get('MEDIA_BUCKET', '')

# Key detection mapping (Krumhansl-Schmuckler key profiles)
KEY_NAMES = ['C', 'C#', 'D', 'D#', 'E', 'F', 'F#', 'G', 'G#', 'A', 'A#', 'B']
CAMELOT_MAJOR = {'C': '8B', 'G': '9B', 'D': '10B', 'A': '11B', 'E': '12B', 'B': '1B', 
                 'F#': '2B', 'C#': '3B', 'G#': '4B', 'D#': '5B', 'A#': '6B', 'F': '7B'}
CAMELOT_MINOR = {'A': '8A', 'E': '9A', 'B': '10A', 'F#': '11A', 'C#': '12A', 'G#': '1A',
                 'D#': '2A', 'A#': '3A', 'F': '4A', 'C': '5A', 'G': '6A', 'D': '7A'}

def detect_key(y, sr):
    """Detect musical key using chroma features."""
    chroma = librosa.feature.chroma_cqt(y=y, sr=sr)
    chroma_avg = np.mean(chroma, axis=1)
    
    # Krumhansl-Schmuckler key profiles
    major_profile = np.array([6.35, 2.23, 3.48, 2.33, 4.38, 4.09, 2.52, 5.19, 2.39, 3.66, 2.29, 2.88])
    minor_profile = np.array([6.33, 2.68, 3.52, 5.38, 2.60, 3.53, 2.54, 4.75, 3.98, 2.69, 3.34, 3.17])
    
    major_corrs = [np.corrcoef(chroma_avg, np.roll(major_profile, i))[0, 1] for i in range(12)]
    minor_corrs = [np.corrcoef(chroma_avg, np.roll(minor_profile, i))[0, 1] for i in range(12)]
    
    major_key = np.argmax(major_corrs)
    minor_key = np.argmax(minor_corrs)
    
    if major_corrs[major_key] > minor_corrs[minor_key]:
        key_name = KEY_NAMES[major_key]
        mode = 'major'
        camelot = CAMELOT_MAJOR[key_name]
    else:
        key_name = KEY_NAMES[minor_key]
        mode = 'minor'
        camelot = CAMELOT_MINOR[key_name]
    
    return f"{key_name} {mode}", camelot

def handler(event, context):
    """Lambda handler for audio feature extraction."""
    track_id = event.get('trackId', '')
    user_id = event.get('userId', '')
    s3_key = event.get('s3Key', '')
    
    result = {'trackId': track_id, 'userId': user_id}
    
    if not MEDIA_BUCKET or not s3_key:
        result['error'] = 'missing configuration'
        return result
    
    try:
        # Download audio file
        with NamedTemporaryFile(suffix='.mp3', delete=False) as tmp:
            s3.download_file(MEDIA_BUCKET, s3_key, tmp.name)
            tmp_path = tmp.name
        
        # Load audio (30 sec sample for speed)
        y, sr = librosa.load(tmp_path, sr=22050, duration=60)
        
        # Detect BPM
        tempo, _ = librosa.beat.beat_track(y=y, sr=sr)
        bpm = round(float(tempo[0]) if hasattr(tempo, '__iter__') else float(tempo))
        
        # Detect key
        key, camelot = detect_key(y, sr)
        
        # Cleanup
        os.unlink(tmp_path)
        
        result['bpm'] = bpm
        result['key'] = key
        result['camelotCode'] = camelot
        
    except Exception as e:
        result['error'] = str(e)
    
    return result
