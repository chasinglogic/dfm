use std::io;
use std::path::Path;
use std::process;

use regex;
use serde::Deserialize;

use crate::link;

#[derive(Deserialize, Clone)]
pub struct MappingConfig {
    #[serde(rename = "match")]
    match_str: String,
    target_dir: Option<String>,
    target_os: Option<Vec<String>>,
    skip: Option<bool>,
}

pub struct Mapping {
    rgx: regex::Regex,
    config: MappingConfig,
}

impl From<MappingConfig> for Mapping {
    fn from(cfg: MappingConfig) -> Mapping {
        Mapping {
            rgx: match regex::Regex::new(&cfg.match_str) {
                Ok(r) => r,
                // TODO: Maybe don't crash here?
                Err(e) => {
                    println!("unable to compile regex for mapping: {}", e);
                    process::exit(1);
                }
            },
            config: cfg,
        }
    }
}

impl From<&MappingConfig> for Mapping {
    fn from(cfg: &MappingConfig) -> Mapping {
        Mapping::from(cfg.clone())
    }
}

impl Mapping {
    pub fn matches(&self, path: &Path) -> bool {
        self.rgx.is_match(&path.as_os_str().to_string_lossy())
    }

    pub fn change(&self, info: &mut link::Info) -> Option<link::Info> {
        if self.config.skip.unwrap_or(false) {
            return None;
        }

        if let Some(target_os) = &self.config.target_os {
            if cfg!(target_os = "macos") && !target_os.contains(&"Darwin".to_string()) {
                return None;
            }

            if cfg!(target_os = "linux") && !target_os.contains(&"Linux".to_string()) {
                return None;
            }

            if cfg!(windows) && !target_os.contains(&"Windows".to_string()) {
                return None;
            }
        }

        if let Some(new_target) = &self.config.target_dir {
            info.retarget(Path::new(&new_target));
        }

        Some(info.clone())
    }
}

pub struct Mappings {
    mappings: Vec<Mapping>,
}

impl From<Vec<MappingConfig>> for Mappings {
    fn from(cfgs: Vec<MappingConfig>) -> Mappings {
        Mappings {
            mappings: cfgs.iter().map(Mapping::from).collect(),
        }
    }
}

impl Mappings {
    pub fn link(&self, from: &Path, target_dir: &Path) -> Result<Vec<link::Info>, io::Error> {
        let wkd: Vec<Result<walkdir::DirEntry, walkdir::Error>> = walkdir::WalkDir::new(from)
            .into_iter()
            .filter_entry(|e| e.file_type().is_file())
            .collect();
        let mut links = Vec::with_capacity(wkd.len());

        'dir: for dir_entry in wkd.iter() {
            if let Ok(path) = dir_entry {
                let src = path.path();
                let mut info = link::Info::new(src, from, target_dir);

                for mapping in self.mappings.iter() {
                    if !mapping.matches(src) {
                        continue;
                    }

                    match mapping.change(&mut info) {
                        Some(new_info) => {
                            info = new_info;
                            break;
                        }
                        None => continue 'dir,
                    }
                }

                links.push(info);
            }
        }

        Ok(links)
    }
}

#[cfg(test)]
pub mod test {
    use super::*;
    use crate::link::Info;
    use serde_yaml::from_str;
    use std::path::PathBuf;

    fn link_info_fixture(src: &Path) -> Info {
        Info::new(
            src,
            Path::new("/home/foo/.config/dfm/profile/bar"),
            Path::new("/home/foo"),
        )
    }

    #[test]
    fn test_matches() {
        let path = Path::new("/home/foo/.config/dfm/profile/bar/path");
        let mapping_config: MappingConfig = from_str(
            "match: path
skip: true
",
        )
        .unwrap();
        let mapping = Mapping::from(mapping_config);
        assert!(mapping.matches(path));
        let info = link_info_fixture(&path);
        assert!(mapping.change(info).is_none());
    }

    #[test]
    fn test_skips() {
        let path = Path::new("/some/foo/bar/path");
        let mapping_config: MappingConfig = from_str(
            "match: path
skip: true
",
        )
        .unwrap();
        let mapping = Mapping::from(mapping_config);
        assert!(mapping.matches(path));
    }

    #[test]
    fn test_transforms_target_dir() {
        let mapping_config: MappingConfig = from_str(
            "match: mongod
target_dir: /etc
",
        )
        .unwrap();
        let mapping = Mapping::from(mapping_config);
        let src = Path::new("/home/foo/.config/dfm/profile/bar/mongod.conf");
        let info = link_info_fixture(&src);
        assert!(mapping.matches(&src));
        let new_info = mapping.change(info);
        assert!(new_info.is_some());
        assert_eq!(
            new_info.unwrap().get_dst(),
            PathBuf::from("/etc/mongod.conf")
        );
    }

    #[test]
    fn test_target_os() {
        let src = Path::new("/home/foo/.config/dfm/profile/bar/.bashrc");
        let info = link_info_fixture(&src);

        let darwin_mapping_config: MappingConfig = from_str(
            "match: bashrc
target_os: 
    - Darwin
",
        )
        .unwrap();
        let darwin_mapping = Mapping::from(darwin_mapping_config);
        assert!(darwin_mapping.matches(&src));
        if cfg!(target_os = "macos") {
            assert!(darwin_mapping.change(info.clone()).is_some());
        } else {
            assert!(darwin_mapping.change(info.clone()).is_none());
        }

        let linux_mapping_config: MappingConfig = from_str(
            "match: bashrc
target_os: 
    - Linux
",
        )
        .unwrap();
        let linux_mapping = Mapping::from(linux_mapping_config);
        assert!(linux_mapping.matches(&src));
        if cfg!(target_os = "linux") {
            assert!(linux_mapping.change(info.clone()).is_some());
        } else {
            assert!(linux_mapping.change(info.clone()).is_none());
        }

        let windows_mapping_config: MappingConfig = from_str(
            "match: bashrc
target_os: 
    - Windows
",
        )
        .unwrap();
        let windows_mapping = Mapping::from(windows_mapping_config);
        assert!(windows_mapping.matches(&src));
        if cfg!(target_os = "windows") {
            assert!(windows_mapping.change(info.clone()).is_some());
        } else {
            assert!(windows_mapping.change(info.clone()).is_none());
        }

        let unix_mapping_config: MappingConfig = from_str(
            "match: bashrc
target_os: 
    - Linux
    - Darwin
",
        )
        .unwrap();
        let unix_mapping = Mapping::from(unix_mapping_config);
        assert!(unix_mapping.matches(&src));
        if cfg!(unix) {
            assert!(unix_mapping.change(info.clone()).is_some());
        } else {
            assert!(unix_mapping.change(info.clone()).is_none());
        }
    }
}
