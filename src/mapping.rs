use std::io;
use std::path::Path;
use std::process;

use regex;
use serde::Deserialize;

use crate::link;

#[derive(Debug, Deserialize, Clone)]
pub struct MappingConfig {
    #[serde(rename = "match")]
    match_str: String,
    target_dir: Option<String>,
    target_os: Option<Vec<String>>,
    skip: Option<bool>,
}

#[derive(Debug, Clone)]
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

impl From<&str> for Mapping {
    fn from(rgx: &str) -> Mapping {
        Mapping::from(MappingConfig {
            match_str: rgx.to_string(),
            skip: Some(true),
            target_dir: None,
            target_os: None,
        })
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

#[derive(Debug, Clone)]
pub struct Mappings {
    mappings: Vec<Mapping>,
}

impl From<Vec<MappingConfig>> for Mappings {
    fn from(cfgs: Vec<MappingConfig>) -> Mappings {
        let mut defaults = vec![
            Mapping::from("/\\.git/"),
            Mapping::from("\\.dfm\\.yml$"),
            Mapping::from("\\.gitignore$"),
            Mapping::from("LICENSE(\\.md)?$"),
            Mapping::from("README(\\.md)?$"),
        ];
        let mut user_mappings = cfgs.iter().map(Mapping::from).collect();
        defaults.append(&mut user_mappings);
        Mappings { mappings: defaults }
    }
}

impl Mappings {
    pub fn link(&self, from: &Path, target_dir: &Path, overwrite: bool) -> Result<Vec<link::Info>, io::Error> {
        let mut wkd = walkdir::WalkDir::new(from).into_iter();
        let mut links = Vec::new();

        'dir: loop {
            let dir_entry = match wkd.next() {
                Some(Ok(d)) => d,
                Some(Err(e)) => {
                    let str_err = format!("{}", e);
                    return Err(e
                        .into_io_error()
                        .unwrap_or(io::Error::new(io::ErrorKind::Other, str_err)));
                }
                None => return Ok(links),
            };

            let src = dir_entry.path();
            let mut info = link::Info::new(src, from, target_dir);
            info.overwrite = overwrite;

            'mapping: for mapping in self.mappings.iter() {
                if !mapping.matches(src) {
                    continue 'mapping;
                }

                match mapping.change(&mut info) {
                    Some(new_info) => {
                        info = new_info;
                        break 'mapping;
                    }
                    None => {
                        if dir_entry.file_type().is_dir() {
                            wkd.skip_current_dir();
                        }

                        continue 'dir;
                    }
                }
            }

            if dir_entry.file_type().is_dir() {
                continue;
            }

            links.push(info);
        }
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
        let mut info = link_info_fixture(&path);
        assert!(mapping.change(&mut info).is_none());
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
        let mut info = link_info_fixture(&src);
        assert!(mapping.matches(&src));
        let new_info = mapping.change(&mut info);
        assert!(new_info.is_some());
        assert_eq!(new_info.unwrap().dst, PathBuf::from("/etc/mongod.conf"));
    }

    #[test]
    fn test_target_os() {
        let src = Path::new("/home/foo/.config/dfm/profile/bar/.bashrc");
        let mut info = link_info_fixture(&src);

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
            assert!(darwin_mapping.change(&mut info).is_some());
        } else {
            assert!(darwin_mapping.change(&mut info).is_none());
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
            assert!(linux_mapping.change(&mut info).is_some());
        } else {
            assert!(linux_mapping.change(&mut info).is_none());
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
            assert!(windows_mapping.change(&mut info).is_some());
        } else {
            assert!(windows_mapping.change(&mut info).is_none());
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
            assert!(unix_mapping.change(&mut info).is_some());
        } else {
            assert!(unix_mapping.change(&mut info).is_none());
        }
    }
}
