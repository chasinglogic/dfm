use serde;

use super::hooks::Hooks;

#[derive(Debug, serde::Deserialize)]
#[serde(rename_all = "snake_case")]
enum LinkMode {
    Default,
    None,
    Pre,
    Post,
}

impl Default for LinkMode {
    fn default() -> Self {
        LinkMode::Default
    }
}

fn default_off() -> bool {
    false
}

#[derive(Debug, serde::Deserialize)]
struct DFMConfig {
    #[serde(default = "default_off")]
    prompt_for_commit_message: bool,
    #[serde(default = "default_off")]
    pull_only: bool,
    #[serde(default)]
    link: LinkMode,
    #[serde(default)]
    location: String,

    hooks: Hooks,

    modules: Vec<DFMConfig>,
}

struct Profile {}
