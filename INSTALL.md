# How to Install Repocheck

Since `repocheck` is a private tool you built yourself, it is **not** on Chocolatey or any public app store yet. You have to "install" it manually.

## Option 1: Run it from the current folder (Easiest)
You must be in the folder where `repocheck.exe` exists.

1.  Open Command Prompt or PowerShell.
2.  `cd` to your project folder:
    ```powershell
    cd c:\Users\GS2330\Documents\Repocheck-go
    ```
3.  Run it using `.\`:
    ```powershell
    .\repocheck.exe scan .
    ```

## Option 2: Add to User PATH (Office Friendly)
This method **does not require Admin rights** and is safe for office PCs.

1.  Create a folder in your user directory, e.g., `C:\Users\GS2330\bin`.
2.  Copy `repocheck.exe` into this new folder.
3.  Add it to your **User PATH**:
    *   Press **Win+S**, type **"env"**, and select **"Edit environment variables for your account"**. (Make sure it says "for your account", not "system").
    *   In the top box (**User variables for GS2330**), find **Path** and double-click it.
    *   Click **New** and paste the full path: `C:\Users\GS2330\bin`.
    *   Click **OK** on all windows.
4.  **Restart your Terminal** (close and reopen VS Code/PowerShell).
5.  Try it: Type `repocheck` in the terminal.

### Why this is safe
*   You are only changing settings for **your** user.
*   You are putting files in **your** folder.
*   It does not affect other users or system stability.

## Why is it not in Chocolatey?
Chocolatey is a package manager for published software. To get `repocheck` on Chocolatey, we would need to:
1.  Release a public version on GitHub.
2.  Write a Chocolatey install script.
3.  Submit it to the Chocolatey repository.
4.  Wait for approval.

For now, using **Option 2** is the standard way to install your own tools.
