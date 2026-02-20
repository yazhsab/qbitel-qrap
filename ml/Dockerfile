FROM python:3.12-slim

WORKDIR /app

COPY ml/pyproject.toml .
COPY ml/src/ src/

RUN pip install --no-cache-dir .

EXPOSE 8084
CMD ["uvicorn", "qrap_ml.api.app:app", "--host", "0.0.0.0", "--port", "8084"]
