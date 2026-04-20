# 📧 mailc - Simple CLI Email Sender

mailc is a lightweight command-line email tool written in Go. It allows you to send emails directly from the terminal with SMTP support, attachments, file-based content, and recipient lists.

---

# 🚀 Features

- Simple CLI interface
- SMTP authentication support
- TOML configuration file support
- Attachments (MIME Base64 encoded)
- Multiple recipients support
- File-based body input
- Lightweight single binary

---

# ⚙️ Configuration

Default config location:

- `~/.config/mailc/config.toml`
- or `/etc/mailc.conf` (root usage only, optional)

Example:

```toml
username = "your@email.com"
password = "your_app_password"
from = "your@email.com"

protocol = "smtp"

smtp_host = "smtp.gmail.com"
smtp_port = 587
use_tls = true

pop3_host = ""
pop3_port = 0

imap_host = ""
imap_port = 0

log_file = "~/.local/share/mailc/mailc.log"
debug = false
timeout = 10
retry = 3
```
# Basic usage
```toml
mailc [subject] [content] [recipient]
mailc "Hi" "Hello world" user@example.com
-s string    Email subject
-b string    Email body (text)
-f file      Email body from file
-r string    Recipient email
-l file      Recipient list file (one per line)
-a file      Attachment (can be used multiple times)
```
