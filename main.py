import os
from network import post_to_discord
from logger import log_status, log_finish, log_success, log_error

def main():
    log_status("Starting Discord Watering Job...")
    
    webhook_url = os.environ.get("DISCORD_WEBHOOK")
    
    if not webhook_url:
        log_error("DISCORD_WEBHOOK environment variable not set.")
        return

    content = "รดน้ำกันจ้า <@650661678316388372> <@739506315784749097>"
    params = {
        "username": "รดน้ำ",
        "avatar_url": "https://raw.githubusercontent.com/tsphere101/watering-everyday/refs/heads/main/sanrio-hello-kitty-gogogal.jpg",
        "allowed_mentions": {
            "users": ["650661678316388372", "739506315784749097"]
        },
    }

    success, message = post_to_discord(webhook_url, content, params)
    
    if success:
        log_success(message)
    else:
        log_error(message)
        
    log_finish("Job complete.")

if __name__ == "__main__":
    main()
