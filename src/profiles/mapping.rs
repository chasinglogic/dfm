use std::{fmt::Display, path::PathBuf};

use log::debug;
use regex::Regex;

use crate::utils;

#[derive(Debug, Copy, Clone, PartialEq, Eq, serde::Deserialize, serde::Serialize)]
pub enum OS {
    Linux,
    Darwin,
    Windows,
}

impl Display for OS {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.write_str(match self {
            OS::Linux => "Linux",
            OS::Darwin => "Darwin",
            OS::Windows => "Windows",
        })
    }
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
    fn is_this_os(&self) -> bool {
        match self {
            TargetOS::All => true,
            TargetOS::Vec(targets) => targets.contains(&CURRENT_OS),
            TargetOS::String(desired) => *desired == CURRENT_OS,
        }
    }

    fn is_specifically_this_os(&self) -> bool {
        match self {
            TargetOS::All => false,
            _ => self.is_this_os(),
        }
    }
}

#[derive(Debug, Clone, serde::Deserialize, serde::Serialize)]
pub struct Mapping {
    #[serde(rename = "match", with = "serde_regex")]
    term: Regex,
    #[serde(default)]
    link_as_dir: Option<bool>,
    #[serde(default)]
    skip: Option<bool>,
    #[serde(default)]
    target_os: TargetOS,
    #[serde(default)]
    dest: Option<String>,
    #[serde(default)]
    target_dir: Option<String>,
}

impl Mapping {
    fn new(term: &str) -> Mapping {
        Mapping {
            term: Regex::new(term).expect("Unable to compile regex!"),
            link_as_dir: None,
            skip: None,
            target_os: TargetOS::All,
            dest: None,
            target_dir: None,
        }
    }

    fn skip(term: &str) -> Mapping {
        let mut mapping = Mapping::new(term);
        mapping.skip = Some(true);
        mapping
    }

    fn does_match(&self, path: &str) -> bool {
        let is_match = self.term.is_match(path);
        debug!("{} matches {:?} = {}", path, self.term, is_match);
        is_match

        // if is_match {
        //     TargetOS::is_this_os(&self.target_os)
        // } else {
        //     false
        // }
    }
}

#[derive(Default, Debug, Clone, PartialEq, Eq)]
pub enum MapAction {
    NewDest(PathBuf),
    NewTargetDir(PathBuf),
    LinkAsDir,
    Skip,
    #[default]
    None,
}

impl Display for MapAction {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        let repr = match self {
            Self::NewDest(value) => format!("MapAction::NewDest({})", value.display()),
            Self::NewTargetDir(value) => format!("MapAction::NewTargetDir({})", value.display()),
            Self::LinkAsDir => "MapAction::LinkAsDir".to_string(),
            Self::Skip => "MapAction::Skip".to_string(),
            Self::None => "MapAction::None".to_string(),
        };

        f.write_str(&repr)
    }
}

impl From<Mapping> for MapAction {
    fn from(mapping: Mapping) -> MapAction {
        MapAction::from(&mapping)
    }
}

fn check_if_skip_action(mapping: &Mapping) -> Option<MapAction> {
    match mapping.skip {
        Some(true) if mapping.target_os.is_this_os() => Some(MapAction::Skip),
        Some(_) => Some(MapAction::None),
        None => None,
    }
}

fn check_if_dest_action(mapping: &Mapping) -> Option<MapAction> {
    match mapping.dest {
        Some(ref dest) if mapping.target_os.is_this_os() => {
            Some(MapAction::NewDest(utils::expand_path(dest)))
        }
        Some(_) => Some(MapAction::None),
        None => None,
    }
}

fn check_if_target_dir_action(mapping: &Mapping) -> Option<MapAction> {
    match mapping.target_dir {
        Some(ref target_dir) if mapping.target_os.is_this_os() => {
            Some(MapAction::NewTargetDir(utils::expand_path(target_dir)))
        }
        Some(_) => Some(MapAction::None),
        None => None,
    }
}

fn check_if_link_as_dir_action(mapping: &Mapping) -> Option<MapAction> {
    match mapping.link_as_dir {
        Some(true) if mapping.target_os.is_this_os() => Some(MapAction::LinkAsDir),
        Some(_) => Some(MapAction::None),
        None => None,
    }
}

fn check_if_link_on_target_os_action(mapping: &Mapping) -> Option<MapAction> {
    if !mapping.target_os.is_specifically_this_os() {
        Some(MapAction::Skip)
    } else {
        None
    }
}

