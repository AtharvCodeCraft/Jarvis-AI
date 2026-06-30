import os
import subprocess
import urllib.parse
import webbrowser

def open_application(app_name: str) -> str:
    app_lower = app_name.lower()
    app_map = {
        "chrome": "chrome.exe",
        "notepad": "notepad.exe",
        "calculator": "calc.exe",
        "calc": "calc.exe",
        "settings": "ms-settings:"
    }
    target = app_map.get(app_lower, f"{app_name}.exe")
    try:
        subprocess.Popen(target, shell=True)
        return f"Opening {app_name}."
    except Exception as e:
        return f"Failed to open {app_name}: {e}"

def google_search(query: str) -> str:
    url = f"https://www.google.com/search?q={urllib.parse.quote(query)}"
    webbrowser.open(url)
    return f"Searching Google for {query}."

def youtube_search(query: str) -> str:
    url = f"https://www.youtube.com/results?search_query={urllib.parse.quote(query)}"
    webbrowser.open(url)
    return f"Searching YouTube for {query}."

def open_website(url: str) -> str:
    if not url.startswith("http"):
        url = "https://" + url
    webbrowser.open(url)
    return f"Opening {url}."

def open_folder(folder_name: str) -> str:
    user_profile = os.environ.get("USERPROFILE", "")
    folder_map = {
        "downloads": os.path.join(user_profile, "Downloads"),
        "documents": os.path.join(user_profile, "Documents"),
        "pictures": os.path.join(user_profile, "Pictures"),
        "desktop": os.path.join(user_profile, "Desktop"),
        "music": os.path.join(user_profile, "Music"),
        "videos": os.path.join(user_profile, "Videos")
    }
    path = folder_map.get(folder_name.lower())
    if path and os.path.exists(path):
        os.startfile(path)
        return f"Opening {folder_name} folder."
    return f"Could not find folder {folder_name}."

def control_volume(action: str) -> str:
    if action == "up":
        # using pycaw or similar in production, fallback to powershell
        subprocess.run(["powershell", "-Command", "(New-Object -ComObject WScript.Shell).SendKeys([char]175)"], creationflags=subprocess.CREATE_NO_WINDOW)
        return "Increasing volume"
    elif action == "down":
        subprocess.run(["powershell", "-Command", "(New-Object -ComObject WScript.Shell).SendKeys([char]174)"], creationflags=subprocess.CREATE_NO_WINDOW)
        return "Decreasing volume"
    elif action == "mute":
        subprocess.run(["powershell", "-Command", "(New-Object -ComObject WScript.Shell).SendKeys([char]173)"], creationflags=subprocess.CREATE_NO_WINDOW)
        return "Toggling mute"
    return "Volume command not recognized"

def power_control(state: str) -> str:
    if state == "lock":
        subprocess.run(["rundll32.exe", "user32.dll,LockWorkStation"], creationflags=subprocess.CREATE_NO_WINDOW)
        return "Locking the screen"
    elif state == "sleep":
        subprocess.run(["rundll32.exe", "powrprof.dll,SetSuspendState", "0", "1", "0"], creationflags=subprocess.CREATE_NO_WINDOW)
        return "Putting the system to sleep"
    elif state == "restart":
        subprocess.run(["shutdown", "/r", "/t", "0"], creationflags=subprocess.CREATE_NO_WINDOW)
        return "Restarting the system"
    elif state == "shutdown":
        subprocess.run(["shutdown", "/s", "/t", "0"], creationflags=subprocess.CREATE_NO_WINDOW)
        return "Shutting down the system"
    return "Invalid power command"

def whatsapp_chat(contact_name: str) -> str:
    url = f"https://web.whatsapp.com/send?text=&phone=&text={urllib.parse.quote(contact_name)}"
    webbrowser.open(url)
    return f"Opening WhatsApp chat for {contact_name}."

def whatsapp_call(contact_name: str) -> str:
    url = f"https://web.whatsapp.com/send?text=&phone=&text={urllib.parse.quote(contact_name)}"
    webbrowser.open(url)
    return f"Opening WhatsApp call for {contact_name}."
