use std::collections::HashMap;

use serde;

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
}
