use regex::Regex;

#[derive(Debug, Clone, PartialEq, Eq, serde::Deserialize, serde::Serialize)]
pub enum OS {
    Linux,
    Darwin,
    Windows,
}

#[derive(Default, Debug, Clone, PartialEq, Eq, serde::Deserialize, serde::Serialize)]
#[serde(untagged)]
pub enum TargetOS {
    String(OS),
    Vec(Vec<OS>),
    #[default]
    All,
}

#[cfg(target_os = "linux")]
const CURRENT_OS: OS = OS::Linux;
#[cfg(target_os = "macos")]
const CURRENT_OS: OS = OS::Darwin;
#[cfg(target_os = "windows")]
const CURRENT_OS: OS = OS::Windows;

impl TargetOS {
    fn is_this_os(target: &TargetOS) -> bool {
        match target {
            &TargetOS::All => true,
            TargetOS::Vec(targets) => targets.iter().any(|t| *t == CURRENT_OS),
            TargetOS::String(desired) => *desired == CURRENT_OS,
        }
    }
}

fn default_off() -> bool {
    false
}

#[derive(Debug, Clone, serde::Deserialize, serde::Serialize)]
pub struct Mapping {
    #[serde(rename = "match", with = "serde_regex")]
    term: Regex,
    #[serde(default = "default_off")]
    link_as_dir: bool,
    #[serde(default = "default_off")]
    skip: bool,
    #[serde(default)]
    target_os: TargetOS,
    #[serde(default)]
    dest: String,
    #[serde(default)]
    target_dir: String,
}

impl Mapping {
    fn new(term: &str) -> Mapping {
        Mapping {
            term: Regex::new(term).expect("Unable to compile regex!"),
            link_as_dir: false,
            skip: false,
            target_os: TargetOS::All,
            dest: "".to_string(),
            target_dir: "".to_string(),
        }
    }

    fn skip(term: &str) -> Mapping {
        let mut mapping = Mapping::new(term);
        mapping.skip = true;
        mapping
    }

    fn does_match(&self, path: &str) -> bool {
        if self.term.is_match(path) {
            return TargetOS::is_this_os(&self.target_os);
        }

        false
    }
}

#[derive(Default, Debug, Clone, PartialEq, Eq)]
pub enum MapAction {
    NewDest(String),
    NewTargetDir(String),
    LinkAsDir,
    Skip,
    #[default]
    None,
}

impl From<Mapping> for MapAction {
    fn from(mapping: Mapping) -> MapAction {
        MapAction::from(&mapping)
    }
}

impl From<&Mapping> for MapAction {
    fn from(mapping: &Mapping) -> MapAction {
        if mapping.skip {
            return MapAction::Skip;
        }

        if !mapping.dest.is_empty() {
            return MapAction::NewDest(shellexpand::tilde(mapping.dest.as_str()).into_owned());
        }

        if !mapping.target_dir.is_empty() {
            return MapAction::NewTargetDir(
                shellexpand::tilde(mapping.target_dir.as_str()).into_owned(),
            );
        }

        if mapping.link_as_dir {
            return MapAction::LinkAsDir;
        }

        MapAction::default()
    }
}

pub struct Mapper {
    mappings: Vec<Mapping>,
}

impl From<Vec<Mapping>> for Mapper {
    fn from(mappings: Vec<Mapping>) -> Mapper {
        Mapper { mappings }
    }
}

impl From<Option<Vec<Mapping>>> for Mapper {
    fn from(mappings: Option<Vec<Mapping>>) -> Mapper {
        let mut configured = mappings.unwrap_or_default();

        let default_mappings = vec![
            Mapping::skip("README.*"),
            Mapping::skip("LICENSE"),
            Mapping::skip(".gitignore$"),
            Mapping::skip(".git/?$"),
            Mapping::skip(".dfm.yml"),
        ];
        configured.extend(default_mappings);

        Mapper::from(configured)
    }
}

impl Mapper {
    pub fn get_mapped_action(&self, relative_path: &str) -> MapAction {
        for mapping in &self.mappings {
            if mapping.does_match(relative_path) {
                return MapAction::from(mapping);
            }
        }

        MapAction::None
    }
}

#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn test_skip_map_action_from_mapping() {
        assert_eq!(MapAction::Skip, MapAction::from(Mapping::skip("README.*")))
    }

    #[test]
    fn test_link_as_dir_map_action_from_mapping() {
        let config = r#"match: .*snippets.*
link_as_dir: true"#;
        let mapping: Mapping = serde_yaml::from_str(config).expect("invalid yaml config in test!");
        assert_eq!(MapAction::LinkAsDir, MapAction::from(mapping))
    }

    #[test]
    fn test_new_dest_map_action_from_mapping() {
        let config = r#"match: LICENSE
dest: /some/new/path.txt"#;
        let mapping: Mapping = serde_yaml::from_str(config).expect("invalid yaml config in test!");
        assert_eq!(
            MapAction::NewDest("/some/new/path.txt".to_string()),
            MapAction::from(mapping)
        )
    }

    #[test]
    fn test_new_target_dir_map_action_from_mapping() {
        let config = r#"match: LICENSE
target_dir: /some/new/"#;
        let mapping: Mapping = serde_yaml::from_str(config).expect("invalid yaml config in test!");
        assert_eq!(
            MapAction::NewTargetDir("/some/new/".to_string()),
            MapAction::from(mapping)
        )
    }

    #[test]
    fn test_new_dest_map_action_expands_tilde() {
        let config = r#"match: LICENSE
dest: ~/.LICENSE.txt"#;
        let mapping: Mapping = serde_yaml::from_str(config).expect("invalid yaml config in test!");
        let action = MapAction::from(mapping);

        match action {
            MapAction::NewDest(value) => {
                assert_ne!(value, "~/.LICENSE.txt");
                assert!(value.starts_with("/"));
            }
            _ => panic!("Reached what should be an unreachable path!"),
        }
    }

    #[test]
    fn test_new_target_dir_map_action_expands_tilde() {
        let config = r#"match: LICENSE
target_dir: ~/some/subfolder"#;
        let mapping: Mapping = serde_yaml::from_str(config).expect("invalid yaml config in test!");
        let action = MapAction::from(mapping);

        match action {
            MapAction::NewTargetDir(value) => {
                assert_ne!(value, "~/some/subfolder");
                assert!(value.starts_with("/"));
            }
            _ => panic!("Reached what should be an unreachable path!"),
        }
    }
}
