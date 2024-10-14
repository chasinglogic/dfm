use std::{collections::HashMap, io, path::Path, process::Command};

use log::debug;

#[derive(Debug, Clone, serde::Deserialize, serde::Serialize)]
pub struct HookDefinition {
    interpreter: String,
    script: String,
}

#[derive(Debug, Clone, serde::Deserialize, serde::Serialize)]
#[serde(untagged)]
pub enum Hook {
    String(String),
    HookDefinition(HookDefinition),
}

#[derive(Debug, Clone, serde::Deserialize, serde::Serialize)]
pub struct Hooks(HashMap<String, Vec<Hook>>);

impl Hooks {
    pub fn new() -> Hooks {
        Hooks(HashMap::new())
    }

    pub fn run_hook(&self, name: &str, working_directory: &Path) -> Result<(), io::Error> {
        debug!("running hook {}", name);

        match self.0.get(name) {
            Some(hooks) => {
                for hook in hooks {
                    let (interpreter_command, script): (&str, &str) = match hook {
                        Hook::String(script) => ("sh -c", script.as_ref()),
                        Hook::HookDefinition(HookDefinition {
                            interpreter,
                            script,
                        }) => (interpreter.as_ref(), script.as_ref()),
                    };

                    debug!("hook: {} {}", interpreter_command, script);

                    let mut argv = shlex::split(interpreter_command).ok_or_else(|| {
                        io::Error::new(
                            io::ErrorKind::InvalidInput,
                            format!("malformed interpreter: {}", &interpreter_command),
                        )
                    })?;
                    argv.push(script.to_string());

                    let shell = argv
                        .drain(0..1)
                        .next()
                        .expect("Unable to determine interpreter!");

                    Command::new(shell)
                        .args(&argv)
                        .current_dir(working_directory)
                        .spawn()
                        .expect("Unable to start shell!")
                        .wait()?;
                }

                Ok(())
            }
            None => Ok(()),
        }
    }
}
