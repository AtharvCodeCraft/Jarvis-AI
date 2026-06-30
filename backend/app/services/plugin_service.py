import logging
from app.services.tool_registry import global_registry, ToolDefinition, ToolParameter
from app.services import system_controls

logger = logging.getLogger(__name__)

def load_all_plugins():
    logger.info("Loading system tools and plugins...")

    global_registry.register(
        ToolDefinition(
            name="open_application",
            description="Launch a Windows desktop application (e.g. chrome, notepad, calculator, settings, etc.)",
            parameters={
                "app_name": ToolParameter(type="string", description="The name of the application to launch (e.g. 'chrome', 'notepad', 'calc')", required=True)
            },
            requires_hitl=False
        ),
        lambda app_name: system_controls.open_application(app_name)
    )

    global_registry.register(
        ToolDefinition(
            name="google_search",
            description="Search the web using Google",
            parameters={
                "query": ToolParameter(type="string", description="The search terms or query to run", required=True)
            },
            requires_hitl=False
        ),
        lambda query: system_controls.google_search(query)
    )

    global_registry.register(
        ToolDefinition(
            name="youtube_search",
            description="Search for a video or music track on YouTube",
            parameters={
                "query": ToolParameter(type="string", description="The video name, channel name, or search query", required=True)
            },
            requires_hitl=False
        ),
        lambda query: system_controls.youtube_search(query)
    )

    global_registry.register(
        ToolDefinition(
            name="open_website",
            description="Navigate to a website URL in the default browser",
            parameters={
                "url": ToolParameter(type="string", description="The URL of the website to open (e.g. 'github.com', 'gmail.com')", required=True)
            },
            requires_hitl=False
        ),
        lambda url: system_controls.open_website(url)
    )

    global_registry.register(
        ToolDefinition(
            name="open_folder",
            description="Open a common Windows user folder path (downloads, documents, pictures, desktop, music, videos)",
            parameters={
                "folder_name": ToolParameter(
                    type="string", 
                    description="The target folder name (downloads, documents, pictures, desktop, music, videos)", 
                    enum=["downloads", "documents", "pictures", "desktop", "music", "videos"], 
                    required=True
                )
            },
            requires_hitl=False
        ),
        lambda folder_name: system_controls.open_folder(folder_name)
    )

    global_registry.register(
        ToolDefinition(
            name="control_volume",
            description="Adjust or mute system audio volume",
            parameters={
                "action": ToolParameter(
                    type="string", 
                    description="The action to execute: 'up' (increase), 'down' (decrease), or 'mute' (toggle mute)", 
                    enum=["up", "down", "mute"], 
                    required=True
                )
            },
            requires_hitl=False
        ),
        lambda action: system_controls.control_volume(action)
    )

    global_registry.register(
        ToolDefinition(
            name="power_control",
            description="Trigger OS power states (lock, sleep, restart, shutdown). Warning: Destructive actions require user prompt confirmation.",
            parameters={
                "state": ToolParameter(
                    type="string", 
                    description="The desired system power state: 'lock', 'sleep', 'restart', or 'shutdown'", 
                    enum=["lock", "sleep", "restart", "shutdown"], 
                    required=True
                )
            },
            requires_hitl=True
        ),
        lambda state: system_controls.power_control(state)
    )

    global_registry.register(
        ToolDefinition(
            name="whatsapp_action",
            description="Initiate a chat or call with a contact on WhatsApp Web",
            parameters={
                "contact_name": ToolParameter(type="string", description="The name of the contact/group to call or chat with", required=True),
                "action": ToolParameter(type="string", description="The WhatsApp action: 'chat' or 'call'", enum=["chat", "call"], required=True)
            },
            requires_hitl=False
        ),
        lambda contact_name, action: system_controls.whatsapp_call(contact_name) if action == "call" else system_controls.whatsapp_chat(contact_name)
    )
