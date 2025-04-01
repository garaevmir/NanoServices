FROM python:3.9-slim

RUN apt-get update && apt-get install -y \
    build-essential \
    curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY proto/post.proto ./proto/

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

RUN mkdir -p generated
RUN python -m grpc_tools.protoc \
    -I proto/ \
    --python_out=generated \
    --grpc_python_out=generated \
    proto/post.proto

COPY events_server.py .
COPY repository.py .

RUN sed -i 's/import post_pb2/from generated import post_pb2/' generated/post_pb2_grpc.py

ENV PYTHONPATH=/app

CMD ["python", "events_server.py"]