use std::collections::HashMap;
use std::io::Error;
use std::path::Path;
use std::process;

pub type HookConfig = HashMap<String, Vec<String>>;

#[derive(Debug, Clone)]
pub struct Hooks {
    config: HookConfig,
}

impl From<HookConfig> for Hooks {
    fn from(cfg: HookConfig) -> Hooks {
        Hooks { config: cfg }
    }
}

impl Hooks {
    pub fn run(&self, wd: &Path, name: &str) -> Result<(), Error> {
        let cmds = match self.config.get(name) {
            Some(c) => c,
            None => return Ok(()),
        };

        for cmd in cmds.iter() {
            let mut child = process::Command::new("sh");
            child.stdin(process::Stdio::inherit());
            child.stdout(process::Stdio::inherit());
            child.stderr(process::Stdio::inherit());
            child.args(vec!["-c", cmd]);
            let mut proc = child
                .current_dir(wd)
                .spawn()
                .expect("failed to start process");
            // TODO: handle error better here
            proc.wait().expect("failed to execute process");
        }

        Ok(())
    }

    pub fn pre(&self, wd: &Path, name: &str) -> Result<(), Error> {
        self.run(wd, &format!("before_{}", name))
    }

    pub fn post(&self, wd: &Path, name: &str) -> Result<(), Error> {
        self.run(wd, &format!("after_{}", name))
    }
}
