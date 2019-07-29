use std::io;
use std::path::{Path, PathBuf};
use std::process;

#[derive(Debug, Clone)]
pub struct Repo {
    pub path: PathBuf,
}

impl Repo {
    pub fn new(path: &Path) -> Repo {
        Repo {
            path: path.to_path_buf(),
        }
    }

    pub fn git(&self, cmd: &[&str]) -> Result<(), io::Error> {
        let mut child = process::Command::new("git");
        child.stdin(process::Stdio::inherit());
        child.stdout(process::Stdio::inherit());
        child.stderr(process::Stdio::inherit());
        let args = cmd;
        child.args(args);
        let mut proc = child.current_dir(&self.path).spawn()?;
        match proc.wait() {
            Ok(_) => Ok(()),
            Err(e) => Err(e),
        }
    }

    // TODO: would using libgit be better here / in general?
    pub fn is_dirty(&self) -> bool {
        match process::Command::new("git")
            .args(&["status", "--porcelain"])
            .current_dir(&self.path)
            .output()
        {
            Ok(proc) => {
                let output = std::str::from_utf8(&proc.stdout).unwrap_or("");
                output != ""
            }
            Err(_) => false,
        }
    }

    fn has_remote(&self) -> bool {
        match process::Command::new("git")
            .args(&["remote", "--verbose"])
            .current_dir(&self.path)
            .output()
        {
            Ok(proc) => {
                let output = std::str::from_utf8(&proc.stdout).unwrap_or_default();
                output.contains("origin")
            }
            Err(_) => false,
        }
    }

    pub fn sync(&self, msg: &str, pull_only: bool) -> Result<(), io::Error> {
        let repo_is_dirty = self.is_dirty();
        if repo_is_dirty && !pull_only {
            self.git(&["add", "--all"])?;
            self.git(&["commit", "--message", &msg])?;
        }

        if self.has_remote() {
            self.git(&["pull", "--rebase", "origin", "master"])?;

            if repo_is_dirty && !pull_only {
                self.git(&["push", "origin", "master"])?;
            }
        }

        Ok(())
    }
}
