import urllib.request
import json
import urllib.error

data = json.dumps({
    "full_name": "Test User",
    "username": "testuser",
    "email": "test@example.com",
    "password": "Password123!"
}).encode('utf-8')

req = urllib.request.Request("http://localhost:8080/api/auth/register", data=data, headers={'Content-Type': 'application/json'})
try:
    with urllib.request.urlopen(req) as res:
        print("Status:", res.status)
        print("Response:", res.read().decode('utf-8'))
except urllib.error.HTTPError as e:
    print("Status:", e.code)
    print("Response:", e.read().decode('utf-8'))
