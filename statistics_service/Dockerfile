FROM python:3.9-slim

RUN apt-get update && apt-get install -y \
    build-essential \
    curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY proto/statistics.proto ./proto/

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

RUN mkdir -p generated
RUN python -m grpc_tools.protoc \
    -I proto/ \
    --python_out=generated \
    --grpc_python_out=generated \
    proto/statistics.proto

RUN sed -i 's/import statistics_pb2/from generated import statistics_pb2/' generated/statistics_pb2_grpc.py

COPY stats_server.py .
COPY repository.py .
COPY init.sql .

ENV PYTHONPATH=/app
ENV DATABASE_URL=postgresql+asyncpg://stats_user:stats_pass@stats_db:5432/stats_db
ENV KAFKA_BOOTSTRAP_SERVERS=kafka:9092

CMD ["python", "stats_server.py"]
