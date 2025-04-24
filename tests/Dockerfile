FROM python:3.13-slim

WORKDIR /tests

COPY . .
COPY req.txt .

RUN pip install --no-cache-dir -r req.txt

CMD ["pytest", "-s", "/tests/test.py", "-v"]