use regex::Regex;

#[derive(Default, Debug, Clone, PartialEq, Eq, serde::Deserialize, serde::Serialize)]
pub enum TargetOS {
    Darwin,
    Linux,
    Windows,
    Vec(Vec<TargetOS>),
    #[default]
    All,
}

impl TargetOS {
    #[cfg(target_os = "macos")]
    fn current() -> TargetOS {
        TargetOS::Darwin
    }

    #[cfg(target_os = "linux")]
    fn current() -> TargetOS {
        TargetOS::Linux
    }

    #[cfg(target_os = "windows")]
    fn current() -> TargetOS {
        TargetOS::Windows
    }

    fn is_this_os(target: &TargetOS) -> bool {
        match target {
            &TargetOS::All => true,
            TargetOS::Vec(targets) => targets
                .into_iter()
                .find(|t| TargetOS::is_this_os(*t))
                .is_some(),
            desired => desired == &TargetOS::current(),
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

    fn skip(mut self) -> Mapping {
        self.skip = true;
        self
    }

    fn link_as_dir(mut self) -> Mapping {
        self.link_as_dir = true;
        self
    }

    fn dest(mut self, new_dest: String) -> Mapping {
        self.dest = new_dest;
        self
    }

    fn target_os(mut self, new_target: TargetOS) -> Mapping {
        self.target_os = new_target;
        self
    }

    fn does_match(&self, path: &str) -> bool {
        if self.term.is_match(path) {
            return TargetOS::is_this_os(&self.target_os);
        }

        false
    }
}

pub struct Mapper {
    mappings: Vec<Mapping>,
}

pub enum MapAction {
    NewDest(String),
    NewTargetDir(String),
    Skip,
    None,
}

impl From<Vec<Mapping>> for Mapper {
    fn from(mappings: Vec<Mapping>) -> Mapper {
        Mapper { mappings }
    }
}

impl From<Option<Vec<Mapping>>> for Mapper {
    fn from(mappings: Option<Vec<Mapping>>) -> Mapper {
        let configured: Vec<Mapping> = match mappings {
            Some(configured) => configured,
            None => vec![
                Mapping::new("README.*").skip(),
                Mapping::new("LICENSE").skip(),
            ],
        };

        Mapper::from(configured)
    }
}

impl Mapper {
    fn get_mapped_action(&self, relative_path: &str) -> MapAction {
        for mapping in &self.mappings {
            if !mapping.does_match(relative_path) {
                continue;
            }
        }

        MapAction::None
    }
}
