# Repocheck

**Repocheck** is a Go based command-line interface (CLI) tool that analyzes a source code repository and reports on its **readiness, hygiene, and risk level**. It helps developers ensure their projects are consistent, secure, and ready for production.

**Version**: `v0.1.0`  
**Website**: [repocheck.co.in](https://repocheck.co.in)

---

## 🌟 Features

### Core (Free)
- **Repository Discovery**: Automatically identifies language, project type, and structure.
- **Onboarding Readiness**: Checks for `README.md`, setup instructions, and run commands.
- **Dependency Overview**: Counts dependencies and identifies package managers.
- **Hygiene Checks**: Scans for secrets, `.env.example`, and configuration consistency.
- **Configurable**: Ignore specific paths or rules via `.repocheck.yaml`.

### Pro (Paid)
- **Executive Summary**: High-level overview of repo health and top issues.
- **License Risk Analysis**: Detects high-risk licenses (e.g., GPL) and missing license files.
- **Readiness Score**: 0-100 "Credit Score" for your repo health with detailed breakdown.
- **Diff Mode**: Compare scans to track regressions.
- **CI Mode**: Strict exit codes and machine-readable output for pipelines.
- **Ownership**: Validate `CODEOWNERS` and maintainer definitions.
- **Extended Hygiene**: Checks for `.gitignore` and ignored environment files.

### UX (V2)
- **Colorized Output**: Semantic coloring for readability (Green used for Good, Red for Bad).
- **Auto-Detection**: Colors disable automatically in CI or when piped.

---

## Installation Guide
:white_check_mark: **Quick Verification**
After installing, run this command to make sure everything is working:

```bash
repocheck --version
```

### :window: Windows
**1. Download**
Download `repocheck_windows_amd64.exe` from the [Releases Page](https://github.com/bilal-sayyed/repocheck-go/releases).

**2. Open Terminal in Folder**
Go to your Downloads folder, hold `Shift + Right Click` on the background, and select "**Open PowerShell window here**".

**3. Run Repocheck**
Rename the downloaded file to `repocheck.exe` for convenience, then run:

```powershell
.\repocheck.exe scan .
```

**4. Global Access (Optional)**
To run it from anywhere, move `repocheck.exe` to a permanent folder (e.g., `C:\Tools`). Then search Windows for "**Edit the system environment variables**" -> Environment Variables -> Select Path -> Edit -> New -> Paste the folder path (`C:\Tools`).

### :computer: macOS & Linux
**1. Download & Unzip**
Download the binary (e.g., `repocheck_darwin_arm64` for Apple Silicon) from the [Releases Page](https://github.com/bilal-sayyed/repocheck-go/releases).

**2. Make Executable**
Open Terminal, navigate to downloads, and run:

```bash
chmod +x repocheck_darwin_arm64
```

**3. Install to PATH (Recommended)**
Move it to a global folder so you can run it from anywhere:

```bash
sudo mv repocheck_darwin_arm64 /usr/local/bin/repocheck
```
Now you can just type `repocheck scan .` in any project!

### :question: Troubleshooting
**Command not found?**
If verify fails, make sure you are in the same folder as the downloaded file (Windows) or that you successfully moved it to `/usr/local/bin` (Mac/Linux).

**MacOS Security Warning?**
If Mac says "**unidentified developer**", go to System Settings > Privacy & Security and click "**Allow Anyway**" for Repocheck.

---

## ⚡ Quick Start

```bash
# Initialize configuration (optional)
repocheck config init

# Scan the current directory
repocheck scan .
```

---

## 🔥 Free vs Pro

| Feature | Free Tier | Pro Tier Individual ($49/yr / ₹1999/yr), Team for 5 users ($199/yr / ₹7999/yr) |
| :--- | :---: | :---: |
| **Repository Discovery** | ✅ | ✅ |
| **Onboarding Checks** | ✅ | ✅ |
| **Dependency Count** | ✅ | ✅ |
| **Basic Hygiene (.env, secrets)** | ✅ | ✅ |
| **Readiness Score (0-100)** | ❌ | ✅ |
| **License Risk Analysis** | ❌ | ✅ |
| **CI/CD Integration (Exit Codes)** | ❌ | ✅ |
| **Regression Testing (`diff`)** | ❌ | ✅ |
| **Executive Summary** | ❌ | ✅ |

### 💎 Upgrade to Pro
Visit [repocheck.co.in](https://repocheck.co.in) to purchase a license key.
    ```bash
    repocheck license add RC1.eyJ...
    ```

---

## 🛠️ Command Usage

### 🔍 `repocheck scan`
Primary analysis command.
```bash
repocheck scan [path] [flags]
```
**Flags:**
- `--json`: Output results in JSON format (useful for piping).
- `--summary`: Show a brief summary only.
- `--ci`: Run in strict CI mode (Exit codes: 0=Pass, 1=Warn, 2=Fail).
- `--no-color`: Disable color output.

### 📜 `repocheck summary` (Pro)
Print a high-level executive summary of the repository status, including Top Issues and Recommendations.
```bash
repocheck summary
```

### ↔️ `repocheck diff` (Pro)
Compare two scan JSON reports to find regressions (e.g., Score dropped, New Secrets).
```bash
# 1. Generate baseline
repocheck scan . --json > baseline.json

# 2. Make changes...

# 3. Generate new scan
repocheck scan . --json > current.json

# 4. Compare
repocheck diff baseline.json current.json
```

### ⚙️ `repocheck config`
Manage local configuration.
```bash
repocheck config init   # Create .repocheck.yaml
repocheck config show   # Display current config
```

### 🔑 `repocheck license`
Manage your Pro license.
```bash
repocheck license status
repocheck license add <KEY>
repocheck license remove
```

---

## ⚙️ Configuration
Create a `.repocheck.yaml` file in your project root:

```yaml
strict_mode: true
ignore_paths:
  - "vendor/"
  - "node_modules/"
  - "dist/"
```

## 📸 Screenshots
### Repocheck Core Feature Scan Output

<img width="500" height="545" alt="Repocheck Core Feature Scan Output" src="https://github.com/user-attachments/assets/6c0fe636-9fbd-45d8-a2df-c4f6e1104c6b" />

### Repocheck Pro Feature Scan Output

<img width="842" height="951" alt="Repocheck Pro Feature Scan Output" src="https://github.com/user-attachments/assets/df1dd829-26b7-4758-a726-6c6962f085e8" />

### Repocheck Core Feature JSON Output

<img width="402" height="473" alt="Repocheck Core Feature JSON Output" src="https://github.com/user-attachments/assets/9da379b8-ad0c-42b1-b448-0e8fd8476fba" />

### Repocheck Pro Feature JSON Output

<img width="840" height="901" alt="Repocheck Pro Feature JSON Output" src="https://github.com/user-attachments/assets/38591e51-2583-4fe3-a775-c9ae3dc8b0ff" />


## 📄 License & Usage

Copyright (c) 2026 Bilal Sayyed. All rights reserved.

Repocheck source code is publicly visible, but usage is governed by the Repocheck License.

- The **Free tier** may be used for personal and evaluation purposes.
- **Pro features** require a valid paid license.
- **Commercial, CI, and team usage** of Pro features requires a paid license.
- A license key is required to unlock Pro functionality.

See the website for full license terms and pricing.

## 🔒 Privacy

**All scanning happens locally. No code leaves your machine.**

Repocheck does not upload source code, track usage, or collect telemetry.  
Your repositories are analyzed entirely on your system.

## � Contact Us

For support, feedback, or enterprise inquiries, please contact us at: **repochecks@gmail.com**
