use std::{fs::File, io::BufReader, path::Path};

use serde;

use super::hooks::Hooks;

#[derive(Debug, Clone, PartialEq, Eq, serde::Deserialize, serde::Serialize)]
#[serde(rename_all = "snake_case")]
pub enum LinkMode {
    Pre,
    Post,
    None,
}

impl Default for LinkMode {
    fn default() -> Self {
        LinkMode::Post
    }
}

fn default_off() -> bool {
    false
}

#[derive(Debug, Clone, serde::Deserialize, serde::Serialize)]
pub struct DFMConfig {
    #[serde(default, alias = "repository")]
    pub repo: String,
    #[serde(default)]
    pub location: String,
    #[serde(default = "Hooks::new")]
    pub hooks: Hooks,

    #[serde(default = "default_off")]
    pub prompt_for_commit_message: bool,
    #[serde(default = "default_off")]
    pub pull_only: bool,
    #[serde(default)]
    pub link: LinkMode,
    #[serde(default = "Vec::new")]
    pub modules: Vec<DFMConfig>,
}

impl Default for DFMConfig {
    fn default() -> Self {
        DFMConfig {
            prompt_for_commit_message: false,
            pull_only: false,
            link: LinkMode::default(),
            repo: "".to_string(),
            location: "".to_string(),
            hooks: Hooks::new(),
            modules: Vec::new(),
        }
    }
}

impl DFMConfig {
    pub fn load(file: &Path) -> DFMConfig {
        let fh = File::open(file).expect(
            format!(
                "Unexpected error reading {}",
                file.to_str().unwrap_or(".dfm.yml")
            )
            .as_str(),
        );
        let reader = BufReader::new(fh);
        let mut config: DFMConfig = serde_yaml::from_reader(reader).expect("Malformed .dfm.yml");
        if config.location == "" {
            config.location = file
                .parent()
                .expect("Unexpected error getting profile location!")
                .to_str()
                .expect("Unexpected error turning profile location to a string!")
                .to_string();
        }

        for module in &mut config.modules {
            module.expand_module();
        }

        config
    }

    fn expand_module(&mut self) {
        if self.location.starts_with("~") {
            self.location = shellexpand::tilde(&self.location).to_string();
        }

        if self.location == "" {
            let name = self
                .repo
                .split("/")
                .last()
                .expect("A module must define a repository!")
                .replace(".git", "");

            let mut module_dir = crate::cli::state::modules_dir();
            module_dir.push(name);

            self.location = module_dir.to_string_lossy().to_string();
        }
    }
}
