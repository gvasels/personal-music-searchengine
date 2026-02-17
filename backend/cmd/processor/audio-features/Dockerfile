FROM public.ecr.aws/lambda/python:3.12

# Install system dependencies
RUN dnf install -y libsndfile && dnf clean all

# Use pre-built binary wheels
RUN pip install --upgrade pip && \
    pip install --no-cache-dir \
    numpy \
    scipy \
    librosa \
    soundfile \
    boto3

# Copy handler
COPY handler.py .

CMD ["handler.handler"]