impl From<&Mapping> for MapAction {
    fn from(mapping: &Mapping) -> MapAction {
        let mapping_parsers = [
            check_if_skip_action,
            check_if_dest_action,
            check_if_target_dir_action,
            check_if_link_as_dir_action,
            // Must come last
            check_if_link_on_target_os_action,
        ];

        for f in mapping_parsers {
            if let Some(action) = f(mapping) {
                return action;
            }
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
            Mapping::skip("^README.[a-z]+$"),
            Mapping::skip("^LICENSE$"),
            Mapping::skip("^\\.gitignore$"),
            Mapping::skip("^\\.git$"),
            Mapping::skip("^\\.dfm\\.yml"),
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

    fn get_invalid_target_oses() -> Vec<OS> {
        vec![OS::Linux, OS::Darwin, OS::Windows]
            .into_iter()
            .filter(|os| *os != CURRENT_OS)
            .collect()
    }

    fn get_invalid_target_os() -> OS {
        get_invalid_target_oses()[0].clone()
    }

    #[test]
    fn test_skip_map_action_from_mapping() {
        assert_eq!(MapAction::Skip, MapAction::from(Mapping::skip("README.*")))
    }

    #[test]
    fn test_is_not_skip_on_target_os() {
        let config = format!(
            r#"match: .*snippets.*
target_os: {}"#,
            CURRENT_OS
        );
        let mapping: Mapping =
            serde_yaml::from_str(config.as_ref()).expect("invalid yaml config in test!");

        assert_eq!(MapAction::None, MapAction::from(mapping))
    }

    #[test]
    fn test_is_skip_not_on_target_os() {
        env_logger::init();

        let targets = get_invalid_target_oses();
        for target in targets {
            let config = format!(
                r#"match: .*snippets.*
target_os: {}"#,
                target
            );
            let mapping: Mapping =
                serde_yaml::from_str(config.as_ref()).expect("invalid yaml config in test!");

            assert_eq!(MapAction::Skip, MapAction::from(mapping))
        }
    }

    #[test]
    fn test_matches_mapping() {
        let config = r#"match: .config/ghostty/macos-config
skip: true"#;
        let mapping: Mapping = serde_yaml::from_str(config).expect("invalid yaml config in test!");
        assert!(mapping.does_match(".config/ghostty/macos-config"))
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
            MapAction::NewDest("/some/new/path.txt".into()),
            MapAction::from(mapping)
        )
    }

    #[test]
    fn test_new_target_dir_map_action_from_mapping() {
        let config = r#"match: LICENSE
target_dir: /some/new/"#;
        let mapping: Mapping = serde_yaml::from_str(config).expect("invalid yaml config in test!");
        assert_eq!(
            MapAction::NewTargetDir("/some/new/".into()),
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
                assert_ne!(value, PathBuf::from("~/.LICENSE.txt"));
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
                assert_ne!(value, PathBuf::from("~/some/subfolder"));
                assert!(value.starts_with("/"));
            }
            _ => panic!("Reached what should be an unreachable path!"),
        }
    }

    #[test]
    fn test_link_as_dir_map_action_from_mapping_on_wrong_os() {
        let config = format!(
            r#"match: .*snippets.*
link_as_dir: true
target_os: {}"#,
            get_invalid_target_os()
        );
        let mapping: Mapping =
            serde_yaml::from_str(config.as_ref()).expect("invalid yaml config in test!");
        assert_eq!(MapAction::None, MapAction::from(mapping))
    }

    #[test]
    fn test_new_dest_map_action_from_mapping_on_wrong_os() {
        let config = format!(
            r#"match: LICENSE
dest: /some/new/path.txt
target_os: {}"#,
            get_invalid_target_os()
        );
        let mapping: Mapping =
            serde_yaml::from_str(config.as_ref()).expect("invalid yaml config in test!");
        assert_eq!(MapAction::None, MapAction::from(mapping))
    }

    #[test]
    fn test_new_target_dir_map_action_from_mapping_on_wrong_os() {
        let config = format!(
            r#"match: LICENSE
target_dir: /some/new/
target_os: {}"#,
            get_invalid_target_os()
        );
        let mapping: Mapping =
            serde_yaml::from_str(config.as_ref()).expect("invalid yaml config in test!");
        assert_eq!(MapAction::None, MapAction::from(mapping))
    }
}
