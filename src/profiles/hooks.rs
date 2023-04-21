use std::collections::HashMap;

use serde;

#[derive(Debug, serde::Deserialize)]
pub struct HookDefinition {
    interpreter: String,
    script: String,
}

#[derive(Debug, serde::Deserialize)]
#[serde(untagged)]
pub enum Hook {
    String(String),
    HookDefinition(HookDefinition),
}

// TODO: Should this be a struct so the user knows they have useless hooks?
pub type Hooks = HashMap<String, Vec<Hook>>;
