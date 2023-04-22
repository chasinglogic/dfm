use std::{fs::File, io::BufReader, path::Path};

use serde;

use super::hooks::Hooks;

#[derive(Debug, Clone, PartialEq, Eq, serde::Deserialize, serde::Serialize)]
#[serde(rename_all = "snake_case")]
enum LinkMode {
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
            module.expand();
        }

        config
    }

    fn expand(&mut self) {
        if self.location.starts_with("~") {
            self.location = shellexpand::tilde(&self.location).to_string();
        }
    }
}
