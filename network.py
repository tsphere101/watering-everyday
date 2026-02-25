import requests

def post_to_discord(webhook_url, content, params):
    if not content:
        return False, "content is empty."

    payloads = []
    # Discord has a 2000 char limit. Using 1900 to be safe.
    for i in range(0, len(content), 1900):
        chunk = content[i:i+1900]
        if i + 1900 < len(content):
            chunk += "... (see next message)"
        
        payloads.append({
            **params,
            "content": chunk
        })
    
    try:
        last_response = None
        for payload in payloads:
            last_response = requests.post(webhook_url, json=payload, timeout=10)
        
        if last_response and last_response.status_code != 204:
            return False, f"Discord Error {last_response.status_code}: {last_response.text}"
        
        return True, "Successfully posted to Discord!"
            
    except requests.exceptions.RequestException as e:
        return False, f"Request failed: {e}"
