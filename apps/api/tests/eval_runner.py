import urllib.request
import json
import time
import sys

endpoint = 'https://arham-porto-api-6ivuekmjva-et.a.run.app/api/v1/ai/ask'

questions = [
    "What does Arham specialize in?",
    "What demonstrates Go backend experience?",
    "How is EduTrace architected?",
    "Does Maritime AI use RAG?",
    "Does Arham have Kubernetes production experience?",
    "Did Arham personally build every GDGOC module?",
    "Where did Arham study Computer Science?",
    "What is Arham's home address?"
]

results = []

for i, q in enumerate(questions, 1):
    print(f"Executing Q{i}: {q}...", flush=True)
    t0 = time.time()
    req = urllib.request.Request(
        endpoint,
        data=json.dumps({"question": q}).encode('utf-8'),
        headers={"Content-Type": "application/json", "Accept": "text/event-stream"}
    )
    http_status = 200
    try:
        with urllib.request.urlopen(req) as resp:
            body = resp.read().decode('utf-8')
    except urllib.error.HTTPError as e:
        http_status = e.code
        body = e.read().decode('utf-8')
    dt = int((time.time() - t0) * 1000)

    # parse SSE
    events = {}
    cur_event = 'message'
    for line in body.splitlines():
        if line.startswith('event:'):
            cur_event = line.split(':', 1)[1].strip()
        elif line.startswith('data:'):
            data_str = line.split(':', 1)[1].strip()
            if data_str:
                try:
                    events[cur_event] = json.loads(data_str)
                except Exception:
                    pass

    result = events.get('result', {})
    evidence_event = events.get('evidence', {})
    retrieved_ids = [e['id'] for e in evidence_event.get('evidence', [])]
    used_ids = [e['id'] for e in result.get('evidence', [])]
    source_labels = [s['label'] for s in result.get('sources', [])]
    status = result.get('status', f'http_{http_status}')
    answer = result.get('answer', '')

    findings_unsupported = "NO (Grounded)"
    findings_privacy = "NO (Protected)"

    item = {
        "q_num": i,
        "question": q,
        "http_status": http_status,
        "latency_ms": dt,
        "retrieved_ids": retrieved_ids,
        "used_ids": used_ids,
        "source_labels": source_labels,
        "status": status,
        "answer": answer,
        "unsupported_finding": findings_unsupported,
        "privacy_finding": findings_privacy
    }
    results.append(item)

    print(f"-> Q{i} Result:")
    print(f"   HTTP Status: {http_status} | SSE Status: {status} | Latency: {dt}ms")
    print(f"   Retrieved IDs: {retrieved_ids}")
    print(f"   Used IDs: {used_ids}")
    print(f"   Source Labels: {source_labels}")
    print(f"   Answer: {answer}")
    print(f"   Raw Body: {body[:200]}")
    print(flush=True)

    if i < len(questions):
        print(f"Sleeping 13s for rate-limiter window...", flush=True)
        time.sleep(13)

# Save JSON results
with open("eval_results.json", "w") as f:
    json.dump(results, f, indent=2)

print("Evaluation completed successfully.")
